
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

		if record.ID == "" || record.Name == "" {
			return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
	var active []DataRecord
	for _, record := range records {
		if strings.ToLower(record.Active) == "true" || record.Active == "1" {
			active = append(active, record)
		}
	}
	return active
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	activeRecords := FilterActiveRecords(records)
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Active records: %d\n", len(activeRecords))

	for i, record := range activeRecords {
		if !ValidateEmail(record.Email) {
			fmt.Printf("Warning: Invalid email for record %d (ID: %s)\n", i+1, record.ID)
		}
	}
}package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Value string
	Valid bool
}

func ProcessRecords(records []DataRecord) ([]string, error) {
	if len(records) == 0 {
		return nil, errors.New("no records to process")
	}

	var processed []string
	for _, record := range records {
		if !record.Valid {
			continue
		}

		trimmed := strings.TrimSpace(record.Value)
		if trimmed == "" {
			continue
		}

		processed = append(processed, fmt.Sprintf("ID:%d|Value:%s", record.ID, trimmed))
	}

	if len(processed) == 0 {
		return nil, errors.New("no valid records found")
	}

	return processed, nil
}

func ValidateRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("invalid ID")
	}
	if strings.TrimSpace(record.Value) == "" {
		return errors.New("empty value")
	}
	return nil
}

func main() {
	records := []DataRecord{
		{ID: 1, Value: "alpha", Valid: true},
		{ID: 2, Value: "  ", Valid: true},
		{ID: 3, Value: "beta", Valid: false},
		{ID: 4, Value: "gamma", Valid: true},
	}

	result, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Println("Processed records:")
	for _, item := range result {
		fmt.Println(item)
	}
}