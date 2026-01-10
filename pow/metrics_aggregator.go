
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