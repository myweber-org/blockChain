
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}

	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}

	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Email:    strings.TrimSpace(email),
		Username: transformedUsername,
		Age:      age,
	}

	err := ValidateUserData(userData)
	if err != nil {
		return UserData{}, err
	}

	return userData, nil
}
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
	ID        int
	Name      string
	Value     float64
	Validated bool
}

func ProcessCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) < 3 {
			continue
		}

		record, err := parseRecord(row)
		if err != nil {
			fmt.Printf("Warning: line %d - %v\n", lineNumber, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRecord(row []string) (Record, error) {
	var record Record

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID format: %s", row[0])
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value format: %s", row[2])
	}
	record.Value = value

	record.Validated = validateRecord(record)

	return record, nil
}

func validateRecord(r Record) bool {
	if r.ID <= 0 {
		return false
	}
	if len(r.Name) == 0 {
		return false
	}
	if r.Value < 0 {
		return false
	}
	return true
}

func CalculateStatistics(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var validCount int

	for _, record := range records {
		if record.Validated {
			sum += record.Value
			validCount++
		}
	}

	if validCount == 0 {
		return 0, 0, 0
	}

	average := sum / float64(validCount)
	return sum, average, validCount
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ProcessCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d records\n", len(records))

	total, average, validCount := CalculateStatistics(records)
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Total value: %.2f\n", total)
	fmt.Printf("Average value: %.2f\n", average)
}