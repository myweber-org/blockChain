
package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ReadCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}

	// Skip header
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			return nil, errors.New("invalid row format")
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, err
		}

		record := Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		}
		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	for _, record := range records {
		if record.ID <= 0 {
			return errors.New("invalid ID: must be positive")
		}
		if record.Name == "" {
			return errors.New("invalid Name: cannot be empty")
		}
		if record.Value < 0 {
			return errors.New("invalid Value: cannot be negative")
		}
	}
	return nil
}

func FilterByValue(records []Record, minValue float64) []Record {
	var filtered []Record
	for _, record := range records {
		if record.Value >= minValue {
			filtered = append(filtered, record)
		}
	}
	return filtered
}package main

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func TransformPhoneNumber(phone string) (string, error) {
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	if len(cleaned) < 10 {
		return "", errors.New("phone number too short")
	}
	if len(cleaned) > 15 {
		return "", errors.New("phone number too long")
	}
	return "+" + cleaned, nil
}

func SanitizeInput(input string) string {
	re := regexp.MustCompile(`[<>"'&]`)
	return re.ReplaceAllString(input, "")
}