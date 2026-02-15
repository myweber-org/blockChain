
package main

import (
	"fmt"
	"math"
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

	if len(swa.metrics) == 0 {
		return nil, fmt.Errorf("no metrics available in window")
	}

	values := make([]float64, len(swa.metrics))
	for i, metric := range swa.metrics {
		values[i] = metric.Value
	}

	stats := make(map[string]float64)
	stats["count"] = float64(len(values))
	stats["min"] = calculateMin(values)
	stats["max"] = calculateMax(values)
	stats["mean"] = calculateMean(values)
	stats["stddev"] = calculateStdDev(values, stats["mean"])

	for _, p := range swa.percentiles {
		key := fmt.Sprintf("p%.0f", p*100)
		stats[key], _ = calculatePercentile(values, p)
	}

	return stats, nil
}

func calculateMin(values []float64) float64 {
	min := math.MaxFloat64
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func calculateMax(values []float64) float64 {
	max := -math.MaxFloat64
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

func calculateMean(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateStdDev(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0.0
	}
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(values)-1))
}

func calculatePercentile(values []float64, percentile float64) (float64, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("empty values slice")
	}
	if percentile < 0 || percentile > 1 {
		return 0, fmt.Errorf("percentile must be between 0 and 1")
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	index := percentile * float64(len(sorted)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return sorted[lower], nil
	}

	weight := index - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight, nil
}

func main() {
	aggregator := NewSlidingWindowAggregator(5*time.Minute, []float64{0.5, 0.95, 0.99})

	for i := 0; i < 100; i++ {
		value := 50 + 20*math.Sin(float64(i)/10)
		aggregator.AddMetric(value)
		time.Sleep(100 * time.Millisecond)
	}

	stats, err := aggregator.CalculateStats()
	if err != nil {
		fmt.Printf("Error calculating stats: %v\n", err)
		return
	}

	fmt.Println("Aggregated metrics:")
	for key, value := range stats {
		fmt.Printf("%s: %.4f\n", key, value)
	}
}