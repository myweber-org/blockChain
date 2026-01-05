package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID      int
	Name    string
	Email   string
	Active  bool
	Score   float64
}

func ParseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(row) != 5 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 5, got %d", lineNum, len(row))
		}

		record, err := parseRow(row, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNum int) (Record, error) {
	var record Record
	var err error

	record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return Record{}, fmt.Errorf("empty name at line %d", lineNum)
	}

	record.Email = strings.TrimSpace(row[2])
	if !strings.Contains(record.Email, "@") {
		return Record{}, fmt.Errorf("invalid email format at line %d", lineNum)
	}

	record.Active, err = strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid active flag at line %d: %w", lineNum, err)
	}

	record.Score, err = strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid score at line %d: %w", lineNum, err)
	}

	if record.Score < 0 || record.Score > 100 {
		return Record{}, fmt.Errorf("score out of range at line %d: must be between 0 and 100", lineNum)
	}

	return record, nil
}

func ValidateRecords(records []Record) error {
	if len(records) == 0 {
		return errors.New("no records to validate")
	}

	emailSet := make(map[string]bool)
	idSet := make(map[int]bool)

	for _, record := range records {
		if emailSet[record.Email] {
			return fmt.Errorf("duplicate email found: %s", record.Email)
		}
		emailSet[record.Email] = true

		if idSet[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		idSet[record.ID] = true
	}

	return nil
}

func CalculateStatistics(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var maxScore float64
	activeCount := 0

	for _, record := range records {
		sum += record.Score
		if record.Score > maxScore {
			maxScore = record.Score
		}
		if record.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return average, maxScore, activeCount
}