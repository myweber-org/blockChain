package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID      string
	Name    string
	Email   string
	Active  string
	Score   string
}

func cleanString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processCSVFile(inputPath, outputPath string) error {
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

	csvReader := csv.NewReader(inFile)
	csvWriter := csv.NewWriter(outFile)
	defer csvWriter.Flush()

	headers, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if err := csvWriter.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	validCount := 0

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		recordCount++

		if len(row) < 5 {
			continue
		}

		record := DataRecord{
			ID:     cleanString(row[0]),
			Name:   cleanString(row[1]),
			Email:  cleanString(row[2]),
			Active: cleanString(row[3]),
			Score:  cleanString(row[4]),
		}

		if record.ID == "" || record.Name == "" {
			continue
		}

		if !validateEmail(record.Email) {
			continue
		}

		if record.Active != "true" && record.Active != "false" {
			continue
		}

		outputRow := []string{
			record.ID,
			strings.Title(record.Name),
			record.Email,
			record.Active,
			record.Score,
		}

		if err := csvWriter.Write(outputRow); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}

		validCount++
	}

	fmt.Printf("Processed %d records, %d valid records written to %s\n", 
		recordCount, validCount, outputPath)
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := processCSVFile(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}