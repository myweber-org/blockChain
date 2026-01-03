package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

func parseLogLine(line string) (LogEntry, bool) {
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 4 {
		return LogEntry{}, false
	}

	return LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
	}, true
}

func extractErrors(logPath string) []LogEntry {
	file, err := os.Open(logPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	var errors []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entry, valid := parseLogLine(scanner.Text())
		if valid && strings.ToUpper(entry.Level) == "ERROR" {
			errors = append(errors, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	return errors
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile>")
		return
	}

	errors := extractErrors(os.Args[1])
	fmt.Printf("Found %d error entries:\n", len(errors))
	for i, entry := range errors {
		fmt.Printf("%d. [%s] %s - %s\n", i+1, entry.Timestamp, entry.Level, entry.Message)
	}
}