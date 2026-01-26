
package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Metric struct {
	Timestamp time.Time
	Value     float64
}

type SlidingWindowAggregator struct {
	windowSize  time.Duration
	metrics     []Metric
	mu          sync.RWMutex
	percentiles []float64
}

func NewSlidingWindowAggregator(windowSize time.Duration, percentiles []float64) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize:  windowSize,
		metrics:     make([]Metric, 0),
		percentiles: percentiles,
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.metrics = append(swa.metrics, Metric{Timestamp: now, Value: value})
	swa.cleanupOldMetrics(now)
}

func (swa *SlidingWindowAggregator) cleanupOldMetrics(currentTime time.Time) {
	cutoff := currentTime.Add(-swa.windowSize)
	validStart := 0
	for i, metric := range swa.metrics {
		if metric.Timestamp.After(cutoff) {
			validStart = i
			break
		}
	}
	swa.metrics = swa.metrics[validStart:]
}

func (swa *SlidingWindowAggregator) CalculateStats() (map[string]float64, error) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	swa.cleanupOldMetrics(time.Now())

	if len(swa.metrics) == 0 {
		return nil, fmt.Errorf("no metrics available in window")
	}

	values := make([]float64, len(swa.metrics))
	for i, metric := range swa.metrics {
		values[i] = metric.Value
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

	for _, p := range swa.percentiles {
		if p < 0 || p > 100 {
			continue
		}
		index := (p / 100) * float64(len(values)-1)
		lower := int(index)
		upper := lower + 1
		weight := index - float64(lower)

		if upper >= len(values) {
			stats[fmt.Sprintf("p%.1f", p)] = values[lower]
		} else {
			stats[fmt.Sprintf("p%.1f", p)] = values[lower]*(1-weight) + values[upper]*weight
		}
	}

	return stats, nil
}

func main() {
	aggregator := NewSlidingWindowAggregator(5*time.Minute, []float64{50.0, 90.0, 95.0, 99.0})

	for i := 0; i < 100; i++ {
		aggregator.AddMetric(float64(i))
		time.Sleep(100 * time.Millisecond)
	}

	stats, err := aggregator.CalculateStats()
	if err != nil {
		fmt.Printf("Error calculating stats: %v\n", err)
		return
	}

	for k, v := range stats {
		fmt.Printf("%s: %.2f\n", k, v)
	}
}
package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Metric struct {
	Timestamp time.Time
	Value     float64
}

type Aggregator struct {
	mu           sync.RWMutex
	windowSize   time.Duration
	metrics      []Metric
	percentiles  []float64
}

func NewAggregator(windowSize time.Duration, percentiles []float64) *Aggregator {
	return &Aggregator{
		windowSize:  windowSize,
		percentiles: percentiles,
		metrics:     make([]Metric, 0),
	}
}

func (a *Aggregator) AddMetric(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	a.metrics = append(a.metrics, Metric{Timestamp: now, Value: value})
	a.cleanupOldMetrics(now)
}

func (a *Aggregator) cleanupOldMetrics(currentTime time.Time) {
	cutoff := currentTime.Add(-a.windowSize)
	i := 0
	for i < len(a.metrics) && a.metrics[i].Timestamp.Before(cutoff) {
		i++
	}
	if i > 0 {
		a.metrics = a.metrics[i:]
	}
}

func (a *Aggregator) GetStats() map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.metrics) == 0 {
		return make(map[string]float64)
	}

	values := make([]float64, len(a.metrics))
	for i, m := range a.metrics {
		values[i] = m.Value
	}

	stats := make(map[string]float64)
	stats["count"] = float64(len(values))
	stats["min"], stats["max"], stats["avg"] = calculateBasicStats(values)

	sort.Float64s(values)
	for _, p := range a.percentiles {
		key := fmt.Sprintf("p%.0f", p*100)
		stats[key] = calculatePercentile(values, p)
	}

	return stats
}

func calculateBasicStats(values []float64) (min, max, avg float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}
	min = values[0]
	max = values[0]
	sum := 0.0

	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}
	avg = sum / float64(len(values))
	return min, max, avg
}

func calculatePercentile(sortedValues []float64, percentile float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}
	index := percentile * float64(len(sortedValues)-1)
	lower := int(index)
	upper := lower + 1

	if upper >= len(sortedValues) {
		return sortedValues[lower]
	}

	weight := index - float64(lower)
	return sortedValues[lower]*(1-weight) + sortedValues[upper]*weight
}

func main() {
	agg := NewAggregator(5*time.Minute, []float64{0.5, 0.95, 0.99})

	for i := 0; i < 100; i++ {
		agg.AddMetric(float64(i) * 1.5)
		time.Sleep(100 * time.Millisecond)
	}

	stats := agg.GetStats()
	for k, v := range stats {
		fmt.Printf("%s: %.2f\n", k, v)
	}
}