package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(input))

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
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

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	var unique []string
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	replacer := strings.NewReplacer("\n", " ", "\r", " ", "\t", " ")
	return replacer.Replace(trimmed)
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{
		"user@example.com",
		"  USER@EXAMPLE.COM  ",
		"invalid-email",
		"another@test.org",
		"user@example.com",
	}
	
	fmt.Println("Original records:", records)
	deduped := cleaner.RemoveDuplicates(records)
	fmt.Println("Deduplicated:", deduped)
	
	for _, email := range deduped {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("Valid email: %s\n", cleaner.SanitizeInput(email))
		} else {
			fmt.Printf("Invalid email: %s\n", email)
		}
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

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var result []DataRecord
	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			result = append(result, record)
		}
	}
	return result
}

func ValidateEmails(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if strings.Contains(record.Email, "@") && len(record.Email) > 3 {
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func PrintRecords(records []DataRecord) {
	for _, record := range records {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", record.ID, record.Email, status)
	}
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "admin@test.org", false},
		{3, "user@example.com", false},
		{4, "invalid-email", false},
		{5, "test@domain.com", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	uniqueRecords := RemoveDuplicates(records)
	fmt.Println("\nAfter deduplication:")
	PrintRecords(uniqueRecords)

	validRecords := ValidateEmails(uniqueRecords)
	fmt.Println("\nValid records only:")
	PrintRecords(validRecords)
}