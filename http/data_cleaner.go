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
	data := []string{"  apple ", "banana", "  apple", "banana ", " ", "cherry"}
	cleaned := CleanData(data)
	fmt.Println("Cleaned data:", cleaned)
}package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Deduplicate(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(item) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && len(item) < 100
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"apple", "Apple", "banana", "", "  BANANA  ", "cherry", "a" + strings.Repeat("b", 99)}
	result := cleaner.Deduplicate(data)
	fmt.Println("Cleaned data:", result)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Deduplicate(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(normalized) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && len(item) < 100
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"apple",
		"Apple",
		"banana",
		"  banana  ",
		"",
		"cherry",
		"cherry",
	}
	
	cleaned := cleaner.Deduplicate(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
	fmt.Printf("Unique count: %d\n", len(cleaned))
	
	cleaner.Reset()
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataCleaner struct {
	skipHeader bool
	delimiter  rune
}

func NewDataCleaner(skipHeader bool, delimiter rune) *DataCleaner {
	return &DataCleaner{
		skipHeader: skipHeader,
		delimiter:  delimiter,
	}
}

func (dc *DataCleaner) CleanCSV(inputPath, outputPath string) error {
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
	reader.Comma = dc.delimiter
	writer := csv.NewWriter(outFile)
	writer.Comma = dc.delimiter
	defer writer.Flush()

	lineNum := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read error at line %d: %w", lineNum+1, err)
		}

		lineNum++
		if dc.skipHeader && lineNum == 1 {
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("write header error: %w", err)
			}
			continue
		}

		cleaned := dc.cleanRecord(record)
		if cleaned == nil {
			continue
		}

		if err := writer.Write(cleaned); err != nil {
			return fmt.Errorf("write error at line %d: %w", lineNum, err)
		}
	}

	return nil
}

func (dc *DataCleaner) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	allEmpty := true

	for i, field := range record {
		field = strings.TrimSpace(field)
		field = strings.ToLower(field)
		if field == "" || field == "null" || field == "n/a" {
			field = "unknown"
		}
		cleaned[i] = field
		if field != "unknown" {
			allEmpty = false
		}
	}

	if allEmpty {
		return nil
	}
	return cleaned
}

func main() {
	cleaner := NewDataCleaner(true, ',')
	err := cleaner.CleanCSV("input.csv", "cleaned.csv")
	if err != nil {
		fmt.Printf("Error cleaning data: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Data cleaning completed successfully")
}