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
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateAndParseJSON(rawData []byte) (*UserData, error) {
	var user UserData
	err := json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if user.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	return &user, nil
}

func main() {
	jsonStr := `{"id": 101, "name": "Alice", "email": "alice@example.com"}`
	user, err := ValidateAndParseJSON([]byte(jsonStr))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Parsed user: %+v\n", user)
}