package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ParseUserData attempts to parse JSON data into a map.
func ParseUserData(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := `{"name": "Alice", "age": 30, "active": true}`

	valid, err := ValidateJSON([]byte(sampleJSON))
	if err != nil {
		log.Printf("Validation error: %v", err)
	} else {
		fmt.Println("JSON is valid:", valid)
	}

	userData, err := ParseUserData(sampleJSON)
	if err != nil {
		log.Printf("Parse error: %v", err)
	} else {
		fmt.Printf("Parsed user data: %+v\n", userData)
	}
}
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Username string
	Email    string
	Age      int
}

func ValidateProfile(profile UserProfile) error {
	if strings.TrimSpace(profile.Username) == "" {
		return errors.New("username cannot be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func NormalizeProfile(profile *UserProfile) {
	profile.Username = strings.TrimSpace(profile.Username)
	profile.Email = strings.ToLower(strings.TrimSpace(profile.Email))
}

func ProcessUserData(profile UserProfile) (UserProfile, error) {
	NormalizeProfile(&profile)

	if err := ValidateProfile(profile); err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
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

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func CalculateStatistics(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
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

	avg, max := CalculateStatistics(records)
	fmt.Printf("Processed %d records successfully\n", len(records))
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
}
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

type DataProcessor struct {
	records []DataRecord
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		records: make([]DataRecord, 0),
	}
}

func (dp *DataProcessor) LoadCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	lineNumber := 0
	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if lineNumber == 1 {
			continue
		}

		record, err := parseRow(row)
		if err != nil {
			return fmt.Errorf("parse error at line %d: %w", lineNumber, err)
		}

		dp.records = append(dp.records, record)
	}

	return nil
}

func parseRow(row []string) (DataRecord, error) {
	if len(row) != 4 {
		return DataRecord{}, errors.New("invalid number of columns")
	}

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid ID: %w", err)
	}

	name := strings.TrimSpace(row[1])
	if name == "" {
		return DataRecord{}, errors.New("name cannot be empty")
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid value: %w", err)
	}

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid active flag: %w", err)
	}

	return DataRecord{
		ID:     id,
		Name:   name,
		Value:  value,
		Active: active,
	}, nil
}

func (dp *DataProcessor) FilterActive() []DataRecord {
	var active []DataRecord
	for _, record := range dp.records {
		if record.Active {
			active = append(active, record)
		}
	}
	return active
}

func (dp *DataProcessor) CalculateAverage() float64 {
	if len(dp.records) == 0 {
		return 0
	}

	var sum float64
	for _, record := range dp.records {
		sum += record.Value
	}
	return sum / float64(len(dp.records))
}

func (dp *DataProcessor) FindByName(name string) *DataRecord {
	for _, record := range dp.records {
		if strings.EqualFold(record.Name, name) {
			return &record
		}
	}
	return nil
}

func (dp *DataProcessor) Validate() []error {
	var errors []error

	seenIDs := make(map[int]bool)
	for i, record := range dp.records {
		if seenIDs[record.ID] {
			errors = append(errors, fmt.Errorf("duplicate ID %d at record %d", record.ID, i+1))
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			errors = append(errors, fmt.Errorf("negative value %f for ID %d", record.Value, record.ID))
		}
	}

	return errors
}

func (dp *DataProcessor) GetStats() map[string]interface{} {
	active := dp.FilterActive()
	average := dp.CalculateAverage()

	return map[string]interface{}{
		"total_records":   len(dp.records),
		"active_records":  len(active),
		"average_value":   average,
		"validation_errors": len(dp.Validate()),
	}
}