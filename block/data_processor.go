
package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID      int
    Name    string
    Value   float64
    Active  bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
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

        if len(row) != 4 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
        }

        record, err := parseRow(row, lineNumber)
        if err != nil {
            return nil, err
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func parseRow(row []string, lineNumber int) (DataRecord, error) {
    var record DataRecord

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
    }
    record.ID = id

    name := strings.TrimSpace(row[1])
    if name == "" {
        return record, fmt.Errorf("empty name at line %d", lineNumber)
    }
    record.Name = name

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
    }
    record.Value = value

    active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return record, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
    }
    record.Active = active

    return record, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, []error) {
    var validRecords []DataRecord
    var validationErrors []error

    for i, record := range records {
        if record.ID <= 0 {
            validationErrors = append(validationErrors, fmt.Errorf("record %d: invalid ID %d", i+1, record.ID))
            continue
        }

        if record.Value < 0 {
            validationErrors = append(validationErrors, fmt.Errorf("record %d: negative value %f", i+1, record.Value))
            continue
        }

        if len(record.Name) > 100 {
            validationErrors = append(validationErrors, fmt.Errorf("record %d: name exceeds 100 characters", i+1))
            continue
        }

        validRecords = append(validRecords, record)
    }

    return validRecords, validationErrors
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var activeCount int
    minValue := records[0].Value
    maxValue := records[0].Value

    for _, record := range records {
        sum += record.Value

        if record.Value < minValue {
            minValue = record.Value
        }
        if record.Value > maxValue {
            maxValue = record.Value
        }

        if record.Active {
            activeCount++
        }
    }

    average := sum / float64(len(records))
    return average, maxValue - minValue, activeCount
}