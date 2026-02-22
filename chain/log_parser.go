package main

import (
    "encoding/json"
    "fmt"
    "regexp"
    "strings"
)

type LogEntry struct {
    Timestamp string `json:"timestamp"`
    Level     string `json:"level"`
    Message   string `json:"message"`
    Source    string `json:"source"`
}

type LogParser struct {
    pattern *regexp.Regexp
}

func NewLogParser() *LogParser {
    pattern := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+) - (.+)$`)
    return &LogParser{pattern: pattern}
}

func (p *LogParser) ParseLine(line string) (*LogEntry, error) {
    matches := p.pattern.FindStringSubmatch(strings.TrimSpace(line))
    if matches == nil {
        return nil, fmt.Errorf("invalid log format")
    }
    
    return &LogEntry{
        Timestamp: matches[1],
        Level:     matches[2],
        Source:    matches[3],
        Message:   matches[4],
    }, nil
}

func (p *LogParser) ParseLines(lines []string) []LogEntry {
    var entries []LogEntry
    for _, line := range lines {
        entry, err := p.ParseLine(line)
        if err != nil {
            continue
        }
        entries = append(entries, *entry)
    }
    return entries
}

func (p *LogParser) ToJSON(entries []LogEntry) (string, error) {
    data, err := json.MarshalIndent(entries, "", "  ")
    if err != nil {
        return "", err
    }
    return string(data), nil
}

func main() {
    parser := NewLogParser()
    
    sampleLogs := []string{
        "2024-01-15 10:30:45 [INFO] auth-service - User login successful",
        "2024-01-15 10:31:22 [ERROR] payment-service - Transaction timeout",
        "2024-01-15 10:32:01 [WARN] api-gateway - High latency detected",
    }
    
    entries := parser.ParseLines(sampleLogs)
    
    jsonOutput, err := parser.ToJSON(entries)
    if err != nil {
        fmt.Printf("Error converting to JSON: %v\n", err)
        return
    }
    
    fmt.Println("Parsed log entries:")
    fmt.Println(jsonOutput)
}