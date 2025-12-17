package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func deduplicateRecords(records [][]string) [][]string {
	seen := make(map[string]bool)
	var unique [][]string

	for _, record := range records {
		key := strings.Join(record, "|")
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func normalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func normalizeRecords(records [][]string) [][]string {
	for i := range records {
		for j := range records[i] {
			records[i][j] = normalizeString(records[i][j])
		}
	}
	return records
}

func processCSV(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("empty csv file")
	}

	headers := records[0]
	data := records[1:]

	data = deduplicateRecords(data)
	data = normalizeRecords(data)

	sort.Slice(data, func(i, j int) bool {
		return strings.Join(data[i], "") < strings.Join(data[j], "")
	})

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := processCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Data cleaned successfully. Output saved to %s\n", outputFile)
}