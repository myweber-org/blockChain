package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	header, err := reader.Read()
	if err != nil {
		return err
	}

	seen := make(map[string]bool)
	writer.Write(header)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for i, val := range record {
			record[i] = strings.TrimSpace(val)
		}

		key := strings.Join(record, "|")
		if !seen[key] {
			seen[key] = true
			writer.Write(record)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	err := cleanCSV(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data cleaning completed successfully")
}package main

import (
	"encoding/csv"
	"io"
	"strings"
)

func CleanCSVData(input io.Reader, output io.Writer) error {
	reader := csv.NewReader(input)
	writer := csv.NewWriter(output)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleaned := make([]string, 0, len(record))
		hasData := false

		for _, field := range record {
			trimmed := strings.TrimSpace(field)
			cleaned = append(cleaned, trimmed)
			if trimmed != "" {
				hasData = true
			}
		}

		if hasData {
			if err := writer.Write(cleaned); err != nil {
				return err
			}
		}
	}

	return nil
}