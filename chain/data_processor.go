package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func processCSVFile(inputPath string, outputPath string) error {
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

	csvReader := csv.NewReader(inputFile)
	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	headers, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	outputHeaders := []string{"ID", "Name", "ProcessedValue", "Status"}
	if err := csvWriter.Write(outputHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read row: %w", err)
		}

		if len(row) < 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			continue
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			continue
		}

		active := strings.ToLower(strings.TrimSpace(row[3])) == "true"

		record := Record{
			ID:     id,
			Name:   name,
			Value:  value,
			Active: active,
		}

		processedValue := record.Value * 1.1
		status := "inactive"
		if record.Active {
			status = "active"
		}

		outputRow := []string{
			strconv.Itoa(record.ID),
			record.Name,
			fmt.Sprintf("%.2f", processedValue),
			status,
		}

		if err := csvWriter.Write(outputRow); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}

		recordCount++
	}

	fmt.Printf("Processed %d valid records\n", recordCount)
	return nil
}

func validateFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func main() {
	inputFile := "input_data.csv"
	outputFile := "processed_data.csv"

	if !validateFileExists(inputFile) {
		fmt.Println("Input file does not exist")
		return
	}

	if err := processCSVFile(inputFile, outputFile); err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	fmt.Println("Data processing completed successfully")
}