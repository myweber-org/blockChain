package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
	records = DeduplicateRecords(records)
	for i := range records {
		records[i].Valid = ValidateEmail(records[i].Email)
	}
	return records
}

func main() {
	records := []DataRecord{
		{ID: "1", Email: "test@example.com"},
		{ID: "2", Email: "invalid-email"},
		{ID: "1", Email: "duplicate@example.com"},
		{ID: "3", Email: "another@test.org"},
	}

	cleaned := CleanData(records)
	for _, r := range cleaned {
		fmt.Printf("ID: %s, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}