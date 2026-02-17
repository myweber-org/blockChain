
package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
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
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	replacer := strings.NewReplacer("\n", " ", "\t", " ", "\r", " ")
	return replacer.Replace(trimmed)
}

func main() {
	cleaner := NewDataCleaner()
	
	duplicateData := []string{"john@example.com", "JANE@test.org", "john@example.com", "  alice@domain.com  "}
	uniqueEmails := cleaner.RemoveDuplicates(duplicateData)
	fmt.Println("Unique emails:", uniqueEmails)
	
	testEmails := []string{"valid@example.com", "invalid", "user@domain", "@domain.com"}
	for _, email := range testEmails {
		fmt.Printf("Email '%s' valid: %v\n", email, cleaner.ValidateEmail(email))
	}
	
	dirtyInput := "  Hello\tWorld\nThis is a test\r\n"
	fmt.Println("Sanitized input:", cleaner.SanitizeInput(dirtyInput))
}