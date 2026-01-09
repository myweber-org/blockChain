
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor(pattern string) (*DataProcessor, error) {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &DataProcessor{allowedPattern: compiled}, nil
}

func (dp *DataProcessor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.allowedPattern.FindString(trimmed)
}

func (dp *DataProcessor) Validate(input string) bool {
	return dp.allowedPattern.MatchString(input)
}

func (dp *DataProcessor) ProcessBatch(inputs []string) []string {
	var results []string
	for _, item := range inputs {
		cleaned := dp.CleanInput(item)
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}
package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

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

		if len(row) < 4 {
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

		valid := false
		if row[3] == "true" {
			valid = true
		}

		record := DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
			Valid: valid,
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, error) {
	if len(records) == 0 {
		return nil, errors.New("no records to validate")
	}

	validRecords := make([]DataRecord, 0)
	for _, record := range records {
		if record.ID > 0 && record.Name != "" && record.Value >= 0 {
			validRecords = append(validRecords, record)
		}
	}

	return validRecords, nil
}

func CalculateTotalValue(records []DataRecord) float64 {
	total := 0.0
	for _, record := range records {
		if record.Valid {
			total += record.Value
		}
	}
	return total
}