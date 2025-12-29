package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUPercent  float64
    MemAlloc    uint64
    MemTotal    uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:  time.Now(),
        CPUPercent: getCPUUsage(),
        MemAlloc:   m.Alloc,
        MemTotal:   m.Sys,
        Goroutines: runtime.NumGoroutine(),
    }
}

func getCPUUsage() float64 {
    start := time.Now()
    startGoroutines := runtime.NumGoroutine()

    time.Sleep(100 * time.Millisecond)

    elapsed := time.Since(start).Seconds()
    endGoroutines := runtime.NumGoroutine()

    usage := (float64(endGoroutines-startGoroutines) * 0.1) + (elapsed * 0.5)
    if usage > 100.0 {
        usage = 100.0
    }
    if usage < 0.0 {
        usage = 0.0
    }
    return usage
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("[%s] CPU: %.2f%% | Memory: %.2f MB (Alloc) / %.2f MB (Total) | Goroutines: %d\n",
        metrics.Timestamp.Format("15:04:05"),
        metrics.CPUPercent,
        float64(metrics.MemAlloc)/1024/1024,
        float64(metrics.MemTotal)/1024/1024,
        metrics.Goroutines,
    )
}

func main() {
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    fmt.Println("Starting system metrics collector...")
    for range ticker.C {
        metrics := collectMetrics()
        printMetrics(metrics)
    }
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

    httpRequestCount = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_request_count_total",
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

    wrappedMux := metricsMiddleware(mux)
    http.ListenAndServe(":8080", wrappedMux)
}