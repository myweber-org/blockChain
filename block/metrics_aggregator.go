
package metrics

import (
	"container/list"
	"sort"
	"sync"
	"time"
)

type MetricPoint struct {
	Value     float64
	Timestamp time.Time
}

type SlidingWindowAggregator struct {
	windowSize  time.Duration
	maxPoints   int
	points      *list.List
	mu          sync.RWMutex
	totalSum    float64
	totalCount  int64
}

func NewSlidingWindowAggregator(windowSize time.Duration, maxPoints int) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize: windowSize,
		maxPoints:  maxPoints,
		points:     list.New(),
	}
}

func (swa *SlidingWindowAggregator) Add(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.points.PushBack(MetricPoint{
		Value:     value,
		Timestamp: now,
	})

	swa.totalSum += value
	swa.totalCount++

	swa.evictExpired(now)
	swa.trimToMaxPoints()
}

func (swa *SlidingWindowAggregator) evictExpired(now time.Time) {
	cutoff := now.Add(-swa.windowSize)
	for {
		front := swa.points.Front()
		if front == nil {
			break
		}
		point := front.Value.(MetricPoint)
		if point.Timestamp.After(cutoff) {
			break
		}
		swa.points.Remove(front)
		swa.totalSum -= point.Value
		swa.totalCount--
	}
}

func (swa *SlidingWindowAggregator) trimToMaxPoints() {
	for swa.points.Len() > swa.maxPoints {
		front := swa.points.Front()
		point := front.Value.(MetricPoint)
		swa.points.Remove(front)
		swa.totalSum -= point.Value
		swa.totalCount--
	}
}

func (swa *SlidingWindowAggregator) GetStats() (float64, float64, float64) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if swa.totalCount == 0 {
		return 0, 0, 0
	}

	values := make([]float64, 0, swa.points.Len())
	for e := swa.points.Front(); e != nil; e = e.Next() {
		values = append(values, e.Value.(MetricPoint).Value)
	}
	sort.Float64s(values)

	mean := swa.totalSum / float64(swa.totalCount)
	median := calculatePercentile(values, 50)
	p95 := calculatePercentile(values, 95)

	return mean, median, p95
}

func calculatePercentile(sortedValues []float64, percentile float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}
	index := (percentile / 100) * float64(len(sortedValues)-1)
	lower := int(index)
	upper := lower + 1
	weight := index - float64(lower)

	if upper >= len(sortedValues) {
		return sortedValues[lower]
	}
	return sortedValues[lower]*(1-weight) + sortedValues[upper]*weight
}package main

import (
	"log"
	"sync"
	"time"
)

type Metrics struct {
	mu            sync.RWMutex
	requestCount  int64
	errorCount    int64
	totalLatency  time.Duration
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RecordRequest(latency time.Duration, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requestCount++
	m.totalLatency += latency
	if isError {
		m.errorCount++
	}
}

func (m *Metrics) GetStats() (int64, int64, time.Duration) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.requestCount == 0 {
		return 0, 0, 0
	}
	return m.requestCount, m.errorCount, m.totalLatency / time.Duration(m.requestCount)
}

func main() {
	metrics := NewMetrics()

	metrics.RecordRequest(150*time.Millisecond, false)
	metrics.RecordRequest(200*time.Millisecond, true)
	metrics.RecordRequest(100*time.Millisecond, false)

	total, errors, avgLatency := metrics.GetStats()
	log.Printf("Total requests: %d, Errors: %d, Average latency: %v", total, errors, avgLatency)
}