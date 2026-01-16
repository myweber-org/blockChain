
package main

import "fmt"

func removeDuplicates(nums []int) []int {
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
	cleaned := removeDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
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

func DeduplicateEmails(emails []string) []string {
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

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        record.Email = strings.ToLower(strings.TrimSpace(record.Email))
        
        if ValidateEmail(record.Email) && !emailSet[record.Email] {
            emailSet[record.Email] = true
            record.Valid = true
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", false},
        {2, "USER@example.com", false},
        {3, "invalid-email", false},
        {4, "user@example.com", false},
        {5, "new@domain.org", false},
    }
    
    cleaned := CleanData(sampleData)
    fmt.Printf("Cleaned %d records\n", len(cleaned))
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
    }
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

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		key := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmails(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if strings.Contains(record.Email, "@") && len(record.Email) > 3 {
			record.Valid = true
			valid = append(valid, record)
		}
	}
	return valid
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
	}

	unique := deduplicateRecords(records)
	validated := validateEmails(unique)

	fmt.Printf("Original: %d records\n", len(records))
	fmt.Printf("Deduplicated: %d records\n", len(unique))
	fmt.Printf("Validated: %d records\n", len(validated))

	for _, r := range validated {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}