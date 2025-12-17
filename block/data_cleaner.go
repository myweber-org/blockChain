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
    numbers := []int{1, 2, 2, 3, 4, 4, 5, 5, 5}
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
	records []string
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]string, 0),
	}
}

func (dc *DataCleaner) AddRecord(record string) {
	dc.records = append(dc.records, strings.TrimSpace(record))
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, record := range dc.records {
		if !seen[record] {
			seen[record] = true
			result = append(result, record)
		}
	}

	dc.records = result
	return result
}

func (dc *DataCleaner) ValidateRecords() (valid []string, invalid []string) {
	valid = make([]string, 0)
	invalid = make([]string, 0)

	for _, record := range dc.records {
		if len(record) > 0 && !strings.ContainsAny(record, "!@#$%") {
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

	sampleData := []string{
		"user1@example.com",
		"user2@test.org",
		"user1@example.com",
		"",
		"bad!user@domain.com",
		"  user3@sample.net  ",
	}

	for _, data := range sampleData {
		cleaner.AddRecord(data)
	}

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	unique := cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", len(unique))

	valid, invalid := cleaner.ValidateRecords()
	fmt.Printf("Valid records: %d\n", len(valid))
	fmt.Printf("Invalid records: %d\n", len(invalid))

	fmt.Println("\nValid records:")
	for _, v := range valid {
		fmt.Printf("  - %s\n", v)
	}
}