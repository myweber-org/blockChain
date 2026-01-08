
package main

import (
	"fmt"
	"strings"
)

func CleanStringSlice(input []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; !exists {
			seen[trimmed] = struct{}{}
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	data := []string{" apple ", "banana", " apple", "banana ", " ", "cherry"}
	cleaned := CleanStringSlice(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processed map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processed: make(map[string]bool),
	}
}

func (dc *DataCleaner) Clean(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("empty input")
	}

	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)

	if dc.processed[lower] {
		return "", fmt.Errorf("duplicate entry: %s", trimmed)
	}

	dc.processed[lower] = true
	return trimmed, nil
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func main() {
	cleaner := NewDataCleaner()

	samples := []string{
		"  User@Example.com  ",
		"user@example.com",
		"invalid-email",
		"another@test.org",
	}

	for _, sample := range samples {
		cleaned, err := cleaner.Clean(sample)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		isValid := cleaner.ValidateEmail(cleaned)
		fmt.Printf("Original: '%s' -> Cleaned: '%s' (Valid: %v)\n", sample, cleaned, isValid)
	}
}package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
    if len(nums) == 0 {
        return nums
    }

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
    fmt.Println("Original:", data)
    fmt.Println("Cleaned:", cleaned)
}