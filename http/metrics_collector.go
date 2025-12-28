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
        
        defer func() {
            duration := time.Since(start).Seconds()
            status := http.StatusText(rw.statusCode)
            
            httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
            httpRequestCount.WithLabelValues(r.Method, r.URL.Path, status).Inc()
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
    mux.Handle("/metrics", promhttp.Handler())
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Hello, World!"))
    })

    handler := metricsMiddleware(mux)
    http.ListenAndServe(":8080", handler)
}package main

import (
	"log"
	"net/http"
	"time"
)

type MetricsCollector struct {
	latencies []time.Duration
	errors    int
	requests  int
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		latencies: make([]time.Duration, 0),
	}
}

func (mc *MetricsCollector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		mc.requests++

		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		mc.latencies = append(mc.latencies, duration)

		if lrw.statusCode >= 400 {
			mc.errors++
		}
	})
}

func (mc *MetricsCollector) AverageLatency() time.Duration {
	if len(mc.latencies) == 0 {
		return 0
	}
	var total time.Duration
	for _, lat := range mc.latencies {
		total += lat
	}
	return total / time.Duration(len(mc.latencies))
}

func (mc *MetricsCollector) ErrorRate() float64 {
	if mc.requests == 0 {
		return 0.0
	}
	return float64(mc.errors) / float64(mc.requests) * 100
}

func (mc *MetricsCollector) Reset() {
	mc.latencies = make([]time.Duration, 0)
	mc.errors = 0
	mc.requests = 0
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func main() {
	collector := NewMetricsCollector()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		avgLatency := collector.AverageLatency()
		errorRate := collector.ErrorRate()
		response := map[string]interface{}{
			"total_requests": collector.requests,
			"errors":         collector.errors,
			"error_rate":     errorRate,
			"avg_latency_ms": avgLatency.Milliseconds(),
		}
		json.NewEncoder(w).Encode(response)
	})

	handler := collector.Middleware(mux)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}