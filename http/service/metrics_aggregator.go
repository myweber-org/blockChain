
package metrics

import (
	"container/list"
	"sort"
	"sync"
	"time"
)

type DataPoint struct {
	Value     float64
	Timestamp time.Time
}

type SlidingWindowAggregator struct {
	windowSize  time.Duration
	maxPoints   int
	dataPoints  *list.List
	mu          sync.RWMutex
}

func NewSlidingWindowAggregator(windowSize time.Duration, maxPoints int) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize: windowSize,
		maxPoints:  maxPoints,
		dataPoints: list.New(),
	}
}

func (swa *SlidingWindowAggregator) AddPoint(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.dataPoints.PushBack(DataPoint{
		Value:     value,
		Timestamp: now,
	})

	swa.cleanupOldPoints(now)
	if swa.dataPoints.Len() > swa.maxPoints {
		swa.dataPoints.Remove(swa.dataPoints.Front())
	}
}

func (swa *SlidingWindowAggregator) cleanupOldPoints(now time.Time) {
	cutoff := now.Add(-swa.windowSize)
	for e := swa.dataPoints.Front(); e != nil; {
		next := e.Next()
		if dp := e.Value.(DataPoint); dp.Timestamp.Before(cutoff) {
			swa.dataPoints.Remove(e)
		}
		e = next
	}
}

func (swa *SlidingWindowAggregator) CalculatePercentile(p float64) (float64, bool) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if swa.dataPoints.Len() == 0 {
		return 0, false
	}

	values := make([]float64, 0, swa.dataPoints.Len())
	for e := swa.dataPoints.Front(); e != nil; e = e.Next() {
		values = append(values, e.Value.(DataPoint).Value)
	}

	sort.Float64s(values)
	index := int(float64(len(values)-1) * p / 100.0)
	return values[index], true
}

func (swa *SlidingWindowAggregator) GetStats() (min, max, avg float64, count int) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if swa.dataPoints.Len() == 0 {
		return 0, 0, 0, 0
	}

	var sum float64
	min = swa.dataPoints.Front().Value.(DataPoint).Value
	max = min
	count = 0

	for e := swa.dataPoints.Front(); e != nil; e = e.Next() {
		value := e.Value.(DataPoint).Value
		sum += value
		count++
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	avg = sum / float64(count)
	return min, max, avg, count
}