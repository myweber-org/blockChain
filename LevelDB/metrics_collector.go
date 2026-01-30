package main

import (
	"log"
	"net/http"
	"time"
)

type MetricsCollector struct {
	requestCount    int
	errorCount      int
	totalLatency    time.Duration
	statusCodeCount map[int]int
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		statusCodeCount: make(map[int]int),
	}
}

func (mc *MetricsCollector) RecordRequest(statusCode int, latency time.Duration) {
	mc.requestCount++
	mc.totalLatency += latency
	mc.statusCodeCount[statusCode]++

	if statusCode >= 400 {
		mc.errorCount++
	}
}

func (mc *MetricsCollector) GetAverageLatency() time.Duration {
	if mc.requestCount == 0 {
		return 0
	}
	return mc.totalLatency / time.Duration(mc.requestCount)
}

func (mc *MetricsCollector) GetSuccessRate() float64 {
	if mc.requestCount == 0 {
		return 0.0
	}
	return float64(mc.requestCount-mc.errorCount) / float64(mc.requestCount) * 100
}

func (mc *MetricsCollector) GetStatusCodeDistribution() map[int]int {
	distribution := make(map[int]int)
	for code, count := range mc.statusCodeCount {
		distribution[code] = count
	}
	return distribution
}

func main() {
	collector := NewMetricsCollector()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Metrics endpoint"))
		collector.RecordRequest(http.StatusOK, time.Since(start))
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error endpoint"))
		collector.RecordRequest(http.StatusInternalServerError, time.Since(start))
	})

	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.Printf("Average latency: %v", collector.GetAverageLatency())
			log.Printf("Success rate: %.2f%%", collector.GetSuccessRate())
			log.Printf("Status code distribution: %v", collector.GetStatusCodeDistribution())
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}