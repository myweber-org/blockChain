
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
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) DataRecord {
	if multiplier <= 0 {
		multiplier = 1.0
	}
	
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp.UTC(),
		Tags:      append([]string{"processed"}, record.Tags...),
	}
}

func FilterRecords(records []DataRecord, minValue float64) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Value >= minValue && ValidateRecord(record) == nil {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func CalculateAverage(records []DataRecord) (float64, error) {
	if len(records) == 0 {
		return 0, errors.New("no records to process")
	}
	
	var sum float64
	validCount := 0
	
	for _, record := range records {
		if ValidateRecord(record) == nil {
			sum += record.Value
			validCount++
		}
	}
	
	if validCount == 0 {
		return 0, errors.New("no valid records found")
	}
	
	return sum / float64(validCount), nil
}