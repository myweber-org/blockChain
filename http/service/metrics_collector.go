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

	// Simulate CPU usage calculation
	metrics.CPUUsage = calculateCPUUsage()

	return metrics
}

func calculateCPUUsage() float64 {
	// Placeholder for actual CPU calculation logic
	// In production, this would use proper system monitoring
	return float64(runtime.NumCPU()) * 0.75
}

func printMetrics(metrics SystemMetrics) {
	fmt.Printf("Metrics collected at: %s\n", metrics.Timestamp.Format(time.RFC3339))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUUsage)
	fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
	fmt.Printf("Total Memory: %d bytes\n", metrics.MemoryTotal)
	fmt.Printf("Active Goroutines: %d\n", metrics.GoroutineCount)
	fmt.Println("---")
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 3; i++ {
		select {
		case <-ticker.C:
			metrics := collectMetrics()
			printMetrics(metrics)
		}
	}
}