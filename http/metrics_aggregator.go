package metrics

import (
	"sync"
	"time"
)

type Aggregator struct {
	mu            sync.RWMutex
	latencies     []time.Duration
	errorCount    int
	totalRequests int
	windowSize    int
}

func NewAggregator(windowSize int) *Aggregator {
	return &Aggregator{
		latencies:  make([]time.Duration, 0, windowSize),
		windowSize: windowSize,
	}
}

func (a *Aggregator) Record(latency time.Duration, isError bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.totalRequests++
	if isError {
		a.errorCount++
	}

	if len(a.latencies) >= a.windowSize {
		a.latencies = a.latencies[1:]
	}
	a.latencies = append(a.latencies, latency)
}

func (a *Aggregator) GetStats() (avgLatency time.Duration, errorRate float64) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.latencies) == 0 {
		return 0, 0
	}

	var total time.Duration
	for _, lat := range a.latencies {
		total += lat
	}
	avgLatency = total / time.Duration(len(a.latencies))

	if a.totalRequests > 0 {
		errorRate = float64(a.errorCount) / float64(a.totalRequests)
	}
	return avgLatency, errorRate
}

func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.latencies = make([]time.Duration, 0, a.windowSize)
	a.errorCount = 0
	a.totalRequests = 0
}package aggregator

import (
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
	subscribers []chan AggregatedResult
}

type AggregatedResult struct {
	Average   float64
	Max       float64
	Min       float64
	Count     int
	Timestamp time.Time
}

func NewSlidingWindowAggregator(windowSize time.Duration) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize:  windowSize,
		metrics:     make([]Metric, 0),
		subscribers: make([]chan AggregatedResult, 0),
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.metrics = append(swa.metrics, Metric{Timestamp: now, Value: value})
	swa.cleanupOldMetrics(now)
	swa.notifySubscribers(now)
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

func (swa *SlidingWindowAggregator) notifySubscribers(currentTime time.Time) {
	if len(swa.metrics) == 0 {
		return
	}

	result := swa.calculateAggregates(currentTime)
	for _, ch := range swa.subscribers {
		select {
		case ch <- result:
		default:
		}
	}
}

func (swa *SlidingWindowAggregator) calculateAggregates(currentTime time.Time) AggregatedResult {
	var sum float64
	max := swa.metrics[0].Value
	min := swa.metrics[0].Value

	for _, metric := range swa.metrics {
		sum += metric.Value
		if metric.Value > max {
			max = metric.Value
		}
		if metric.Value < min {
			min = metric.Value
		}
	}

	return AggregatedResult{
		Average:   sum / float64(len(swa.metrics)),
		Max:       max,
		Min:       min,
		Count:     len(swa.metrics),
		Timestamp: currentTime,
	}
}

func (swa *SlidingWindowAggregator) Subscribe() <-chan AggregatedResult {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	ch := make(chan AggregatedResult, 10)
	swa.subscribers = append(swa.subscribers, ch)
	return ch
}

func (swa *SlidingWindowAggregator) GetCurrentAggregates() AggregatedResult {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	return swa.calculateAggregates(time.Now())
}