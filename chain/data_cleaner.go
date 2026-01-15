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
	ID    int
	Name  string
	Email string
	Score float64
}

func cleanCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("cannot read headers: %w", err)
	}

	headers = append(headers, "Valid")
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("cannot write headers: %w", err)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row: %w", err)
		}

		record, validationErr := validateRow(row)
		isValid := validationErr == nil

		outputRow := []string{
			strconv.Itoa(record.ID),
			strings.TrimSpace(record.Name),
			strings.ToLower(strings.TrimSpace(record.Email)),
			fmt.Sprintf("%.2f", record.Score),
			strconv.FormatBool(isValid),
		}

		if err := writer.Write(outputRow); err != nil {
			return fmt.Errorf("error writing row: %w", err)
		}
	}

	return nil
}

func validateRow(row []string) (Record, error) {
	if len(row) < 4 {
		return Record{}, fmt.Errorf("insufficient columns")
	}

	id, err := strconv.Atoi(row[0])
	if err != nil || id <= 0 {
		return Record{}, fmt.Errorf("invalid ID")
	}

	name := strings.TrimSpace(row[1])
	if name == "" {
		return Record{}, fmt.Errorf("empty name")
	}

	email := strings.TrimSpace(row[2])
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return Record{}, fmt.Errorf("invalid email format")
	}

	score, err := strconv.ParseFloat(row[3], 64)
	if err != nil || score < 0 || score > 100 {
		return Record{}, fmt.Errorf("score out of range")
	}

	return Record{
		ID:    id,
		Name:  name,
		Email: email,
		Score: score,
	}, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Data cleaning completed. Output saved to %s\n", outputFile)
}