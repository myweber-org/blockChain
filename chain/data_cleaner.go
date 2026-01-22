package utils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)

	// Remove any HTML/XML tags
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(trimmed, "")

	// Escape potentially dangerous characters
	re = regexp.MustCompile(`['"\\;]`)
	escaped := re.ReplaceAllStringFunc(cleaned, func(match string) string {
		return "\\" + match
	})

	return escaped
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
	data := []string{" apple ", "banana", " apple", "banana ", "  ", "cherry"}
	cleaned := CleanData(data)
	fmt.Println("Cleaned data:", cleaned)
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

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV: %w", err)
        }

        cleaned := make([]string, len(record))
        for i, field := range record {
            cleaned[i] = strings.TrimSpace(field)
        }

        if err := writer.Write(cleaned); err != nil {
            return fmt.Errorf("error writing CSV: %w", err)
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

    if err := cleanCSV(inputFile, outputFile); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully cleaned data. Output saved to %s\n", outputFile)
}