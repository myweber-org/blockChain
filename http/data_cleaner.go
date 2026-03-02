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

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func ProcessRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if ValidateEmail(record.Email) {
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return DeduplicateRecords(validRecords)
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "admin@test.org", false},
		{3, "USER@example.com", false},
		{4, "invalid-email", false},
		{5, "test@domain", false},
	}

	cleaned := ProcessRecords(sampleData)
	fmt.Printf("Processed %d records\n", len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
	}
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
    numbers := []int{1, 2, 2, 3, 4, 4, 5, 5, 5}
    uniqueNumbers := RemoveDuplicates(numbers)
    fmt.Println("Original:", numbers)
    fmt.Println("Unique:", uniqueNumbers)
    
    strings := []string{"apple", "banana", "apple", "orange", "banana"}
    uniqueStrings := RemoveDuplicates(strings)
    fmt.Println("Original:", strings)
    fmt.Println("Unique:", uniqueStrings)
}