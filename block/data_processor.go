
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type DataProcessor struct {
	InputPath  string
	OutputPath string
}

func NewDataProcessor(input, output string) *DataProcessor {
	return &DataProcessor{
		InputPath:  input,
		OutputPath: output,
	}
}

func (dp *DataProcessor) ValidateAndClean() error {
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
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	cleanedHeaders := dp.cleanHeaders(headers)
	if err := writer.Write(cleanedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("warning: skipping invalid record: %v", err)
			continue
		}

		cleanedRecord := dp.cleanRecord(record)
		if dp.isValidRecord(cleanedRecord) {
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			recordCount++
		}
	}

	log.Printf("processed %d valid records", recordCount)
	return nil
}

func (dp *DataProcessor) cleanHeaders(headers []string) []string {
	cleaned := make([]string, len(headers))
	for i, header := range headers {
		cleaned[i] = strings.TrimSpace(header)
		cleaned[i] = strings.ToLower(cleaned[i])
		cleaned[i] = strings.ReplaceAll(cleaned[i], " ", "_")
	}
	return cleaned
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		cleaned[i] = strings.TrimSpace(field)
		if cleaned[i] == "" {
			cleaned[i] = "N/A"
		}
	}
	return cleaned
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
	for _, field := range record {
		if field == "" {
			return false
		}
	}
	return true
}

func main() {
	processor := NewDataProcessor("input.csv", "output.csv")
	if err := processor.ValidateAndClean(); err != nil {
		log.Fatalf("processing failed: %v", err)
	}
	fmt.Println("data processing completed successfully")
}