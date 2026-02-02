package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUPercent  float64
    MemoryUsage uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:   time.Now(),
        CPUPercent:  getCPUUsage(),
        MemoryUsage: m.Alloc,
        Goroutines:  runtime.NumGoroutine(),
    }
}

func getCPUUsage() float64 {
    start := time.Now()
    startCPU := runtime.NumCgoCall()

    time.Sleep(100 * time.Millisecond)

    elapsed := time.Since(start).Seconds()
    endCPU := runtime.NumCgoCall()

    if elapsed > 0 {
        return float64(endCPU-startCPU) / elapsed * 100
    }
    return 0.0
}

func displayMetrics(metrics SystemMetrics) {
    fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
    fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUPercent)
    fmt.Printf("Memory Usage: %d bytes\n", metrics.MemoryUsage)
    fmt.Printf("Active Goroutines: %d\n", metrics.Goroutines)
    fmt.Println("---")
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            displayMetrics(metrics)
        }
    }
}package main

import (
	"log"
	"net/http"
	"time"
)

type MetricsCollector struct {
	requestCount    int
	errorCount      int
	totalLatency    time.Duration
	latencySamples  []time.Duration
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		latencySamples: make([]time.Duration, 0),
	}
}

func (m *MetricsCollector) RecordRequest(latency time.Duration, isError bool) {
	m.requestCount++
	m.totalLatency += latency
	m.latencySamples = append(m.latencySamples, latency)
	
	if isError {
		m.errorCount++
	}
}

func (m *MetricsCollector) GetAverageLatency() time.Duration {
	if m.requestCount == 0 {
		return 0
	}
	return m.totalLatency / time.Duration(m.requestCount)
}

func (m *MetricsCollector) GetErrorRate() float64 {
	if m.requestCount == 0 {
		return 0.0
	}
	return float64(m.errorCount) / float64(m.requestCount)
}

func (m *MetricsCollector) GetPercentileLatency(percentile float64) time.Duration {
	if len(m.latencySamples) == 0 {
		return 0
	}
	
	index := int(float64(len(m.latencySamples)-1) * percentile / 100.0)
	if index < 0 {
		index = 0
	}
	if index >= len(m.latencySamples) {
		index = len(m.latencySamples) - 1
	}
	
	return m.latencySamples[index]
}

func main() {
	collector := NewMetricsCollector()
	
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status": "ok"}`))
		
		latency := time.Since(start)
		collector.RecordRequest(latency, err != nil)
		
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})
	
	log.Println("Starting metrics server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}