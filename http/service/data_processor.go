
package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func readCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func filterRows(records [][]string, filterFunc func([]string) bool) [][]string {
	var filtered [][]string
	for _, row := range records {
		if filterFunc(row) {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func containsKeyword(row []string, keyword string) bool {
	for _, field := range row {
		if strings.Contains(strings.ToLower(field), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func writeCSVFile(filePath string, records [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func processCSVData(inputPath, outputPath, keyword string) error {
	records, err := readCSVFile(inputPath)
	if err != nil {
		return err
	}

	filtered := filterRows(records, func(row []string) bool {
		return containsKeyword(row, keyword)
	})

	if len(filtered) == 0 {
		log.Println("No matching records found")
	}

	return writeCSVFile(outputPath, filtered)
}