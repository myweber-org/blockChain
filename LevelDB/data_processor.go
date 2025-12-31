
package main

import (
	"strings"
	"unicode"
)

// CleanInput removes extra whitespace and normalizes line endings
func CleanInput(input string) string {
	// Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	var result strings.Builder
	result.Grow(len(trimmed))
	spaceFlag := false
	
	for _, r := range trimmed {
		if unicode.IsSpace(r) {
			if !spaceFlag {
				result.WriteRune(' ')
				spaceFlag = true
			}
		} else {
			result.WriteRune(r)
			spaceFlag = false
		}
	}
	
	return result.String()
}

// NormalizeWhitespace converts all whitespace characters to single spaces
func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

// IsValidInput checks if input contains only printable characters
func IsValidInput(input string) bool {
	for _, r := range input {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
package main

import (
	"errors"
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
		return errors.New("record value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) DataRecord {
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp.UTC(),
		Tags:      append([]string{"processed"}, record.Tags...),
	}
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	
	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
}

func FilterByTag(records []DataRecord, tag string) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		for _, t := range record.Tags {
			if t == tag {
				filtered = append(filtered, record)
				break
			}
		}
	}
	return filtered
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func parseCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNum := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d", lineNum)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
		}

		record := DataRecord{
			ID:    id,
			Name:  row[1],
			Value: value,
		}

		if !validateRecord(record) {
			return nil, fmt.Errorf("validation failed for record at line %d", lineNum)
		}

		records = append(records, record)
		lineNum++
	}

	return records, nil
}

func validateRecord(record DataRecord) bool {
	if record.ID <= 0 {
		return false
	}
	if record.Name == "" {
		return false
	}
	if record.Value < 0 {
		return false
	}
	return true
}

func calculateStatistics(records []DataRecord) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}

func processDataFile(inputPath string) error {
	records, err := parseCSVFile(inputPath)
	if err != nil {
		return fmt.Errorf("data processing failed: %w", err)
	}

	avg, max := calculateStatistics(records)
	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)

	return nil
}