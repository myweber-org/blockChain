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
	data := []string{"apple", "banana", "apple", "cherry", "banana", "date"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
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
	fmt.Printf("Cleaned: %v\n", cleaned)
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
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
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
	Valid bool
}

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func processRecords(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord
	
	for _, rec := range records {
		email := strings.ToLower(strings.TrimSpace(rec.Email))
		
		if validateEmail(email) && !emailSet[email] {
			emailSet[email] = true
			rec.Email = email
			rec.Valid = true
			cleaned = append(cleaned, rec)
		}
	}
	return cleaned
}

func main() {
	emails := []string{
		"test@example.com",
		"TEST@example.com",
		"invalid-email",
		"another@test.org",
		"test@example.com",
		"  spaced@domain.com  ",
	}
	
	fmt.Println("Original emails:", emails)
	deduped := deduplicateEmails(emails)
	fmt.Println("Deduplicated:", deduped)
	
	records := []DataRecord{
		{1, "user@domain.com", false},
		{2, "DUPLICATE@domain.com", false},
		{3, "user@domain.com", false},
		{4, "bad-format", false},
		{5, "new@test.io", false},
	}
	
	cleaned := processRecords(records)
	fmt.Printf("\nCleaned records: %d out of %d\n", len(cleaned), len(records))
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", rec.ID, rec.Email, rec.Valid)
	}
}