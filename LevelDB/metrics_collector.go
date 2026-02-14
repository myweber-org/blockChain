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
}package main

import (
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path", "status"},
    )

    httpRequestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_request_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(httpRequestCount)
}

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(rw, r)

        duration := time.Since(start).Seconds()
        status := http.StatusText(rw.statusCode)

        httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
        httpRequestCount.WithLabelValues(r.Method, r.URL.Path, status).Inc()
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    mux.Handle("/metrics", promhttp.Handler())

    handler := metricsMiddleware(mux)
    http.ListenAndServe(":8080", handler)
}package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    CPUUsage      float64
    MemoryUsageMB float64
    GoroutineCount int
}

func collectMetrics() SystemMetrics {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    return SystemMetrics{
        Timestamp:     time.Now(),
        CPUUsage:      calculateCPUUsage(),
        MemoryUsageMB: float64(memStats.Alloc) / 1024 / 1024,
        GoroutineCount: runtime.NumGoroutine(),
    }
}

func calculateCPUUsage() float64 {
    start := time.Now()
    startGoroutines := runtime.NumGoroutine()

    time.Sleep(100 * time.Millisecond)

    elapsed := time.Since(start).Seconds()
    endGoroutines := runtime.NumGoroutine()

    usage := (float64(endGoroutines-startGoroutines) / elapsed) * 10
    if usage < 0 {
        usage = 0
    }
    return usage
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            fmt.Printf("[%s] CPU: %.2f%%, Memory: %.2fMB, Goroutines: %d\n",
                metrics.Timestamp.Format("15:04:05"),
                metrics.CPUUsage,
                metrics.MemoryUsageMB,
                metrics.GoroutineCount)
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

func (mc *MetricsCollector) RecordRequest(latency time.Duration, err error) {
	mc.requestCount++
	mc.totalLatency += latency
	mc.latencySamples = append(mc.latencySamples, latency)
	
	if err != nil {
		mc.errorCount++
	}
}

func (mc *MetricsCollector) AverageLatency() time.Duration {
	if mc.requestCount == 0 {
		return 0
	}
	return mc.totalLatency / time.Duration(mc.requestCount)
}

func (mc *MetricsCollector) ErrorRate() float64 {
	if mc.requestCount == 0 {
		return 0.0
	}
	return float64(mc.errorCount) / float64(mc.requestCount)
}

func (mc *MetricsCollector) PercentileLatency(p float64) time.Duration {
	if len(mc.latencySamples) == 0 {
		return 0
	}
	
	// Simple implementation - for production use proper percentile calculation
	index := int(float64(len(mc.latencySamples)-1) * p / 100.0)
	if index < 0 {
		index = 0
	}
	if index >= len(mc.latencySamples) {
		index = len(mc.latencySamples) - 1
	}
	return mc.latencySamples[index]
}

func main() {
	collector := NewMetricsCollector()
	
	// Simulate some requests
	for i := 0; i < 100; i++ {
		start := time.Now()
		
		// Simulate HTTP request
		_, err := http.Get("http://example.com")
		latency := time.Since(start)
		
		collector.RecordRequest(latency, err)
		time.Sleep(10 * time.Millisecond)
	}
	
	log.Printf("Requests processed: %d", collector.requestCount)
	log.Printf("Average latency: %v", collector.AverageLatency())
	log.Printf("Error rate: %.2f%%", collector.ErrorRate()*100)
	log.Printf("95th percentile latency: %v", collector.PercentileLatency(95))
}package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUPercent  float64
    MemoryAlloc uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:   time.Now().UTC(),
        CPUPercent:  getCPUUsage(),
        MemoryAlloc: m.Alloc,
        Goroutines:  runtime.NumGoroutine(),
    }
}

func getCPUUsage() float64 {
    // Simplified CPU usage calculation
    // In production, use proper system monitoring libraries
    start := time.Now()
    var count int64
    for i := 0; i < 1000000; i++ {
        count += int64(i)
    }
    elapsed := time.Since(start).Seconds()
    
    // Simulate CPU usage based on processing time
    usage := 1.0 / (elapsed + 0.1) * 10
    if usage > 100.0 {
        usage = 100.0
    }
    return usage
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Metrics collected at: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUPercent)
    fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
    fmt.Printf("Active Goroutines: %d\n", metrics.Goroutines)
    fmt.Println("---")
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    fmt.Println("Starting system metrics collector...")
    
    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            printMetrics(metrics)
        }
    }
}