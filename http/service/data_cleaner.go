
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

	strings := []string{"apple", "banana", "apple", "orange"}
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
	processedRecords map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	var unique []string
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	replacer := strings.NewReplacer("\n", " ", "\t", " ", "\r", " ")
	return replacer.Replace(trimmed)
}

func main() {
	cleaner := NewDataCleaner()
	
	duplicateData := []string{"john@example.com", "JANE@test.org", "john@example.com", "  alice@domain.com  "}
	uniqueEmails := cleaner.RemoveDuplicates(duplicateData)
	fmt.Println("Unique emails:", uniqueEmails)
	
	testEmails := []string{"valid@example.com", "invalid", "user@domain", "@domain.com"}
	for _, email := range testEmails {
		fmt.Printf("Email '%s' valid: %v\n", email, cleaner.ValidateEmail(email))
	}
	
	dirtyInput := "  Hello\tWorld\nThis is a test\r\n"
	fmt.Println("Sanitized input:", cleaner.SanitizeInput(dirtyInput))
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	var unique []string
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) CleanPhoneNumber(phone string) string {
	var builder strings.Builder
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{
		"user@example.com",
		"USER@EXAMPLE.COM",
		"invalid-email",
		"  user@example.com  ",
		"test@domain.org",
	}
	
	fmt.Println("Original records:", records)
	unique := cleaner.RemoveDuplicates(records)
	fmt.Println("After deduplication:", unique)
	
	for _, email := range unique {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("Valid email: %s\n", email)
		} else {
			fmt.Printf("Invalid email: %s\n", email)
		}
	}
	
	phone := "+1 (555) 123-4567"
	cleaned := cleaner.CleanPhoneNumber(phone)
	fmt.Printf("Original: %s -> Cleaned: %s\n", phone, cleaned)
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSV(inputPath, outputPath string) error {
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

	trimmedHeaders := make([]string, len(headers))
	for i, h := range headers {
		trimmedHeaders[i] = strings.TrimSpace(h)
	}
	if err := writer.Write(trimmedHeaders); err != nil {
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

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedField := strings.TrimSpace(field)
			cleanedField = strings.ToLower(cleanedField)
			cleanedRecord[i] = cleanedField
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully cleaned data. Output saved to %s\n", outputFile)
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSV(inputPath, outputPath string) error {
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

	trimmedHeaders := make([]string, len(headers))
	for i, h := range headers {
		trimmedHeaders[i] = strings.TrimSpace(h)
	}
	if err := writer.Write(trimmedHeaders); err != nil {
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

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedField := strings.TrimSpace(field)
			cleanedField = strings.ToLower(cleanedField)
			cleanedRecord[i] = cleanedField
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write cleaned record: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully cleaned data. Output saved to %s\n", outputFile)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	Data []string
}

func NewDataCleaner(data []string) *DataCleaner {
	return &DataCleaner{Data: data}
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range dc.Data {
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

func (dc *DataCleaner) Clean() []string {
	return dc.RemoveDuplicates()
}

func main() {
	rawData := []string{"  apple ", "banana", "  apple", "cherry  ", "banana", "  ", "date"}
	cleaner := NewDataCleaner(rawData)
	cleaned := cleaner.Clean()
	fmt.Println("Cleaned data:", cleaned)
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}