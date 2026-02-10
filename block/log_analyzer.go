package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

type LogSummary struct {
	TotalEntries int
	LevelCounts  map[string]int
	StartTime    time.Time
	EndTime      time.Time
}

func parseLogLine(line string) (LogEntry, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return LogEntry{}, fmt.Errorf("invalid log format")
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05Z", parts[0])
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		Timestamp: timestamp,
		Level:     parts[1],
		Message:   parts[2],
	}, nil
}

func analyzeLogs(filePath string) (LogSummary, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return LogSummary{}, err
	}
	defer file.Close()

	summary := LogSummary{
		LevelCounts: make(map[string]int),
	}
	firstEntry := true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}

		summary.TotalEntries++
		summary.LevelCounts[entry.Level]++

		if firstEntry {
			summary.StartTime = entry.Timestamp
			summary.EndTime = entry.Timestamp
			firstEntry = false
		} else {
			if entry.Timestamp.Before(summary.StartTime) {
				summary.StartTime = entry.Timestamp
			}
			if entry.Timestamp.After(summary.EndTime) {
				summary.EndTime = entry.Timestamp
			}
		}
	}

	return summary, scanner.Err()
}

func printSummary(summary LogSummary) {
	fmt.Println("Log Analysis Summary")
	fmt.Println("====================")
	fmt.Printf("Total entries: %d\n", summary.TotalEntries)
	fmt.Printf("Time range: %s to %s\n", summary.StartTime.Format(time.RFC3339), summary.EndTime.Format(time.RFC3339))
	fmt.Println("\nLevel distribution:")
	for level, count := range summary.LevelCounts {
		percentage := float64(count) / float64(summary.TotalEntries) * 100
		fmt.Printf("  %s: %d (%.1f%%)\n", level, count, percentage)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_analyzer <logfile>")
		os.Exit(1)
	}

	summary, err := analyzeLogs(os.Args[1])
	if err != nil {
		fmt.Printf("Error analyzing logs: %v\n", err)
		os.Exit(1)
	}

	printSummary(summary)
}package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"
)

type LogEntry struct {
    Timestamp time.Time
    Level     string
    Message   string
}

type LogSummary struct {
    TotalEntries int
    ErrorCount   int
    WarnCount    int
    InfoCount    int
    StartTime    time.Time
    EndTime      time.Time
}

func parseLogLine(line string) (LogEntry, error) {
    parts := strings.SplitN(line, " ", 3)
    if len(parts) < 3 {
        return LogEntry{}, fmt.Errorf("invalid log format")
    }

    timestamp, err := time.Parse("2006-01-02T15:04:05Z", parts[0])
    if err != nil {
        return LogEntry{}, err
    }

    return LogEntry{
        Timestamp: timestamp,
        Level:     parts[1],
        Message:   parts[2],
    }, nil
}

func analyzeLogs(filePath string) (LogSummary, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return LogSummary{}, err
    }
    defer file.Close()

    var summary LogSummary
    scanner := bufio.NewScanner(file)
    firstEntry := true

    for scanner.Scan() {
        entry, err := parseLogLine(scanner.Text())
        if err != nil {
            continue
        }

        summary.TotalEntries++
        switch entry.Level {
        case "ERROR":
            summary.ErrorCount++
        case "WARN":
            summary.WarnCount++
        case "INFO":
            summary.InfoCount++
        }

        if firstEntry {
            summary.StartTime = entry.Timestamp
            firstEntry = false
        }
        summary.EndTime = entry.Timestamp
    }

    return summary, scanner.Err()
}

func printSummary(summary LogSummary) {
    fmt.Println("Log Analysis Summary")
    fmt.Println("====================")
    fmt.Printf("Total entries: %d\n", summary.TotalEntries)
    fmt.Printf("Error count: %d\n", summary.ErrorCount)
    fmt.Printf("Warning count: %d\n", summary.WarnCount)
    fmt.Printf("Info count: %d\n", summary.InfoCount)
    fmt.Printf("Time range: %s to %s\n",
        summary.StartTime.Format("2006-01-02 15:04:05"),
        summary.EndTime.Format("2006-01-02 15:04:05"))
    fmt.Printf("Duration: %v\n", summary.EndTime.Sub(summary.StartTime))
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_analyzer <logfile>")
        os.Exit(1)
    }

    summary, err := analyzeLogs(os.Args[1])
    if err != nil {
        fmt.Printf("Error analyzing logs: %v\n", err)
        os.Exit(1)
    }

    printSummary(summary)
}