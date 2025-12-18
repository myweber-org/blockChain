
package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Score int
}

type DataCleaner struct {
	records map[string]DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make(map[string]DataRecord),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if !strings.Contains(record.Email, "@") {
		return errors.New("invalid email format")
	}
	if record.Score < 0 || record.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}

	if _, exists := dc.records[record.ID]; exists {
		return fmt.Errorf("duplicate record ID: %s", record.ID)
	}

	dc.records[record.ID] = record
	return nil
}

func (dc *DataCleaner) RemoveRecord(id string) bool {
	if _, exists := dc.records[id]; exists {
		delete(dc.records, id)
		return true
	}
	return false
}

func (dc *DataCleaner) GetValidRecords() []DataRecord {
	var validRecords []DataRecord
	for _, record := range dc.records {
		validRecords = append(validRecords, record)
	}
	return validRecords
}

func (dc *DataCleaner) CalculateAverageScore() float64 {
	if len(dc.records) == 0 {
		return 0.0
	}

	total := 0
	for _, record := range dc.records {
		total += record.Score
	}
	return float64(total) / float64(len(dc.records))
}

func main() {
	cleaner := NewDataCleaner()

	sampleRecords := []DataRecord{
		{ID: "001", Email: "user1@example.com", Score: 85},
		{ID: "002", Email: "user2@example.com", Score: 92},
		{ID: "003", Email: "user3@example.com", Score: 78},
	}

	for _, record := range sampleRecords {
		if err := cleaner.AddRecord(record); err != nil {
			fmt.Printf("Failed to add record %s: %v\n", record.ID, err)
		}
	}

	records := cleaner.GetValidRecords()
	fmt.Printf("Total valid records: %d\n", len(records))
	fmt.Printf("Average score: %.2f\n", cleaner.CalculateAverageScore())
}