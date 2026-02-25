
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

func parseLogLine(line string) *LogEntry {
	pattern := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`)
	matches := pattern.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}
	return &LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
	}
}

func filterLogsByLevel(entries []LogEntry, level string) []LogEntry {
	var filtered []LogEntry
	for _, entry := range entries {
		if strings.EqualFold(entry.Level, level) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
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
	for scanner.Scan() {
		entry := parseLogLine(scanner.Text())
		if entry != nil {
			entries = append(entries, *entry)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	errorLogs := filterLogsByLevel(entries, "ERROR")
	fmt.Printf("Found %d error entries:\n", len(errorLogs))
	for _, entry := range errorLogs {
		fmt.Printf("[%s] %s\n", entry.Timestamp, entry.Message)
	}
}
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
        entry, ok := parseLogLine(scanner.Text())
        if ok && strings.ToUpper(entry.Level) == "ERROR" {
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
        os.Exit(1)
    }

    errors := extractErrors(os.Args[1])
    fmt.Printf("Found %d error entries:\n", len(errors))
    for _, entry := range errors {
        fmt.Printf("[%s] %s: %s\n", entry.Timestamp, entry.Level, entry.Message)
    }
}