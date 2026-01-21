package main

import (
	"fmt"
	"regexp"
	"strings"
)

type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

var logPatterns = map[string]*regexp.Regexp{
	"timestamp": regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`),
	"level":     regexp.MustCompile(`(INFO|WARN|ERROR|DEBUG)`),
}

func ParseLogLine(line string) (*LogEntry, error) {
	tsMatch := logPatterns["timestamp"].FindString(line)
	if tsMatch == "" {
		return nil, fmt.Errorf("no timestamp found")
	}

	levelMatch := logPatterns["level"].FindString(line)
	if levelMatch == "" {
		levelMatch = "UNKNOWN"
	}

	message := strings.TrimSpace(line)
	for _, pattern := range []string{tsMatch, levelMatch} {
		message = strings.Replace(message, pattern, "", 1)
	}
	message = strings.TrimSpace(message)

	return &LogEntry{
		Timestamp: tsMatch,
		Level:     levelMatch,
		Message:   message,
	}, nil
}

func main() {
	sampleLog := "2023-10-05 14:30:25 INFO User login successful from 192.168.1.10"
	entry, err := ParseLogLine(sampleLog)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	fmt.Printf("Time: %s\nLevel: %s\nMessage: %s\n",
		entry.Timestamp, entry.Level, entry.Message)
}