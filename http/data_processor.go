
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
		return errors.New("record value must be non-negative")
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

	transformed := DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp.UTC(),
		Tags:      append([]string{}, record.Tags...),
	}

	if len(transformed.Tags) == 0 {
		transformed.Tags = append(transformed.Tags, "default")
	}

	return transformed, nil
}

func ProcessBatch(records []DataRecord, multiplier float64) ([]DataRecord, []error) {
	var processed []DataRecord
	var errs []error

	for i, record := range records {
		transformed, err := TransformRecord(record, multiplier)
		if err != nil {
			errs = append(errs, fmt.Errorf("record %d: %w", i, err))
			continue
		}
		processed = append(processed, transformed)
	}

	return processed, errs
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records provided")
	}

	var sum float64
	var count int

	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			continue
		}
		sum += record.Value
		count++
	}

	if count == 0 {
		return 0, 0, errors.New("no valid records found")
	}

	average := sum / float64(count)

	var varianceSum float64
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			continue
		}
		diff := record.Value - average
		varianceSum += diff * diff
	}

	variance := varianceSum / float64(count)

	return average, variance, nil
}