package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
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

	// Test duplicate removal
	records := []string{"user1@example.com", "User1@Example.com", "user2@test.org", "user1@example.com"}
	unique := cleaner.RemoveDuplicates(records)
	fmt.Printf("Unique records: %v\n", unique)

	// Test email validation
	emails := []string{"test@example.com", "invalid-email", "user@domain", "@domain.com"}
	for _, email := range emails {
		fmt.Printf("%s valid: %v\n", email, cleaner.ValidateEmail(email))
	}

	// Test input sanitization
	input := "  Hello\tWorld\nThis is a test\r\n"
	fmt.Printf("Sanitized: '%s'\n", cleaner.SanitizeInput(input))
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	normalizeCase bool
	trimSpaces    bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		normalizeCase: true,
		trimSpaces:    true,
	}
}

func (dc *DataCleaner) NormalizeString(input string) string {
	result := input

	if dc.trimSpaces {
		result = strings.TrimSpace(result)
	}

	if dc.normalizeCase {
		result = strings.ToLower(result)
	}

	return result
}

func (dc *DataCleaner) DeduplicateStrings(items []string) []string {
	seen := make(map[string]struct{})
	var unique []string

	for _, item := range items {
		normalized := dc.NormalizeString(item)
		if _, exists := seen[normalized]; !exists {
			seen[normalized] = struct{}{}
			unique = append(unique, normalized)
		}
	}

	return unique
}

func main() {
	cleaner := NewDataCleaner()

	data := []string{
		"  Apple  ",
		"apple",
		" BANANA ",
		"banana",
		"Cherry",
		"cherry ",
	}

	cleaned := cleaner.DeduplicateStrings(data)

	fmt.Println("Original data:", data)
	fmt.Println("Cleaned data:", cleaned)
}
package main

import (
    "fmt"
    "strings"
)

func DeduplicateStrings(slice []string) []string {
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

func ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    if len(parts[0]) == 0 || len(parts[1]) == 0 {
        return false
    }
    return true
}

func main() {
    emails := []string{
        "test@example.com",
        "user@domain.org",
        "test@example.com",
        "invalid-email",
        "another@test.net",
        "user@domain.org",
    }

    fmt.Println("Original list:", emails)
    
    uniqueEmails := DeduplicateStrings(emails)
    fmt.Println("Deduplicated list:", uniqueEmails)
    
    fmt.Println("\nEmail validation:")
    for _, email := range uniqueEmails {
        valid := ValidateEmail(email)
        fmt.Printf("%s: %v\n", email, valid)
    }
}