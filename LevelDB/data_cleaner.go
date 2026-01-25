
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Name  string
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(id int, email, name string) {
	record := DataRecord{
		ID:    id,
		Email: strings.ToLower(strings.TrimSpace(email)),
		Name:  strings.TrimSpace(name),
	}
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]bool)
	unique := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := fmt.Sprintf("%s|%s", record.Email, record.Name)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}

	dc.records = unique
	return unique
}

func (dc *DataCleaner) ValidateEmails() (valid, invalid []DataRecord) {
	for _, record := range dc.records {
		if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}
	return valid, invalid
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(1, "user@example.com", "John Doe")
	cleaner.AddRecord(2, "user@example.com", "John Doe")
	cleaner.AddRecord(3, "invalid-email", "Jane Smith")
	cleaner.AddRecord(4, "another@test.org", "Alice Johnson")

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	unique := cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", len(unique))

	valid, invalid := cleaner.ValidateEmails()
	fmt.Printf("Valid emails: %d, Invalid emails: %d\n", len(valid), len(invalid))

	for _, record := range valid {
		fmt.Printf("Valid: ID=%d, Email=%s, Name=%s\n", record.ID, record.Email, record.Name)
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package utils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove any HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(trimmed, "")
	
	// Escape special characters for SQL (basic prevention)
	cleaned = strings.ReplaceAll(cleaned, "'", "''")
	cleaned = strings.ReplaceAll(cleaned, "\"", "\"\"")
	
	// Limit length to prevent buffer overflow attacks
	maxLength := 255
	if len(cleaned) > maxLength {
		cleaned = cleaned[:maxLength]
	}
	
	return cleaned
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
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

func deduplicateEmails(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && isValidEmail(email) {
			seen[email] = true
			record.Email = email
			record.Valid = true
			unique = append(unique, record)
		}
	}
	return unique
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanData(records []DataRecord) []DataRecord {
	cleaned := deduplicateEmails(records)
	fmt.Printf("Cleaned %d records to %d unique entries\n", len(records), len(cleaned))
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
		{5, "test@domain.org", false},
	}

	result := cleanData(sampleData)
	for _, r := range result {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func validateRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		record.Valid = validateEmail(record.Email)
		valid = append(valid, record)
	}
	return valid
}

func cleanData(records []DataRecord) []DataRecord {
	unique := deduplicateRecords(records)
	validated := validateRecords(unique)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob Wilson", "invalid-email", false},
		{5, "Alice Brown", "alice@example", false},
	}

	cleaned := cleanData(sampleData)

	fmt.Println("Cleaned Records:")
	for _, record := range cleaned {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, Status: %s\n",
			record.ID, record.Name, record.Email, status)
	}
}