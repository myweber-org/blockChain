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
}

func (m *MetricsCollector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		m.requestCount++
		m.totalLatency += duration
		
		if recorder.statusCode >= 400 {
			m.errorCount++
		}
		
		log.Printf("Request processed: %s %s - Status: %d, Duration: %v", 
			r.Method, r.URL.Path, recorder.statusCode, duration)
	})
}

func (m *MetricsCollector) GetMetrics() map[string]interface{} {
	avgLatency := time.Duration(0)
	if m.requestCount > 0 {
		avgLatency = m.totalLatency / time.Duration(m.requestCount)
	}
	
	errorRate := 0.0
	if m.requestCount > 0 {
		errorRate = float64(m.errorCount) / float64(m.requestCount) * 100
	}
	
	return map[string]interface{}{
		"total_requests":   m.requestCount,
		"error_count":      m.errorCount,
		"error_rate":       errorRate,
		"avg_latency_ns":   avgLatency.Nanoseconds(),
		"avg_latency_str":  avgLatency.String(),
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func main() {
	collector := &MetricsCollector{}
	
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := collector.GetMetrics()
		w.Header().Set("Content-Type", "application/json")
		// In real implementation, use json.Marshal
		w.Write([]byte("Metrics endpoint - implement JSON serialization"))
	})
	
	handler := collector.Middleware(mux)
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}