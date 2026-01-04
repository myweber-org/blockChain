
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
}package main

import (
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	ID   int
	Name string
	Tags []string
}

func RemoveDuplicates(records []Record) []Record {
	seen := make(map[int]bool)
	result := []Record{}
	for _, rec := range records {
		if !seen[rec.ID] {
			seen[rec.ID] = true
			result = append(result, rec)
		}
	}
	return result
}

func ValidateRecord(rec Record) error {
	if rec.ID <= 0 {
		return errors.New("invalid record ID")
	}
	if strings.TrimSpace(rec.Name) == "" {
		return errors.New("record name cannot be empty")
	}
	return nil
}

func CleanData(records []Record) ([]Record, error) {
	cleaned := RemoveDuplicates(records)
	for _, rec := range cleaned {
		if err := ValidateRecord(rec); err != nil {
			return nil, fmt.Errorf("validation failed for record %d: %w", rec.ID, err)
		}
	}
	return cleaned, nil
}

func main() {
	sampleData := []Record{
		{ID: 1, Name: "Alpha", Tags: []string{"tag1"}},
		{ID: 2, Name: "Beta", Tags: []string{"tag2"}},
		{ID: 1, Name: "Alpha", Tags: []string{"tag1"}},
		{ID: 3, Name: "", Tags: []string{"tag3"}},
	}

	cleaned, err := CleanData(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Name: %s\n", rec.ID, rec.Name)
	}
}