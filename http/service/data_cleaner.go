
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Score float64
}

func cleanCSVData(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	validCount := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row: %w", err)
		}

		recordCount++
		cleanedRow, isValid := validateAndCleanRow(row)

		if isValid {
			if err := writer.Write(cleanedRow); err != nil {
				return fmt.Errorf("error writing row: %w", err)
			}
			validCount++
		}
	}

	fmt.Printf("Processed %d records, %d valid records written to %s\n", recordCount, validCount, outputPath)
	return nil
}

func validateAndCleanRow(row []string) ([]string, bool) {
	if len(row) != 4 {
		return nil, false
	}

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil || id <= 0 {
		return nil, false
	}

	name := strings.TrimSpace(row[1])
	if name == "" || len(name) > 100 {
		return nil, false
	}

	email := strings.TrimSpace(row[2])
	if !strings.Contains(email, "@") || strings.Contains(email, " ") {
		return nil, false
	}

	score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
	if err != nil || score < 0 || score > 100 {
		return nil, false
	}

	return []string{
		strconv.Itoa(id),
		name,
		strings.ToLower(email),
		fmt.Sprintf("%.2f", score),
	}, true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run data_cleaner.go <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSVData(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
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

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) ProcessBatch(items []string) []string {
	var unique []string
	for _, item := range items {
		if !dc.IsDuplicate(item) {
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", "apple ", " BANANA", "banana", "Cherry"}
	
	fmt.Println("Original data:", data)
	
	cleaned := cleaner.ProcessBatch(data)
	fmt.Println("Cleaned data:", cleaned)
	
	cleaner.Reset()
	
	moreData := []string{"grape", "Grape", "PEACH"}
	secondBatch := cleaner.ProcessBatch(moreData)
	fmt.Println("Second batch:", secondBatch)
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func removeDuplicates(inputFile, outputFile string) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	records := make(map[string]bool)
	reader := csv.NewReader(in)

	var cleanedData [][]string
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	cleanedData = append(cleanedData, headers)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		key := fmt.Sprintf("%v", record)
		if !records[key] {
			records[key] = true
			cleanedData = append(cleanedData, record)
		}
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	return writer.WriteAll(cleanedData)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	err := removeDuplicates(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Duplicate removal completed successfully")
}