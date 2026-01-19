
package main

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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
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

    strings := []string{"apple", "banana", "apple", "orange", "banana"}
    uniqueStrings := RemoveDuplicates(strings)
    fmt.Println("Original:", strings)
    fmt.Println("Unique:", uniqueStrings)
}package main

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

func DeduplicateEmails(emails []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, email := range emails {
		if _, exists := seen[email]; !exists {
			seen[email] = struct{}{}
			result = append(result, email)
		}
	}
	return result
}

func ValidateEmail(email string) error {
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

func CleanRecords(records []Record) ([]Record, error) {
	cleaned := []Record{}
	emailSet := make(map[string]struct{})

	for _, rec := range records {
		if err := ValidateEmail(rec.Email); err != nil {
			fmt.Printf("Skipping record %d: %v\n", rec.ID, err)
			continue
		}

		if _, exists := emailSet[rec.Email]; exists {
			fmt.Printf("Duplicate email found for record %d: %s\n", rec.ID, rec.Email)
			continue
		}

		emailSet[rec.Email] = struct{}{}
		rec.Valid = true
		cleaned = append(cleaned, rec)
	}

	if len(cleaned) == 0 {
		return cleaned, errors.New("no valid records after cleaning")
	}

	return cleaned, nil
}

func main() {
	records := []Record{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "user@example.com", false},
		{4, "another@domain.org", false},
	}

	cleaned, err := CleanRecords(records)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Cleaned %d records:\n", len(cleaned))
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", rec.ID, rec.Email, rec.Valid)
	}
}