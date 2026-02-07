package main

import (
	"fmt"
	"runtime"
	"time"
)

type SystemMetrics struct {
	Timestamp    time.Time
	CPUUsage     float64
	MemoryAlloc  uint64
	MemoryTotal  uint64
	GoroutineCount int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := SystemMetrics{
		Timestamp:     time.Now(),
		MemoryAlloc:   m.Alloc,
		MemoryTotal:   m.TotalAlloc,
		GoroutineCount: runtime.NumGoroutine(),
	}

	metrics.CPUUsage = calculateCPUUsage()
	return metrics
}

func calculateCPUUsage() float64 {
	start := time.Now()
	runtime.Gosched()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start).Seconds()
	return (100.0 - (elapsed * 1000)) / 100.0
}

func displayMetrics(metrics SystemMetrics) {
	fmt.Printf("Time: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUUsage)
	fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
	fmt.Printf("Total Memory: %d bytes\n", metrics.MemoryTotal)
	fmt.Printf("Active Goroutines: %d\n", metrics.GoroutineCount)
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
    "fmt"
    "net/http"
    "time"
)

var (
    requestCount    = make(map[string]int)
    requestDuration = make(map[string]time.Duration)
    statusCodes     = make(map[int]int)
)

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        path := r.URL.Path

        recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(recorder, r)

        duration := time.Since(start)
        requestCount[path]++
        requestDuration[path] += duration
        statusCodes[recorder.statusCode]++
    })
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
    rr.statusCode = code
    rr.ResponseWriter.WriteHeader(code)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintln(w, "HTTP Request Metrics")
    fmt.Fprintln(w, "====================")

    fmt.Fprintln(w, "\nRequest Count by Path:")
    for path, count := range requestCount {
        fmt.Fprintf(w, "  %s: %d\n", path, count)
    }

    fmt.Fprintln(w, "\nAverage Duration by Path:")
    for path, totalDuration := range requestDuration {
        count := requestCount[path]
        avg := totalDuration / time.Duration(count)
        fmt.Fprintf(w, "  %s: %v\n", path, avg)
    }

    fmt.Fprintln(w, "\nStatus Code Distribution:")
    for code, count := range statusCodes {
        fmt.Fprintf(w, "  %d: %d\n", code, count)
    }
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    time.Sleep(time.Duration(50+time.Now().UnixNano()%100) * time.Millisecond)
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Request processed successfully")
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", mainHandler)
    mux.HandleFunc("/metrics", metricsHandler)

    wrappedMux := metricsMiddleware(mux)

    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", wrappedMux)
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.Write([]byte("Hello, metrics!"))
	})
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    ":8080",
		Handler: instrumentedHandler(mux),
	}
	server.ListenAndServe()
}