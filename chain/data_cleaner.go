
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func cleanCSVData(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleanedRecord := make([]string, len(record))
		allEmpty := true
		for i, field := range record {
			trimmed := strings.TrimSpace(field)
			cleanedRecord[i] = trimmed
			if trimmed != "" {
				allEmpty = false
			}
		}

		if !allEmpty {
			err = writer.Write(cleanedRecord)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	err := cleanCSVData("input.csv", "output.csv")
	if err != nil {
		log.Fatal("Failed to process CSV:", err)
	}
}package main

import (
	"regexp"
	"strings"
)

func SanitizeString(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	normalized := re.ReplaceAllString(trimmed, " ")
	
	// Convert to lowercase for consistency
	lowercased := strings.ToLower(normalized)
	
	return lowercased
}

func RemoveSpecialCharacters(input string) string {
	// Keep only alphanumeric characters and spaces
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return re.ReplaceAllString(input, "")
}

func NormalizeWhitespace(input string) string {
	// Replace various whitespace characters with standard space
	re := regexp.MustCompile(`[\t\n\r\f\v]+`)
	return re.ReplaceAllString(input, " ")
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSVData(inputPath, outputPath string) error {
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
		fmt.Println("Usage: go run data_cleaner.go <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSVData(inputFile, outputFile); err != nil {
		fmt.Printf("Error cleaning data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Data cleaned successfully. Output written to %s\n", outputFile)
}package main

import (
	"fmt"
	"strings"
)

// TrimSpaces removes leading and trailing whitespace from each string in the slice.
func TrimSpaces(input []string) []string {
	trimmed := make([]string, len(input))
	for i, s := range input {
		trimmed[i] = strings.TrimSpace(s)
	}
	return trimmed
}

func main() {
	data := []string{"  apple ", "banana  ", "  cherry  ", "date"}
	cleaned := TrimSpaces(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}