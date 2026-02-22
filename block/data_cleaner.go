
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func MarkInvalidRecords(records []DataRecord) []DataRecord {
	for i := range records {
		records[i].Valid = ValidateEmail(records[i].Email)
	}
	return records
}

func CleanDataPipeline(records []DataRecord) []DataRecord {
	records = DeduplicateRecords(records)
	records = MarkInvalidRecords(records)
	return records
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob Wilson", "invalid-email", false},
		{5, "Alice Brown", "alice@test", false},
	}

	cleaned := CleanDataPipeline(sampleData)

	fmt.Println("Cleaned Records:")
	for _, record := range cleaned {
		status := "Valid"
		if !record.Valid {
			status = "Invalid"
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, Status: %s\n",
			record.ID, record.Name, record.Email, status)
	}
}
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

func cleanCSVData(inputPath string, outputPath string) error {
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

	headers = append(headers, "Valid")
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
			return fmt.Errorf("failed to read row: %w", err)
		}

		recordCount++
		valid := validateRecord(row)

		row = append(row, strconv.FormatBool(valid))
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}

		if valid {
			validCount++
		}
	}

	fmt.Printf("Processed %d records, %d valid, %d invalid\n",
		recordCount, validCount, recordCount-validCount)

	return nil
}

func validateRecord(row []string) bool {
	if len(row) < 4 {
		return false
	}

	if _, err := strconv.Atoi(row[0]); err != nil {
		return false
	}

	if strings.TrimSpace(row[1]) == "" {
		return false
	}

	if !strings.Contains(row[2], "@") {
		return false
	}

	if _, err := strconv.ParseFloat(row[3], 64); err != nil {
		return false
	}

	return true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSVData(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data cleaning completed successfully")
}