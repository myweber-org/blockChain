package main

import (
	"encoding/csv"
	"io"
	"strings"
)

func FilterCSVRows(input io.Reader, output io.Writer, filterFunc func([]string) bool) error {
	reader := csv.NewReader(input)
	writer := csv.NewWriter(output)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if filterFunc(record) {
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}
	return nil
}

func DefaultFilter(record []string) bool {
	for _, field := range record {
		trimmed := strings.TrimSpace(field)
		if trimmed == "" || strings.Contains(trimmed, "NULL") {
			return false
		}
	}
	return true
}package main

import (
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

func ValidateEmails(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if strings.Contains(email, "@") && strings.Contains(email, ".") {
			record.Valid = true
			valid = append(valid, record)
		}
	}
	return valid
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
	}

	unique := RemoveDuplicates(records)
	validated := ValidateEmails(unique)

	fmt.Printf("Original: %d records\n", len(records))
	fmt.Printf("After deduplication: %d records\n", len(unique))
	fmt.Printf("Valid emails: %d records\n", len(validated))

	for _, r := range validated {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}