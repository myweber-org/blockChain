
package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

type MetricPoint struct {
	Timestamp time.Time
	Value     float64
}

type SlidingWindowAggregator struct {
	mu          sync.RWMutex
	windowSize  time.Duration
	points      []MetricPoint
	percentiles []float64
}

func NewSlidingWindowAggregator(windowSize time.Duration, percentiles []float64) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize:  windowSize,
		points:      make([]MetricPoint, 0),
		percentiles: percentiles,
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.points = append(swa.points, MetricPoint{Timestamp: now, Value: value})
	swa.cleanupOldPoints(now)
}

func (swa *SlidingWindowAggregator) cleanupOldPoints(currentTime time.Time) {
	cutoff := currentTime.Add(-swa.windowSize)
	i := 0
	for i < len(swa.points) && swa.points[i].Timestamp.Before(cutoff) {
		i++
	}
	if i > 0 {
		swa.points = swa.points[i:]
	}
}

func (swa *SlidingWindowAggregator) CalculateStats() map[string]float64 {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	values := make([]float64, len(swa.points))
	for i, p := range swa.points {
		values[i] = p.Value
	}

	if len(values) == 0 {
		return map[string]float64{}
	}

	stats := make(map[string]float64)
	stats["count"] = float64(len(values))

	sum := 0.0
	min := math.MaxFloat64
	max := -math.MaxFloat64

	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	stats["sum"] = sum
	stats["avg"] = sum / float64(len(values))
	stats["min"] = min
	stats["max"] = max

	sort.Float64s(values)

	for _, p := range swa.percentiles {
		if p < 0 || p > 100 {
			continue
		}
		key := fmt.Sprintf("p%.0f", p)
		stats[key] = calculatePercentile(values, p)
	}

	return stats
}

func calculatePercentile(sortedValues []float64, percentile float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}

	index := (percentile / 100) * float64(len(sortedValues)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return sortedValues[lower]
	}

	weight := index - float64(lower)
	return sortedValues[lower]*(1-weight) + sortedValues[upper]*weight
}

func main() {
	aggregator := NewSlidingWindowAggregator(5*time.Minute, []float64{50, 95, 99})

	for i := 0; i < 100; i++ {
		aggregator.AddMetric(float64(i))
		time.Sleep(100 * time.Millisecond)
	}

	stats := aggregator.CalculateStats()
	for k, v := range stats {
		fmt.Printf("%s: %.2f\n", k, v)
	}
}