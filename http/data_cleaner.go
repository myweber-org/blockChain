package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import "fmt"

func removeDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Value float64
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

	if !isValidEmail(record.Email) {
		return errors.New("invalid email format")
	}

	if record.Value < 0 {
		return errors.New("value cannot be negative")
	}

	if _, exists := dc.records[record.ID]; exists {
		return fmt.Errorf("duplicate record ID: %s", record.ID)
	}

	dc.records[record.ID] = record
	return nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (dc *DataCleaner) GetRecord(id string) (DataRecord, error) {
	record, exists := dc.records[id]
	if !exists {
		return DataRecord{}, fmt.Errorf("record not found: %s", id)
	}
	return record, nil
}

func (dc *DataCleaner) RemoveRecord(id string) error {
	if _, exists := dc.records[id]; !exists {
		return fmt.Errorf("record not found: %s", id)
	}
	delete(dc.records, id)
	return nil
}

func (dc *DataCleaner) TotalRecords() int {
	return len(dc.records)
}

func (dc *DataCleaner) AverageValue() float64 {
	if len(dc.records) == 0 {
		return 0
	}

	var sum float64
	for _, record := range dc.records {
		sum += record.Value
	}
	return sum / float64(len(dc.records))
}

func main() {
	cleaner := NewDataCleaner()

	records := []DataRecord{
		{ID: "001", Email: "user1@example.com", Value: 42.5},
		{ID: "002", Email: "user2@domain.org", Value: 78.9},
		{ID: "003", Email: "user3@test.net", Value: 15.2},
	}

	for _, record := range records {
		if err := cleaner.AddRecord(record); err != nil {
			fmt.Printf("Failed to add record %s: %v\n", record.ID, err)
		}
	}

	fmt.Printf("Total records: %d\n", cleaner.TotalRecords())
	fmt.Printf("Average value: %.2f\n", cleaner.AverageValue())

	if record, err := cleaner.GetRecord("002"); err == nil {
		fmt.Printf("Record 002: %+v\n", record)
	}
}