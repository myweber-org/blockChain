
package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNum++
		if lineNum == 1 {
			continue
		}

		if len(line) != 4 {
			return nil, errors.New("invalid column count on line " + strconv.Itoa(lineNum))
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, errors.New("invalid ID on line " + strconv.Itoa(lineNum))
		}

		name := line[1]

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, errors.New("invalid value on line " + strconv.Itoa(lineNum))
		}

		valid := line[3] == "true"

		record := DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
			Valid: valid,
		}
		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if record.Valid && record.Value > 0 && record.Name != "" {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}

	total := 0.0
	for _, record := range records {
		total += record.Value
	}
	return total / float64(len(records))
}