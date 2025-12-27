package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ProcessRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("invalid record ID")
	}

	if strings.TrimSpace(record.Name) == "" {
		return errors.New("record name cannot be empty")
	}

	if record.Value < 0 {
		return errors.New("record value cannot be negative")
	}

	fmt.Printf("Processing record %d: %s (%.2f)\n", record.ID, record.Name, record.Value)
	return nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, error) {
	var validRecords []DataRecord
	var validationErrors []string

	for _, record := range records {
		err := ProcessRecord(record)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("Record %d: %v", record.ID, err))
			continue
		}
		validRecords = append(validRecords, record)
	}

	if len(validationErrors) > 0 {
		return validRecords, fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
	}

	return validRecords, nil
}

func CalculateTotal(records []DataRecord) float64 {
	var total float64
	for _, record := range records {
		total += record.Value
	}
	return total
}