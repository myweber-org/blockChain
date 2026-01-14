package main

import (
	"fmt"
	"sort"
)

type Record struct {
	ID   int
	Name string
}

func cleanData(records []Record) []Record {
	seen := make(map[int]bool)
	var unique []Record

	for _, r := range records {
		if !seen[r.ID] {
			seen[r.ID] = true
			unique = append(unique, r)
		}
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].ID < unique[j].ID
	})

	return unique
}

func main() {
	data := []Record{
		{3, "Charlie"},
		{1, "Alice"},
		{2, "Bob"},
		{1, "Alice"},
		{4, "David"},
	}

	cleaned := cleanData(data)
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
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
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, DataRecord{
				ID:    record.ID,
				Email: email,
				Valid: record.Valid,
			})
		}
	}
	return unique
}

func validateEmailFormat(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func filterValidRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if record.Valid && validateEmailFormat(record.Email) {
			valid = append(valid, record)
		}
	}
	return valid
}

func processDataset(records []DataRecord) []DataRecord {
	deduped := deduplicateRecords(records)
	validated := filterValidRecords(deduped)
	return validated
}

func main() {
	dataset := []DataRecord{
		{1, "user@example.com", true},
		{2, "user@example.com", true},
		{3, "invalid-email", true},
		{4, "test@domain.org", false},
		{5, "another@test.net", true},
	}

	cleaned := processDataset(dataset)
	fmt.Printf("Original: %d records\n", len(dataset))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))

	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s\n", record.ID, record.Email)
	}
}