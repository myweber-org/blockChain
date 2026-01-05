
package main

import (
	"strings"
)

func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func TrimAll(slice []string) []string {
	result := make([]string, len(slice))
	for i, item := range slice {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func CleanData(data []string) []string {
	trimmed := TrimAll(data)
	unique := RemoveDuplicates(trimmed)
	return unique
}
package main

import (
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record
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
		return errors.New("email must contain @ symbol")
	}
	if !strings.Contains(email, ".") {
		return errors.New("email must contain domain")
	}
	return nil
}

func CleanData(records []Record) ([]Record, error) {
	records = DeduplicateRecords(records)
	for i := range records {
		err := ValidateEmail(records[i].Email)
		records[i].Valid = err == nil
		if err != nil {
			fmt.Printf("Warning: Record ID %d has invalid email: %v\n", records[i].ID, err)
		}
	}
	return records, nil
}

func main() {
	sampleData := []Record{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "ANOTHER@TEST.ORG", false},
	}

	cleaned, err := CleanData(sampleData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Cleaned %d records\n", len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", r.ID, r.Email, r.Valid)
	}
}