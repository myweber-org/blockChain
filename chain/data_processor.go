
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID        int
	Name      string
	Value     float64
	Validated bool
}

func processCSVFile(filename string) ([]DataRecord, error) {
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

		if len(row) < 3 {
			continue
		}

		record, err := parseRow(row)
		if err != nil {
			fmt.Printf("skipping invalid row at line %d: %v\n", lineNumber, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID format: %w", err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value format: %w", err)
	}
	record.Value = value

	record.Validated = validateRecord(record)

	return record, nil
}

func validateRecord(record DataRecord) bool {
	if record.ID <= 0 {
		return false
	}
	if len(record.Name) == 0 {
		return false
	}
	if record.Value < 0 {
		return false
	}
	return true
}

func calculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var validCount int

	for _, record := range records {
		if record.Validated {
			sum += record.Value
			validCount++
		}
	}

	if validCount == 0 {
		return 0, 0, 0
	}

	average := sum / float64(validCount)
	return sum, average, validCount
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("error processing file: %v\n", err)
		os.Exit(1)
	}

	total, average, validCount := calculateStatistics(records)

	fmt.Printf("processed %d records\n", len(records))
	fmt.Printf("valid records: %d\n", validCount)
	fmt.Printf("total value: %.2f\n", total)
	fmt.Printf("average value: %.2f\n", average)
}