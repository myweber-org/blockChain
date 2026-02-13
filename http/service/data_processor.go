
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

type DataRecord struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
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

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, errors.New("no valid records found in file")
	}

	return records, nil
}

func parseRow(row []string, lineNumber int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNumber)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
	}
	record.Active = active

	return record, nil
}

func ValidateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID %d found", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("record %d has negative value: %f", record.ID, record.Value)
		}
	}

	return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	minValue := records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value < minValue {
			minValue = record.Value
		}
		if record.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return average, minValue, activeCount
}
package data_processor

import (
	"regexp"
	"strings"
	"unicode"
)

type DataCleaner struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dc *DataCleaner) NormalizeString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := dc.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return normalized
}

func (dc *DataCleaner) RemoveSpecialChars(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, input)
}

func (dc *DataCleaner) ProcessInput(rawInput string) (string, error) {
	if rawInput == "" {
		return "", nil
	}
	
	cleaned := dc.NormalizeString(rawInput)
	final := dc.RemoveSpecialChars(cleaned)
	
	return final, nil
}package main

import (
    "errors"
    "strings"
)

type UserData struct {
    Username string
    Email    string
    Age      int
}

func ValidateUserData(data UserData) error {
    if strings.TrimSpace(data.Username) == "" {
        return errors.New("username cannot be empty")
    }
    if !strings.Contains(data.Email, "@") {
        return errors.New("invalid email format")
    }
    if data.Age < 0 || data.Age > 150 {
        return errors.New("age must be between 0 and 150")
    }
    return nil
}

func TransformUsername(username string) string {
    return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(rawUsername, rawEmail string, rawAge int) (UserData, error) {
    user := UserData{
        Username: TransformUsername(rawUsername),
        Email:    strings.TrimSpace(rawEmail),
        Age:      rawAge,
    }
    
    if err := ValidateUserData(user); err != nil {
        return UserData{}, err
    }
    
    return user, nil
}