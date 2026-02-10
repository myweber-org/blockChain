
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func cleanCSVData(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleanedRecord := make([]string, len(record))
		allEmpty := true
		for i, field := range record {
			trimmed := strings.TrimSpace(field)
			cleanedRecord[i] = trimmed
			if trimmed != "" {
				allEmpty = false
			}
		}

		if !allEmpty {
			err = writer.Write(cleanedRecord)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	err := cleanCSVData("input.csv", "output.csv")
	if err != nil {
		log.Fatal("Failed to process CSV:", err)
	}
}