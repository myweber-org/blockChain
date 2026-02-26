package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput removes potentially harmful characters and trims whitespace
func SanitizeInput(input string) string {
	// Trim leading and trailing whitespace
	cleaned := strings.TrimSpace(input)
	
	// Remove any HTML/XML tags
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned = re.ReplaceAllString(cleaned, "")
	
	// Remove control characters except newline and tab
	re = regexp.MustCompile(`[\x00-\x08\x0B-\x0C\x0E-\x1F\x7F]`)
	cleaned = re.ReplaceAllString(cleaned, "")
	
	// Limit length to prevent excessive input
	maxLength := 1000
	if len(cleaned) > maxLength {
		cleaned = cleaned[:maxLength]
	}
	
	return cleaned
}

// NormalizeWhitespace converts multiple whitespace characters to single spaces
func NormalizeWhitespace(input string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(input, " ")
}
package main

import "fmt"

func removeDuplicates(input []int) []int {
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
	data := []int{5, 2, 8, 2, 5, 9, 1, 8}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
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

func (dc *DataCleaner) AddRecord(record DataRecord) {
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]bool)
	unique := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := fmt.Sprintf("%d|%s", record.ID, strings.ToLower(record.Email))
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
		if strings.Contains(record.Email, "@") && len(record.Name) > 0 {
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

	cleaner.AddRecord(DataRecord{ID: 1, Email: "user@example.com", Name: "John"})
	cleaner.AddRecord(DataRecord{ID: 2, Email: "user@example.com", Name: "Jane"})
	cleaner.AddRecord(DataRecord{ID: 3, Email: "test@domain.com", Name: "Bob"})
	cleaner.AddRecord(DataRecord{ID: 4, Email: "invalid-email", Name: "Alice"})

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	unique := cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", len(unique))

	valid := cleaner.ValidateEmails()
	fmt.Printf("Valid records: %d\n", len(valid))
}package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned:  %v\n", cleaned)
}