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
	Category  string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Category == "" {
		return errors.New("category cannot be empty")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("timestamp cannot be in the future")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	record.Category = strings.ToUpper(strings.TrimSpace(record.Category))
	record.ID = strings.ReplaceAll(record.ID, " ", "_")
	if record.Value > 1000 {
		record.Value = record.Value / 1000
	}
	return record
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

func main() {
	records := []DataRecord{
		{ID: "rec 001", Value: 1500.5, Timestamp: time.Now().Add(-time.Hour), Category: "category A"},
		{ID: "rec002", Value: 500.0, Timestamp: time.Now().Add(-2 * time.Hour), Category: "category B"},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, rec := range processed {
		fmt.Printf("Processed: %+v\n", rec)
	}
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func ValidateJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("JSON object is empty")
	}
	return result, nil
}

func main() {
	jsonData := `{"name": "test", "value": 42}`
	parsed, err := ValidateJSON([]byte(jsonData))
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}
	fmt.Printf("Parsed data: %v\n", parsed)
}