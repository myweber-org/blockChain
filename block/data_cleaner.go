
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID   int
	Name string
	Age  int
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[int]bool)
	var result []DataRecord
	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			result = append(result, record)
		}
	}
	return result
}

func ValidateRecord(record DataRecord) error {
	if record.ID <= 0 {
		return fmt.Errorf("invalid ID: %d", record.ID)
	}
	if strings.TrimSpace(record.Name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if record.Age < 0 || record.Age > 150 {
		return fmt.Errorf("age out of range: %d", record.Age)
	}
	return nil
}

func CleanData(records []DataRecord) ([]DataRecord, []string) {
	var cleaned []DataRecord
	var errors []string

	uniqueRecords := RemoveDuplicates(records)

	for _, record := range uniqueRecords {
		if err := ValidateRecord(record); err != nil {
			errors = append(errors, fmt.Sprintf("Record ID %d: %v", record.ID, err))
		} else {
			cleaned = append(cleaned, record)
		}
	}

	return cleaned, errors
}

func main() {
	sampleData := []DataRecord{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 3, Name: "", Age: 40},
		{ID: 4, Name: "Charlie", Age: 200},
	}

	cleaned, errors := CleanData(sampleData)

	fmt.Println("Cleaned Records:")
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", record.ID, record.Name, record.Age)
	}

	if len(errors) > 0 {
		fmt.Println("\nValidation Errors:")
		for _, err := range errors {
			fmt.Println(err)
		}
	}
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
	numbers := []int{1, 2, 2, 3, 4, 4, 5, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange", "banana"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}