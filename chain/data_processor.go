package main

import (
	"fmt"
)

func MovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
		return nil
	}

	result := make([]float64, len(data)-windowSize+1)
	for i := 0; i < len(result); i++ {
		sum := 0.0
		for j := 0; j < windowSize; j++ {
			sum += data[i+j]
		}
		result[i] = sum / float64(windowSize)
	}
	return result
}

func main() {
	sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}
	averages := MovingAverage(sampleData, 3)
	fmt.Println("Moving averages:", averages)
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
	Category  string
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
	if strings.TrimSpace(record.Category) == "" {
		return errors.New("category cannot be empty")
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
		Category:  strings.ToLower(strings.TrimSpace(record.Category)),
	}
}

func ProcessRecords(records []DataRecord, multiplier float64) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, err
		}
		processed = append(processed, TransformRecord(record, multiplier))
	}
	return processed, nil
}