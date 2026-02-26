
package metrics

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	windowSize   time.Duration
	percentiles  []float64
	measurements []measurement
	mu           sync.RWMutex
}

type measurement struct {
	timestamp time.Time
	value     float64
}

func NewAggregator(windowSize time.Duration, percentiles []float64) *Aggregator {
	for _, p := range percentiles {
		if p < 0 || p > 100 {
			panic("percentile must be between 0 and 100")
		}
	}
	sort.Float64s(percentiles)
	
	return &Aggregator{
		windowSize:  windowSize,
		percentiles: percentiles,
		measurements: make([]measurement, 0),
	}
}

func (a *Aggregator) Record(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	now := time.Now()
	a.measurements = append(a.measurements, measurement{
		timestamp: now,
		value:     value,
	})
	a.cleanup(now)
}

func (a *Aggregator) cleanup(now time.Time) {
	cutoff := now.Add(-a.windowSize)
	i := 0
	for i < len(a.measurements) && a.measurements[i].timestamp.Before(cutoff) {
		i++
	}
	if i > 0 {
		a.measurements = a.measurements[i:]
	}
}

func (a *Aggregator) GetStats() map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	now := time.Now()
	a.cleanup(now)
	
	if len(a.measurements) == 0 {
		return make(map[string]float64)
	}
	
	values := make([]float64, len(a.measurements))
	for i, m := range a.measurements {
		values[i] = m.value
	}
	sort.Float64s(values)
	
	stats := make(map[string]float64)
	stats["count"] = float64(len(values))
	stats["min"] = values[0]
	stats["max"] = values[len(values)-1]
	
	var sum float64
	for _, v := range values {
		sum += v
	}
	stats["mean"] = sum / float64(len(values))
	
	for _, p := range a.percentiles {
		key := formatPercentileKey(p)
		stats[key] = calculatePercentile(values, p)
	}
	
	return stats
}

func formatPercentileKey(p float64) string {
	if p == 50 {
		return "median"
	}
	return formatFloat(p) + "th"
}

func formatFloat(f float64) string {
	if f == float64(int(f)) {
		return string(rune(int(f)))
	}
	return string(rune(f))
}

func calculatePercentile(sorted []float64, p float64) float64 {
	if len(sorted) == 1 {
		return sorted[0]
	}
	
	index := (p / 100) * float64(len(sorted)-1)
	lower := int(index)
	upper := lower + 1
	
	if upper >= len(sorted) {
		return sorted[lower]
	}
	
	weight := index - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}