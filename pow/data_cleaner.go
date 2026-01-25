package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
)

func removeDuplicates(records [][]string) [][]string {
	seen := make(map[string]bool)
	var result [][]string

	for _, record := range records {
		key := fmt.Sprintf("%v", record)
		if !seen[key] {
			seen[key] = true
			result = append(result, record)
		}
	}
	return result
}

func sortRecords(records [][]string) {
	sort.Slice(records, func(i, j int) bool {
		return records[i][0] < records[j][0]
	})
}

func processCSV(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("empty CSV file")
	}

	headers := records[0]
	data := records[1:]

	cleanedData := removeDuplicates(data)
	sortRecords(cleanedData)

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

	return writer.WriteAll(cleanedData)
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

	fmt.Printf("Successfully cleaned data. Output saved to %s\n", outputFile)
}