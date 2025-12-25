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
}