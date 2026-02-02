package datautils

import (
	"regexp"
	"strings"
)

func SanitizeString(input string) string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(trimmed, " ")
	
	return cleaned
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

	cleanedHeaders := make([]string, len(headers))
	for i, h := range headers {
		cleanedHeaders[i] = strings.TrimSpace(strings.ToLower(h))
	}
	if err := writer.Write(cleanedHeaders); err != nil {
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
			if cleanedField == "" {
				cleanedField = "N/A"
			}
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
}package main

import (
	"crypto/md5"
	"fmt"
	"strings"
)

type Record struct {
	ID   int
	Name string
	Data string
}

func deduplicateRecords(records []Record) []Record {
	seen := make(map[[16]byte]bool)
	var unique []Record

	for _, rec := range records {
		hash := md5.Sum([]byte(fmt.Sprintf("%d%s%s", rec.ID, rec.Name, rec.Data)))
		if !seen[hash] {
			seen[hash] = true
			unique = append(unique, rec)
		}
	}
	return unique
}

func validateRecord(rec Record) error {
	if rec.ID <= 0 {
		return fmt.Errorf("invalid ID: %d", rec.ID)
	}
	if strings.TrimSpace(rec.Name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(rec.Data) > 1000 {
		return fmt.Errorf("data exceeds maximum length")
	}
	return nil
}

func cleanData(records []Record) ([]Record, []string) {
	var cleaned []Record
	var errors []string

	unique := deduplicateRecords(records)

	for _, rec := range unique {
		if err := validateRecord(rec); err != nil {
			errors = append(errors, fmt.Sprintf("Record %d: %v", rec.ID, err))
			continue
		}
		cleaned = append(cleaned, rec)
	}
	return cleaned, errors
}

func main() {
	sample := []Record{
		{1, "Alpha", "Sample data"},
		{2, "Beta", "Another sample"},
		{1, "Alpha", "Sample data"},
		{3, "", "Invalid record"},
		{4, "Gamma", strings.Repeat("x", 2000)},
	}

	cleaned, errs := cleanData(sample)
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	fmt.Printf("Errors found: %d\n", len(errs))
	for _, e := range errs {
		fmt.Println(e)
	}
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

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}