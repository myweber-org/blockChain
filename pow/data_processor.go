
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID      string
    Name    string
    Email   string
    Active  string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []DataRecord
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if lineNumber == 1 {
            continue
        }

        if len(row) < 4 {
            return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
        }

        record := DataRecord{
            ID:     strings.TrimSpace(row[0]),
            Name:   strings.TrimSpace(row[1]),
            Email:  strings.TrimSpace(row[2]),
            Active: strings.TrimSpace(row[3]),
        }

        if !isValidRecord(record) {
            return nil, fmt.Errorf("invalid data at line %d", lineNumber)
        }

        records = append(records, record)
    }

    return records, nil
}

func isValidRecord(record DataRecord) bool {
    if record.ID == "" || record.Name == "" || record.Email == "" {
        return false
    }
    if record.Active != "true" && record.Active != "false" {
        return false
    }
    return true
}

func FilterActiveUsers(records []DataRecord) []DataRecord {
    var activeUsers []DataRecord
    for _, record := range records {
        if record.Active == "true" {
            activeUsers = append(activeUsers, record)
        }
    }
    return activeUsers
}

func GenerateReport(records []DataRecord) {
    fmt.Printf("Total records processed: %d\n", len(records))
    activeUsers := FilterActiveUsers(records)
    fmt.Printf("Active users: %d\n", len(activeUsers))
    fmt.Printf("Inactive users: %d\n", len(records)-len(activeUsers))
    
    fmt.Println("\nSample records:")
    displayCount := 3
    if len(records) < displayCount {
        displayCount = len(records)
    }
    for i := 0; i < displayCount; i++ {
        fmt.Printf("  ID: %s, Name: %s, Email: %s\n", 
            records[i].ID, records[i].Name, records[i].Email)
    }
}