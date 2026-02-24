
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
}package metrics

import (
	"sync"
	"time"
)

type Aggregator struct {
	mu            sync.RWMutex
	latencySum    time.Duration
	requestCount  int64
	errorCount    int64
	lastResetTime time.Time
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		lastResetTime: time.Now(),
	}
}

func (a *Aggregator) Record(latency time.Duration, isError bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.latencySum += latency
	a.requestCount++
	if isError {
		a.errorCount++
	}
}

func (a *Aggregator) GetStats() (avgLatency time.Duration, errorRate float64, totalRequests int64) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.requestCount == 0 {
		return 0, 0, 0
	}

	avgLatency = a.latencySum / time.Duration(a.requestCount)
	errorRate = float64(a.errorCount) / float64(a.requestCount)
	return avgLatency, errorRate, a.requestCount
}

func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.latencySum = 0
	a.requestCount = 0
	a.errorCount = 0
	a.lastResetTime = time.Now()
}

func (a *Aggregator) TimeSinceReset() time.Duration {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return time.Since(a.lastResetTime)
}package metrics

import (
	"sync"
	"time"
)

type Aggregator struct {
	mu          sync.RWMutex
	windowSize  time.Duration
	metrics     []float64
	timestamps  []time.Time
	maxSamples  int
}

func NewAggregator(windowSize time.Duration, maxSamples int) *Aggregator {
	return &Aggregator{
		windowSize: windowSize,
		maxSamples: maxSamples,
		metrics:    make([]float64, 0, maxSamples),
		timestamps: make([]time.Time, 0, maxSamples),
	}
}

func (a *Aggregator) Add(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	a.cleanup(now)

	if len(a.metrics) < a.maxSamples {
		a.metrics = append(a.metrics, value)
		a.timestamps = append(a.timestamps, now)
	}
}

func (a *Aggregator) cleanup(currentTime time.Time) {
	cutoff := currentTime.Add(-a.windowSize)
	validStart := 0

	for i, ts := range a.timestamps {
		if ts.After(cutoff) {
			validStart = i
			break
		}
	}

	if validStart > 0 {
		a.metrics = a.metrics[validStart:]
		a.timestamps = a.timestamps[validStart:]
	}
}

func (a *Aggregator) Average() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	a.cleanup(time.Now())

	if len(a.metrics) == 0 {
		return 0
	}

	var sum float64
	for _, v := range a.metrics {
		sum += v
	}
	return sum / float64(len(a.metrics))
}

func (a *Aggregator) Count() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	a.cleanup(time.Now())
	return len(a.metrics)
}

func (a *Aggregator) Percentile(p float64) float64 {
	if p < 0 || p > 100 {
		return 0
	}

	a.mu.RLock()
	defer a.mu.RUnlock()
	a.cleanup(time.Now())

	if len(a.metrics) == 0 {
		return 0
	}

	values := make([]float64, len(a.metrics))
	copy(values, a.metrics)

	sorted := sortSlice(values)
	index := int(float64(len(sorted)-1) * p / 100)
	return sorted[index]
}

func sortSlice(s []float64) []float64 {
	sorted := make([]float64, len(s))
	copy(sorted, s)

	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j] < sorted[i] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}