
package metrics

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	windowSize  time.Duration
	maxSamples  int
	mu          sync.RWMutex
	samples     []float64
	timestamps  []time.Time
	percentiles []float64
}

func NewAggregator(windowSize time.Duration, maxSamples int, percentiles []float64) *Aggregator {
	return &Aggregator{
		windowSize:  windowSize,
		maxSamples:  maxSamples,
		percentiles: percentiles,
		samples:     make([]float64, 0, maxSamples),
		timestamps:  make([]time.Time, 0, maxSamples),
	}
}

func (a *Aggregator) Add(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	a.samples = append(a.samples, value)
	a.timestamps = append(a.timestamps, now)

	a.evictOldSamples(now)
	if len(a.samples) > a.maxSamples {
		a.samples = a.samples[1:]
		a.timestamps = a.timestamps[1:]
	}
}

func (a *Aggregator) evictOldSamples(now time.Time) {
	cutoff := now.Add(-a.windowSize)
	i := 0
	for i < len(a.timestamps) && a.timestamps[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		a.samples = a.samples[i:]
		a.timestamps = a.timestamps[i:]
	}
}

func (a *Aggregator) GetStats() map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.samples) == 0 {
		return nil
	}

	stats := make(map[string]float64)
	sorted := make([]float64, len(a.samples))
	copy(sorted, a.samples)
	sort.Float64s(sorted)

	stats["count"] = float64(len(a.samples))
	stats["min"] = sorted[0]
	stats["max"] = sorted[len(sorted)-1]

	var sum float64
	for _, v := range a.samples {
		sum += v
	}
	stats["mean"] = sum / float64(len(a.samples))

	for _, p := range a.percentiles {
		if p < 0 || p > 100 {
			continue
		}
		idx := int(float64(len(sorted)-1) * p / 100.0)
		key := "p" + formatPercentile(p)
		stats[key] = sorted[idx]
	}

	return stats
}

func formatPercentile(p float64) string {
	if p == float64(int(p)) {
		return formatInt(int(p))
	}
	return formatFloat(p)
}

func formatInt(n int) string {
	return string(rune('0' + n/10)) + string(rune('0' + n%10))
}

func formatFloat(f float64) string {
	s := formatInt(int(f))
	dec := int((f - float64(int(f))) * 100)
	if dec > 0 {
		s += "_" + formatInt(dec)
	}
	return s
}
package metrics

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	windowSize   time.Duration
	dataPoints   []float64
	timestamps   []time.Time
	mu           sync.RWMutex
	percentiles  []float64
}

func NewAggregator(windowSize time.Duration, percentiles []float64) *Aggregator {
	return &Aggregator{
		windowSize:  windowSize,
		dataPoints:  make([]float64, 0),
		timestamps:  make([]time.Time, 0),
		percentiles: percentiles,
	}
}

func (a *Aggregator) Add(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	now := time.Now()
	a.dataPoints = append(a.dataPoints, value)
	a.timestamps = append(a.timestamps, now)
	
	a.cleanup(now)
}

func (a *Aggregator) cleanup(currentTime time.Time) {
	cutoff := currentTime.Add(-a.windowSize)
	startIdx := 0
	
	for i, ts := range a.timestamps {
		if ts.After(cutoff) {
			startIdx = i
			break
		}
	}
	
	a.dataPoints = a.dataPoints[startIdx:]
	a.timestamps = a.timestamps[startIdx:]
}

func (a *Aggregator) GetStats() map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	if len(a.dataPoints) == 0 {
		return make(map[string]float64)
	}
	
	stats := make(map[string]float64)
	
	// Calculate basic statistics
	var sum float64
	minVal := a.dataPoints[0]
	maxVal := a.dataPoints[0]
	
	dataCopy := make([]float64, len(a.dataPoints))
	copy(dataCopy, a.dataPoints)
	sort.Float64s(dataCopy)
	
	for _, val := range a.dataPoints {
		sum += val
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}
	
	stats["count"] = float64(len(a.dataPoints))
	stats["min"] = minVal
	stats["max"] = maxVal
	stats["mean"] = sum / float64(len(a.dataPoints))
	
	// Calculate percentiles
	for _, p := range a.percentiles {
		if p < 0 || p > 100 {
			continue
		}
		
		idx := int(float64(len(dataCopy)-1) * p / 100.0)
		key := "p" + formatPercentile(p)
		stats[key] = dataCopy[idx]
	}
	
	return stats
}

func formatPercentile(p float64) string {
	if p == float64(int(p)) {
		return string(rune('0' + int(p)/10)) + string(rune('0' + int(p)%10))
	}
	return string(rune('0' + int(p)/10)) + string(rune('0' + int(p)%10)) + "5"
}

func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	a.dataPoints = make([]float64, 0)
	a.timestamps = make([]time.Time, 0)
}