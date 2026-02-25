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
}package utils

import (
	"regexp"
	"strings"
)

func SanitizeString(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(trimmed, " ")
	
	return cleaned
}

func SanitizeAndLower(input string) string {
	cleaned := SanitizeString(input)
	return strings.ToLower(cleaned)
}package main

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
	numbers := []int{1, 2, 2, 3, 4, 4, 5, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange", "banana"}
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
	processedRecords int
	duplicatesFound  int
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: 0,
		duplicatesFound:  0,
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	seen := make(map[string]bool)
	var unique []string

	for _, record := range records {
		trimmed := strings.TrimSpace(record)
		if trimmed == "" {
			continue
		}

		if !seen[trimmed] {
			seen[trimmed] = true
			unique = append(unique, trimmed)
		} else {
			dc.duplicatesFound++
		}
		dc.processedRecords++
	}

	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	if parts[0] == "" || parts[1] == "" {
		return false
	}

	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) Stats() string {
	return fmt.Sprintf("Processed: %d, Duplicates: %d", 
		dc.processedRecords, dc.duplicatesFound)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"user@example.com",
		"user@example.com",
		"invalid-email",
		"another@test.org",
		"",
		"  user@example.com  ",
	}

	unique := cleaner.RemoveDuplicates(data)
	fmt.Println("Unique records:", unique)
	fmt.Println(cleaner.Stats())

	for _, email := range unique {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("Valid email: %s\n", email)
		} else {
			fmt.Printf("Invalid email: %s\n", email)
		}
	}
}