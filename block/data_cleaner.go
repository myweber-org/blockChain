
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    ID    int
    Name  string
    Email string
    Score float64
}

func cleanCSV(inputPath, outputPath string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    headers := []string{"ID", "Name", "Email", "Score"}
    if err := writer.Write(headers); err != nil {
        return fmt.Errorf("failed to write headers: %w", err)
    }

    lineNumber := 0
    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("warning: line %d: %v\n", lineNumber, err)
            continue
        }

        if len(row) != 4 {
            fmt.Printf("warning: line %d: invalid column count %d\n", lineNumber, len(row))
            continue
        }

        record, err := parseRecord(row)
        if err != nil {
            fmt.Printf("warning: line %d: %v\n", lineNumber, err)
            continue
        }

        if !validateRecord(record) {
            fmt.Printf("warning: line %d: validation failed\n", lineNumber)
            continue
        }

        outputRow := []string{
            strconv.Itoa(record.ID),
            strings.TrimSpace(record.Name),
            strings.ToLower(strings.TrimSpace(record.Email)),
            fmt.Sprintf("%.2f", record.Score),
        }

        if err := writer.Write(outputRow); err != nil {
            return fmt.Errorf("failed to write record: %w", err)
        }
    }

    return nil
}

func parseRecord(row []string) (Record, error) {
    var record Record

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, fmt.Errorf("invalid ID: %w", err)
    }
    record.ID = id

    record.Name = row[1]

    record.Email = row[2]

    score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
    if err != nil {
        return record, fmt.Errorf("invalid score: %w", err)
    }
    record.Score = score

    return record, nil
}

func validateRecord(r Record) bool {
    if r.ID <= 0 {
        return false
    }
    if strings.TrimSpace(r.Name) == "" {
        return false
    }
    if !strings.Contains(r.Email, "@") {
        return false
    }
    if r.Score < 0 || r.Score > 100 {
        return false
    }
    return true
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := cleanCSV(inputFile, outputFile); err != nil {
        fmt.Printf("error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("data cleaning completed successfully")
}