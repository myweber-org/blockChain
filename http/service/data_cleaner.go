package csvutils

import (
	"strings"
	"unicode"
)

// SanitizeString removes potentially problematic characters from CSV field values.
// It trims spaces, removes newlines and carriage returns, and replaces
// consecutive spaces with a single space.
func SanitizeString(input string) string {
	// Remove leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove newlines and carriage returns
	noNewlines := strings.ReplaceAll(trimmed, "\n", "")
	noNewlines = strings.ReplaceAll(noNewlines, "\r", "")
	
	// Replace multiple spaces with single space
	var result strings.Builder
	result.Grow(len(noNewlines))
	
	prevSpace := false
	for _, r := range noNewlines {
		if unicode.IsSpace(r) {
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}
	
	return result.String()
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
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmails(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		record.Valid = strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".")
		valid = append(valid, record)
	}
	return valid
}

func processData(records []DataRecord) []DataRecord {
	deduped := deduplicateRecords(records)
	validated := validateEmails(deduped)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
		{5, "TEST@DOMAIN.ORG", false},
	}

	processed := processData(sampleData)

	for _, record := range processed {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}