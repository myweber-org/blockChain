
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

func TransformRecord(record DataRecord, multiplier float64) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, err
	}

	transformed := record
	transformed.Value = record.Value * multiplier

	for i, tag := range transformed.Tags {
		transformed.Tags[i] = strings.ToUpper(strings.TrimSpace(tag))
	}

	return transformed, nil
}

func ProcessRecords(records []DataRecord, multiplier float64) ([]DataRecord, []error) {
	var processed []DataRecord
	var errs []error

	for _, record := range records {
		transformed, err := TransformRecord(record, multiplier)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to process record %s: %w", record.ID, err))
			continue
		}
		processed = append(processed, transformed)
	}

	return processed, errs
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