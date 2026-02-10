
package data_processor

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

type RecordValidator func([]string) error

func ProcessCSVData(input string, validator RecordValidator) ([][]string, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.TrimLeadingSpace = true

	var processedRecords [][]string
	lineNumber := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNumber++
		if validator != nil {
			if err := validator(record); err != nil {
				return nil, errors.New("validation failed at line " + string(rune(lineNumber)) + ": " + err.Error())
			}
		}

		processedRecords = append(processedRecords, record)
	}

	if len(processedRecords) == 0 {
		return nil, errors.New("no valid records found in CSV data")
	}

	return processedRecords, nil
}

func ValidateRecordLength(expected int) RecordValidator {
	return func(record []string) error {
		if len(record) != expected {
			return errors.New("record length mismatch")
		}
		return nil
	}
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
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

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNum, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNum)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func CalculateStats(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	min := records[0].Value
	max := records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value < min {
			min = record.Value
		}
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max - min, len(records)
}
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type UserData struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func TransformName(name string) string {
	return strings.ToUpper(strings.TrimSpace(name))
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var user UserData
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	if !ValidateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format: %s", user.Email)
	}

	user.Name = TransformName(user.Name)
	user.CreatedAt = time.Now().UTC()

	return &user, nil
}

func FilterActiveUsers(users []UserData) []UserData {
	var activeUsers []UserData
	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers
}

func GenerateReport(users []UserData) string {
	activeUsers := FilterActiveUsers(users)
	return fmt.Sprintf("Total users: %d, Active users: %d", len(users), len(activeUsers))
}