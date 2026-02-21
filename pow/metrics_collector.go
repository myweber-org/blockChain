package main

import (
	"fmt"
	"runtime"
	"time"
)

type SystemMetrics struct {
	Timestamp    time.Time
	CPUPercent   float64
	MemoryAlloc  uint64
	MemoryTotal  uint64
	GoroutineCnt int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		Timestamp:    time.Now(),
		CPUPercent:   getCPUUsage(),
		MemoryAlloc:  m.Alloc,
		MemoryTotal:  m.TotalAlloc,
		GoroutineCnt: runtime.NumGoroutine(),
	}
}

func getCPUUsage() float64 {
	start := time.Now()
	runtime.Gosched()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start).Seconds()
	return (100.0 - (elapsed * 1000 / 100)) / 100.0
}

func printMetrics(metrics SystemMetrics) {
	fmt.Printf("Time: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUPercent)
	fmt.Printf("Memory Allocated: %v bytes\n", metrics.MemoryAlloc)
	fmt.Printf("Total Memory: %v bytes\n", metrics.MemoryTotal)
	fmt.Printf("Goroutines: %d\n", metrics.GoroutineCnt)
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
}