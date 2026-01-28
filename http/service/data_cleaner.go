package main

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

func cleanCSVData(inputPath, outputPath string) error {
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

	seen := make(map[string]bool)
	headers, err := reader.Read()
	if err != nil {
		return err
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		key := strings.Join(record, "|")
		if !seen[key] {
			seen[key] = true
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}
	if err := cleanCSVData(os.Args[1], os.Args[2]); err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}
	println("Data cleaning completed successfully")
}