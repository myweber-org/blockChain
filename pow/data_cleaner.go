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
		if !seen[email] && len(email) > 0 {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	emailSet := make(map[string]bool)

	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		if ValidateEmail(cleanEmail) && !emailSet[cleanEmail] {
			emailSet[cleanEmail] = true
			record.Email = cleanEmail
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "another@domain.org", false},
		{5, "user@example.com", false},
	}

	cleaned := CleanRecords(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID   int
	Name string
	Age  int
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
	result := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := fmt.Sprintf("%d-%s-%d", record.ID, strings.ToLower(record.Name), record.Age)
		if !seen[key] {
			seen[key] = true
			result = append(result, record)
		}
	}

	dc.records = result
	return result
}

func (dc *DataCleaner) ValidateRecords() (valid []DataRecord, invalid []DataRecord) {
	valid = make([]DataRecord, 0)
	invalid = make([]DataRecord, 0)

	for _, record := range dc.records {
		if record.ID > 0 && record.Name != "" && record.Age >= 0 && record.Age <= 120 {
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

	cleaner.AddRecord(DataRecord{ID: 1, Name: "John", Age: 25})
	cleaner.AddRecord(DataRecord{ID: 2, Name: "Alice", Age: 30})
	cleaner.AddRecord(DataRecord{ID: 1, Name: "John", Age: 25})
	cleaner.AddRecord(DataRecord{ID: 3, Name: "", Age: 35})
	cleaner.AddRecord(DataRecord{ID: 4, Name: "Bob", Age: 150})

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", cleaner.GetRecordCount())

	valid, invalid := cleaner.ValidateRecords()
	fmt.Printf("Valid records: %d\n", len(valid))
	fmt.Printf("Invalid records: %d\n", len(invalid))

	for _, record := range valid {
		fmt.Printf("Valid: %+v\n", record)
	}
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct{}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (dc *DataCleaner) TrimWhitespace(items []string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func main() {
	cleaner := &DataCleaner{}
	data := []string{" apple ", "banana", " apple ", "  cherry  ", "banana"}

	fmt.Println("Original:", data)
	trimmed := cleaner.TrimWhitespace(data)
	fmt.Println("Trimmed:", trimmed)
	unique := cleaner.RemoveDuplicates(trimmed)
	fmt.Println("Cleaned:", unique)
}