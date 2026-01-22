package data_processor

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID        int
	Name      string
	Value     float64
	Validated bool
}

func ParseCSVData(reader io.Reader) ([]DataRecord, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRecord
	for i, row := range records {
		if len(row) < 3 {
			return nil, errors.New("insufficient columns in row " + strconv.Itoa(i))
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, errors.New("invalid ID in row " + strconv.Itoa(i))
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, errors.New("empty name in row " + strconv.Itoa(i))
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, errors.New("invalid value in row " + strconv.Itoa(i))
		}

		validated := false
		if len(row) > 3 {
			validated = strings.ToLower(strings.TrimSpace(row[3])) == "true"
		}

		data = append(data, DataRecord{
			ID:        id,
			Name:      name,
			Value:     value,
			Validated: validated,
		})
	}

	return data, nil
}

func ValidateRecords(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		if record.ID > 0 && record.Value >= 0 {
			record.Validated = true
			validated = append(validated, record)
		}
	}
	return validated
}

func CalculateTotal(records []DataRecord) float64 {
	var total float64
	for _, record := range records {
		if record.Validated {
			total += record.Value
		}
	}
	return total
}