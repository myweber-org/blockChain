
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	InputPath  string
	OutputPath string
	Delimiter  rune
}

func NewDataProcessor(input, output string) *DataProcessor {
	return &DataProcessor{
		InputPath:  input,
		OutputPath: output,
		Delimiter:  ',',
	}
}

func (dp *DataProcessor) ValidateRow(row []string) bool {
	if len(row) == 0 {
		return false
	}
	for _, field := range row {
		if strings.TrimSpace(field) == "" {
			return false
		}
	}
	return true
}

func (dp *DataProcessor) CleanField(field string) string {
	cleaned := strings.TrimSpace(field)
	cleaned = strings.ToLower(cleaned)
	return cleaned
}

func (dp *DataProcessor) Process() error {
	inputFile, err := os.Open(dp.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dp.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	reader.Comma = dp.Delimiter

	writer := csv.NewWriter(outputFile)
	writer.Comma = dp.Delimiter
	defer writer.Flush()

	headerWritten := false
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		if !headerWritten {
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing header: %w", err)
			}
			headerWritten = true
			continue
		}

		if !dp.ValidateRow(record) {
			continue
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedRecord[i] = dp.CleanField(field)
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1], os.Args[2])
	if err := processor.Process(); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data processing completed successfully")
}