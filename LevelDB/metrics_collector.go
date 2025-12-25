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
}