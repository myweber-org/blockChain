
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
        w.Write([]byte("Hello, metrics!"))
    })
    mux.Handle("/metrics", promhttp.Handler())

    wrappedMux := metricsMiddleware(mux)
    http.ListenAndServe(":8080", wrappedMux)
}package main

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
	runtime.Gosched()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)
	return float64(elapsed) / float64(time.Second) * 100
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

	for i := 0; i < 3; i++ {
		metrics := collectMetrics()
		displayMetrics(metrics)
		<-ticker.C
	}
}package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUUsage    float64
    MemoryUsage uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    return SystemMetrics{
        Timestamp:   time.Now(),
        CPUUsage:    calculateCPUUsage(),
        MemoryUsage: memStats.Alloc,
        Goroutines:  runtime.NumGoroutine(),
    }
}

func calculateCPUUsage() float64 {
    start := time.Now()
    runtime.Gosched()
    time.Sleep(50 * time.Millisecond)
    elapsed := time.Since(start)

    return float64(elapsed) / float64(time.Second) * 100
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUUsage)
    fmt.Printf("Memory Usage: %d bytes\n", metrics.MemoryUsage)
    fmt.Printf("Active Goroutines: %d\n", metrics.Goroutines)
    fmt.Println("---")
}

func main() {
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        metrics := collectMetrics()
        printMetrics(metrics)
    }
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

    httpRequestTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(httpRequestTotal)
}

func instrumentedHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        defer func() {
            duration := time.Since(start).Seconds()
            status := http.StatusText(rw.statusCode)
            
            httpRequestDuration.WithLabelValues(
                r.Method,
                r.URL.Path,
                status,
            ).Observe(duration)
            
            httpRequestTotal.WithLabelValues(
                r.Method,
                r.URL.Path,
                status,
            ).Inc()
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

    server := &http.Server{
        Addr:    ":8080",
        Handler: instrumentedHandler(mux),
    }

    server.ListenAndServe()
}package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    GoroutineCount int
    MemoryAlloc   uint64
    MemoryTotal   uint64
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return SystemMetrics{
        Timestamp:     time.Now().UTC(),
        GoroutineCount: runtime.NumGoroutine(),
        MemoryAlloc:   m.Alloc,
        MemoryTotal:   m.TotalAlloc,
    }
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("Goroutines: %d\n", metrics.GoroutineCount)
    fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
    fmt.Printf("Total Memory: %d bytes\n", metrics.MemoryTotal)
    fmt.Println("---")
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            printMetrics(metrics)
        }
    }
}package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    CPUPercent    float64
    MemoryUsedMB  uint64
    MemoryTotalMB uint64
    GoroutineCount int
}

func collectMetrics() SystemMetrics {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    return SystemMetrics{
        Timestamp:     time.Now(),
        MemoryUsedMB:  memStats.Alloc / 1024 / 1024,
        MemoryTotalMB: memStats.Sys / 1024 / 1024,
        GoroutineCount: runtime.NumGoroutine(),
    }
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("[%s] Goroutines: %d | Memory: %dMB/%dMB\n",
        metrics.Timestamp.Format("15:04:05"),
        metrics.GoroutineCount,
        metrics.MemoryUsedMB,
        metrics.MemoryTotalMB)
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            printMetrics(metrics)
        }
    }
}