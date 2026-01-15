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
		Level:     strings.ToUpper(matches[2]),
		Message:   matches[3],
	}, true
}

func extractErrors(logPath string) ([]LogEntry, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var errors []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entry, valid := parseLogLine(scanner.Text())
		if valid && entry.Level == "ERROR" {
			errors = append(errors, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return errors, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile>")
		os.Exit(1)
	}

	errors, err := extractErrors(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing log: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d error entries:\n", len(errors))
	for i, entry := range errors {
		fmt.Printf("%d. [%s] %s: %s\n", i+1, entry.Timestamp, entry.Level, entry.Message)
	}
}package main

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
		Level:     strings.ToUpper(matches[2]),
		Message:   matches[3],
	}, true
}

func filterErrors(entries []LogEntry) []LogEntry {
	var errors []LogEntry
	for _, entry := range entries {
		if entry.Level == "ERROR" {
			errors = append(errors, entry)
		}
	}
	return errors
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		entry, ok := parseLogLine(line)
		if ok {
			entries = append(entries, entry)
		} else {
			fmt.Printf("Warning: Could not parse line %d: %s\n", lineNumber, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	errorEntries := filterErrors(entries)
	fmt.Printf("Total log entries: %d\n", len(entries))
	fmt.Printf("Error entries: %d\n", len(errorEntries))
	fmt.Println("\nError messages:")

	for _, entry := range errorEntries {
		fmt.Printf("[%s] %s\n", entry.Timestamp, entry.Message)
	}
}