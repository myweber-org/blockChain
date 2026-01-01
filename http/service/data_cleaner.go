
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

func cleanEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}

func validateScore(score float64) bool {
    return score >= 0 && score <= 100
}

func parseCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNum := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        lineNum++
        if lineNum == 1 {
            continue
        }

        if len(line) != 4 {
            continue
        }

        id, err := strconv.Atoi(strings.TrimSpace(line[0]))
        if err != nil {
            continue
        }

        name := strings.TrimSpace(line[1])
        email := cleanEmail(line[2])

        score, err := strconv.ParseFloat(strings.TrimSpace(line[3]), 64)
        if err != nil || !validateScore(score) {
            continue
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Email: email,
            Score: score,
        })
    }

    return records, nil
}

func calculateAverageScore(records []Record) float64 {
    if len(records) == 0 {
        return 0
    }

    total := 0.0
    for _, record := range records {
        total += record.Score
    }
    return total / float64(len(records))
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_cleaner <csv_file>")
        return
    }

    records, err := parseCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    fmt.Printf("Processed %d valid records\n", len(records))
    fmt.Printf("Average score: %.2f\n", calculateAverageScore(records))

    for i, record := range records {
        if i < 3 {
            fmt.Printf("Sample record: %+v\n", record)
        }
    }
}