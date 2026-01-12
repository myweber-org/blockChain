
package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Phone string
}

type DataCleaner struct {
    records []DataRecord
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{
        records: make([]DataRecord, 0),
    }
}

func (dc *DataCleaner) AddRecord(record DataRecord) {
    dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
    seen := make(map[string]bool)
    unique := make([]DataRecord, 0)

    for _, record := range dc.records {
        key := fmt.Sprintf("%s|%s", record.Email, record.Phone)
        if !seen[key] {
            seen[key] = true
            unique = append(unique, record)
        }
    }

    dc.records = unique
    return unique
}

func (dc *DataCleaner) ValidateEmails() []DataRecord {
    valid := make([]DataRecord, 0)

    for _, record := range dc.records {
        if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
            valid = append(valid, record)
        }
    }

    return valid
}

func (dc *DataCleaner) GetRecordCount() int {
    return len(dc.records)
}

func main() {
    cleaner := NewDataCleaner()

    cleaner.AddRecord(DataRecord{ID: 1, Email: "test@example.com", Phone: "1234567890"})
    cleaner.AddRecord(DataRecord{ID: 2, Email: "test@example.com", Phone: "1234567890"})
    cleaner.AddRecord(DataRecord{ID: 3, Email: "invalid-email", Phone: "0987654321"})
    cleaner.AddRecord(DataRecord{ID: 4, Email: "another@test.org", Phone: "5551234567"})

    fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

    unique := cleaner.RemoveDuplicates()
    fmt.Printf("After deduplication: %d\n", len(unique))

    valid := cleaner.ValidateEmails()
    fmt.Printf("Valid email records: %d\n", len(valid))

    for _, record := range valid {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", record.ID, record.Email, record.Phone)
    }
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func cleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	for _, record := range records {
		if validateEmail(record.Email) {
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return deduplicateRecords(cleaned)
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "user@example.com", false},
		{4, "test@domain.org", false},
		{5, "another@test.co.uk", false},
	}

	cleaned := cleanData(sampleData)
	fmt.Printf("Original records: %d\n", len(sampleData))
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
	}
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}