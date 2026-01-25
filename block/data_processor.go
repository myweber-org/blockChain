
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
	ID      int
	Name    string
	Value   float64
	IsValid bool
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

		if len(row) < 4 {
			continue
		}

		record, err := parseRow(row)
		if err != nil {
			fmt.Printf("warning: line %d: %v\n", lineNumber, err)
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
		return record, fmt.Errorf("invalid ID format: %s", row[0])
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value format: %s", row[2])
	}
	record.Value = value

	isValid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid boolean format: %s", row[3])
	}
	record.IsValid = isValid

	return record, nil
}

func ValidateRecords(records []DataRecord) (valid []DataRecord, invalid []DataRecord) {
	for _, record := range records {
		if record.IsValid && record.Value >= 0 {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}
	return valid, invalid
}

func CalculateStatistics(records []DataRecord) (sum float64, average float64, count int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	for _, record := range records {
		if record.IsValid {
			sum += record.Value
			count++
		}
	}

	if count > 0 {
		average = sum / float64(count)
	}

	return sum, average, count
}