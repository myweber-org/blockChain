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
}