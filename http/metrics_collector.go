package main

import (
	"fmt"
	"runtime"
	"time"
)

type SystemMetrics struct {
	Timestamp   time.Time
	CPUPercent  float64
	MemoryAlloc uint64
	NumGoroutine int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		Timestamp:   time.Now(),
		CPUPercent:  getCPUUsage(),
		MemoryAlloc: m.Alloc,
		NumGoroutine: runtime.NumGoroutine(),
	}
}

func getCPUUsage() float64 {
	start := time.Now()
	runtime.Gosched()
	time.Sleep(50 * time.Millisecond)
	elapsed := time.Since(start)
	return float64(elapsed) / float64(time.Second) * 100
}

func displayMetrics(metrics SystemMetrics) {
	fmt.Printf("Time: %s\n", metrics.Timestamp.Format("15:04:05"))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUPercent)
	fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
	fmt.Printf("Active Goroutines: %d\n", metrics.NumGoroutine)
	fmt.Println("---")
}

func main() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 5; i++ {
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
}

func (mc *MetricsCollector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		mc.requestCount++
		mc.totalLatency += duration
		
		if recorder.statusCode >= 400 {
			mc.errorCount++
		}
		
		log.Printf("Request processed: %s %s - Status: %d, Duration: %v", 
			r.Method, r.URL.Path, recorder.statusCode, duration)
	})
}

func (mc *MetricsCollector) GetMetrics() (float64, float64) {
	if mc.requestCount == 0 {
		return 0, 0
	}
	avgLatency := mc.totalLatency.Seconds() / float64(mc.requestCount)
	errorRate := float64(mc.errorCount) / float64(mc.requestCount) * 100
	return avgLatency, errorRate
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func main() {
	collector := &MetricsCollector{}
	
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		avgLatency, errorRate := collector.GetMetrics()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"avg_latency_seconds":` + formatFloat(avgLatency) + 
			`,"error_rate_percent":` + formatFloat(errorRate) + 
			`,"total_requests":` + formatInt(collector.requestCount) + `}`))
	})
	
	handler := collector.Middleware(mux)
	
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func formatFloat(f float64) string {
	return string(strconv.AppendFloat([]byte{}, f, 'f', 4, 64))
}

func formatInt(i int) string {
	return strconv.Itoa(i)
}