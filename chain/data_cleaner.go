
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) ProcessBatch(items []string) []string {
	var uniqueItems []string
	for _, item := range items {
		if !dc.IsDuplicate(item) {
			uniqueItems = append(uniqueItems, item)
		}
	}
	return uniqueItems
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"  Apple  ",
		"apple",
		"BANANA",
		"  banana  ",
		"Cherry",
		"Date",
		"date",
	}
	
	result := cleaner.ProcessBatch(data)
	fmt.Println("Unique items:", result)
	fmt.Println("Total unique count:", len(result))
}package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func CleanData(records []DataRecord) ([]DataRecord, error) {
	cleaned := RemoveDuplicates(records)
	for i := range cleaned {
		if err := ValidateEmail(cleaned[i].Email); err != nil {
			cleaned[i].Valid = false
			fmt.Printf("Warning: Record ID %d has invalid email: %v\n", cleaned[i].ID, err)
		} else {
			cleaned[i].Valid = true
		}
	}
	return cleaned, nil
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "", false},
	}

	cleaned, err := CleanData(sampleData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Cleaned %d records\n", len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}