package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
	UniqueErrors []string
}

func parseLogLine(line string) (*LogEntry, error) {
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid log format")
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", matches[1])
	if err != nil {
		return nil, err
	}

	return &LogEntry{
		Timestamp: timestamp,
		Level:     matches[2],
		Message:   matches[3],
	}, nil
}

func analyzeLogs(filePath string) (*LogSummary, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	summary := &LogSummary{}
	errorSet := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}

		summary.TotalEntries++
		switch strings.ToUpper(entry.Level) {
		case "ERROR":
			summary.ErrorCount++
			if !errorSet[entry.Message] {
				errorSet[entry.Message] = true
				summary.UniqueErrors = append(summary.UniqueErrors, entry.Message)
			}
		case "WARN":
			summary.WarnCount++
		case "INFO":
			summary.InfoCount++
		}
	}

	return summary, scanner.Err()
}

func printSummary(summary *LogSummary) {
	fmt.Println("=== Log Analysis Summary ===")
	fmt.Printf("Total entries: %d\n", summary.TotalEntries)
	fmt.Printf("Info level: %d\n", summary.InfoCount)
	fmt.Printf("Warning level: %d\n", summary.WarnCount)
	fmt.Printf("Error level: %d\n", summary.ErrorCount)
	fmt.Printf("Unique errors: %d\n", len(summary.UniqueErrors))

	if len(summary.UniqueErrors) > 0 {
		fmt.Println("\nUnique error messages:")
		for i, err := range summary.UniqueErrors {
			fmt.Printf("%d. %s\n", i+1, err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_analyzer <log_file_path>")
		os.Exit(1)
	}

	summary, err := analyzeLogs(os.Args[1])
	if err != nil {
		fmt.Printf("Error analyzing logs: %v\n", err)
		os.Exit(1)
	}

	printSummary(summary)
}