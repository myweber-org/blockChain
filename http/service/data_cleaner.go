
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

func DeduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] && strings.Contains(email, "@") {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func ValidateRecords(records []DataRecord) ([]DataRecord, []DataRecord) {
	var valid []DataRecord
	var invalid []DataRecord
	for _, record := range records {
		if record.ID > 0 && strings.Contains(record.Email, "@") {
			record.Valid = true
			valid = append(valid, record)
		} else {
			record.Valid = false
			invalid = append(invalid, record)
		}
	}
	return valid, invalid
}

func main() {
	emails := []string{
		"test@example.com",
		"TEST@example.com",
		"invalid-email",
		"another@test.org",
		"test@example.com",
	}
	uniqueEmails := DeduplicateEmails(emails)
	fmt.Println("Unique emails:", uniqueEmails)

	records := []DataRecord{
		{ID: 1, Email: "alice@example.com"},
		{ID: 0, Email: "bob@example.com"},
		{ID: 2, Email: "invalid"},
		{ID: 3, Email: "charlie@test.org"},
	}
	validRecs, invalidRecs := ValidateRecords(records)
	fmt.Printf("Valid records: %d\n", len(validRecs))
	fmt.Printf("Invalid records: %d\n", len(invalidRecs))
}