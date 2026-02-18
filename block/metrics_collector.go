
package main

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
            Name: "http_request_count_total",
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
        defer func() {
            duration := time.Since(start).Seconds()
            httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.statusCode)).Observe(duration)
            httpRequestCount.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.statusCode)).Inc()
        }()
        next.ServeHTTP(rw, r)
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
        w.WriteHeader(http.StatusOK)
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
    Timestamp   time.Time
    CPUPercent  float64
    MemoryAlloc uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return SystemMetrics{
        Timestamp:   time.Now(),
        MemoryAlloc: m.Alloc,
        Goroutines:  runtime.NumGoroutine(),
    }
}

func startMetricsCollector(interval time.Duration, stopChan <-chan struct{}) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            fmt.Printf("[%s] Memory: %v bytes, Goroutines: %d\n",
                metrics.Timestamp.Format(time.RFC3339),
                metrics.MemoryAlloc,
                metrics.Goroutines)
        case <-stopChan:
            fmt.Println("Metrics collector stopped")
            return
        }
    }
}

func main() {
    stopChan := make(chan struct{})
    go startMetricsCollector(2*time.Second, stopChan)
    
    time.Sleep(10 * time.Second)
    close(stopChan)
    time.Sleep(1 * time.Second)
}package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	httpRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start).Seconds()
		status := http.StatusText(rw.statusCode)
		
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
		httpRequestTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
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
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", metricsMiddleware(handler))
	
	http.ListenAndServe(":8080", nil)
}
package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    Goroutines  int
    AllocBytes  uint64
    TotalAlloc  uint64
    SysBytes    uint64
    NumGC       uint32
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return SystemMetrics{
        Timestamp:   time.Now(),
        Goroutines:  runtime.NumGoroutine(),
        AllocBytes:  m.Alloc,
        TotalAlloc:  m.TotalAlloc,
        SysBytes:    m.Sys,
        NumGC:       m.NumGC,
    }
}

func logMetrics(metrics SystemMetrics) {
    fmt.Printf("[%s] Goroutines: %d, Alloc: %.2f MB, Sys: %.2f MB, GC Cycles: %d\n",
        metrics.Timestamp.Format("15:04:05"),
        metrics.Goroutines,
        float64(metrics.AllocBytes)/1024/1024,
        float64(metrics.SysBytes)/1024/1024,
        metrics.NumGC)
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            logMetrics(metrics)
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

func (mc *MetricsCollector) RecordRequest(latency time.Duration, isError bool) {
	mc.requestCount++
	mc.totalLatency += latency
	mc.latencySamples = append(mc.latencySamples, latency)
	if isError {
		mc.errorCount++
	}
}

func (mc *MetricsCollector) GetAverageLatency() time.Duration {
	if mc.requestCount == 0 {
		return 0
	}
	return mc.totalLatency / time.Duration(mc.requestCount)
}

func (mc *MetricsCollector) GetErrorRate() float64 {
	if mc.requestCount == 0 {
		return 0.0
	}
	return float64(mc.errorCount) / float64(mc.requestCount)
}

func (mc *MetricsCollector) GetPercentile(p float64) time.Duration {
	if len(mc.latencySamples) == 0 {
		return 0
	}
	// Simplified percentile calculation
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

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			latency := time.Since(start)
			collector.RecordRequest(latency, false)
		}()

		avgLatency := collector.GetAverageLatency()
		errorRate := collector.GetErrorRate()
		p95 := collector.GetPercentile(95)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"request_count": ` + string(rune(collector.requestCount)) + `,
			"average_latency_ms": ` + string(rune(avgLatency.Milliseconds())) + `,
			"error_rate": ` + string(rune(errorRate)) + `,
			"p95_latency_ms": ` + string(rune(p95.Milliseconds())) + `
		}`))
	})

	log.Println("Starting metrics server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}