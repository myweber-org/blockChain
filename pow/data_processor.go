
package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Value   int    `json:"value"`
    Active  bool   `json:"active"`
}

func parseCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []Record{}
    lineNumber := 0

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if lineNumber == 0 {
            lineNumber++
            continue
        }

        id, _ := strconv.Atoi(row[0])
        value, _ := strconv.Atoi(row[2])
        active, _ := strconv.ParseBool(row[3])

        record := Record{
            ID:     id,
            Name:   row[1],
            Value:  value,
            Active: active,
        }
        records = append(records, record)
        lineNumber++
    }
    return records, nil
}

func convertToJSON(records []Record) (string, error) {
    jsonData, err := json.MarshalIndent(records, "", "  ")
    if err != nil {
        return "", err
    }
    return string(jsonData), nil
}

func filterActiveRecords(records []Record) []Record {
    var active []Record
    for _, record := range records {
        if record.Active {
            active = append(active, record)
        }
    }
    return active
}

func calculateTotalValue(records []Record) int {
    total := 0
    for _, record := range records {
        total += record.Value
    }
    return total
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        return
    }

    records, err := parseCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error parsing CSV: %v\n", err)
        return
    }

    fmt.Printf("Total records: %d\n", len(records))
    
    activeRecords := filterActiveRecords(records)
    fmt.Printf("Active records: %d\n", len(activeRecords))
    
    totalValue := calculateTotalValue(records)
    fmt.Printf("Total value: %d\n", totalValue)

    jsonOutput, err := convertToJSON(records)
    if err != nil {
        fmt.Printf("Error converting to JSON: %v\n", err)
        return
    }

    outputFile := "output.json"
    err = os.WriteFile(outputFile, []byte(jsonOutput), 0644)
    if err != nil {
        fmt.Printf("Error writing JSON file: %v\n", err)
        return
    }
    
    fmt.Printf("JSON output written to %s\n", outputFile)
}