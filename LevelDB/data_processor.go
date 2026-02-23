
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value cannot be negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
	return transformed
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}

func GenerateReport(records []DataRecord) string {
	var builder strings.Builder
	builder.WriteString("Data Processing Report\n")
	builder.WriteString("======================\n")
	
	totalValue := 0.0
	for _, record := range records {
		builder.WriteString(fmt.Sprintf("ID: %s, Value: %.2f, Tags: %v\n", 
			record.ID, record.Value, record.Tags))
		totalValue += record.Value
	}
	builder.WriteString(fmt.Sprintf("\nTotal Processed Value: %.2f\n", totalValue))
	builder.WriteString(fmt.Sprintf("Records Processed: %d\n", len(records)))
	
	return builder.String()
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec001",
			Value:     100.0,
			Timestamp: time.Now().Add(-24 * time.Hour),
			Tags:      []string{"initial"},
		},
		{
			ID:        "rec002",
			Value:     200.0,
			Timestamp: time.Now().Add(-12 * time.Hour),
			Tags:      []string{"initial", "priority"},
		},
	}
	
	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}
	
	report := GenerateReport(processed)
	fmt.Println(report)
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []Record) []Record {
	var validRecords []Record
	for _, r := range records {
		if r.ID > 0 && r.Name != "" && r.Value >= 0 {
			validRecords = append(validRecords, r)
		}
	}
	return validRecords
}

func CalculateTotal(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		return
	}

	records, err := ProcessCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	validRecords := ValidateRecords(records)
	total := CalculateTotal(validRecords)

	fmt.Printf("Processed %d records, %d valid\n", len(records), len(validRecords))
	fmt.Printf("Total value: %.2f\n", total)
}