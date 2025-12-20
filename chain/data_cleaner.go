
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
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
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
    numbers := []int{1, 2, 2, 3, 4, 4, 5}
    uniqueNumbers := RemoveDuplicates(numbers)
    fmt.Println("Original:", numbers)
    fmt.Println("Unique:", uniqueNumbers)
    
    strings := []string{"apple", "banana", "apple", "orange", "banana"}
    uniqueStrings := RemoveDuplicates(strings)
    fmt.Println("Original:", strings)
    fmt.Println("Unique:", uniqueStrings)
}
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
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type Cleaner struct {
	skipEmptyRows bool
	trimSpaces    bool
}

func NewCleaner(options ...func(*Cleaner)) *Cleaner {
	c := &Cleaner{
		skipEmptyRows: true,
		trimSpaces:    true,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func WithSkipEmptyRows(skip bool) func(*Cleaner) {
	return func(c *Cleaner) {
		c.skipEmptyRows = skip
	}
}

func WithTrimSpaces(trim bool) func(*Cleaner) {
	return func(c *Cleaner) {
		c.trimSpaces = trim
	}
}

func (c *Cleaner) ProcessCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if c.trimSpaces {
		for i := range headers {
			headers[i] = strings.TrimSpace(headers[i])
		}
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		if c.skipEmptyRows && isEmptyRecord(record) {
			continue
		}

		if c.trimSpaces {
			for i := range record {
				record[i] = strings.TrimSpace(record[i])
			}
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func isEmptyRecord(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}

func main() {
	cleaner := NewCleaner(
		WithSkipEmptyRows(true),
		WithTrimSpaces(true),
	)

	err := cleaner.ProcessCSV("input.csv", "output.csv")
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("CSV processing completed successfully")
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
	data := []string{"  apple ", "banana", "  apple ", " cherry", "banana "}

	fmt.Println("Original:", data)
	trimmed := cleaner.TrimWhitespace(data)
	fmt.Println("Trimmed:", trimmed)
	unique := cleaner.RemoveDuplicates(trimmed)
	fmt.Println("Unique:", unique)
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

func DeduplicateEmails(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmailFormat(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func CleanDataset(records []DataRecord) []DataRecord {
	records = DeduplicateEmails(records)
	for i := range records {
		records[i].Valid = ValidateEmailFormat(records[i].Email)
	}
	return records
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "another@domain.org", false},
		{5, "  user@example.com  ", false},
	}

	cleaned := CleanDataset(sampleData)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(sampleData), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}
package main

import (
	"fmt"
	"strings"
)

func CleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	data := []string{" apple", "banana ", "  apple  ", "banana", "", "cherry"}
	cleaned := CleanData(data)
	fmt.Println("Cleaned data:", cleaned)
}