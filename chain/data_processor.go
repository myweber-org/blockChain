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
	if strings.TrimSpace(data.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}
	return nil
}

func NormalizeUserData(data UserData) UserData {
	return UserData{
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Username: strings.TrimSpace(data.Username),
		Age:      data.Age,
	}
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	data := UserData{
		Email:    email,
		Username: username,
		Age:      age,
	}
	
	normalizedData := NormalizeUserData(data)
	
	if err := ValidateUserData(normalizedData); err != nil {
		return UserData{}, err
	}
	
	return normalizedData, nil
}package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Email     string
	Username  string
	Age       int
	Biography string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserProfile(profile UserProfile) error {
	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if strings.TrimSpace(profile.Username) == "" {
		return errors.New("username cannot be empty")
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	if profile.Age < 0 || profile.Age > 120 {
		return errors.New("age must be between 0 and 120")
	}

	if len(profile.Biography) > 500 {
		return errors.New("biography cannot exceed 500 characters")
	}

	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(strings.TrimSpace(profile.Username))
	transformed.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	transformed.Biography = strings.TrimSpace(profile.Biography)

	if transformed.Biography == "" {
		transformed.Biography = "No biography provided."
	}

	return transformed
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateUserProfile(profile); err != nil {
		return UserProfile{}, err
	}

	return TransformProfile(profile), nil
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

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
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

		if len(row) < 4 {
			continue
		}

		record, err := parseRecord(row)
		if err != nil {
			fmt.Printf("Warning: skipping invalid record at line %d: %v\n", lineNumber, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRecord(row []string) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID: %w", err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value: %w", err)
	}
	record.Value = value

	valid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid valid flag: %w", err)
	}
	record.Valid = valid

	return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}

	var sum float64
	count := 0

	for _, record := range records {
		if record.Valid {
			sum += record.Value
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return sum / float64(count)
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func readCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}
	return records, nil
}

func validateCSVData(records [][]string) error {
	if len(records) == 0 {
		return fmt.Errorf("CSV file is empty")
	}

	header := records[0]
	expectedColumns := 3
	if len(header) != expectedColumns {
		return fmt.Errorf("invalid header length: expected %d, got %d", expectedColumns, len(header))
	}

	for i, record := range records[1:] {
		if len(record) != expectedColumns {
			return fmt.Errorf("row %d: invalid column count", i+2)
		}
		for j, field := range record {
			if strings.TrimSpace(field) == "" {
				return fmt.Errorf("row %d, column %d: empty field", i+2, j+1)
			}
		}
	}
	return nil
}

func processCSVData(records [][]string) []map[string]string {
	if len(records) < 2 {
		return []map[string]string{}
	}

	header := records[0]
	var result []map[string]string

	for _, row := range records[1:] {
		rowMap := make(map[string]string)
		for j, value := range row {
			if j < len(header) {
				rowMap[header[j]] = strings.TrimSpace(value)
			}
		}
		result = append(result, rowMap)
	}
	return result
}

func writeProcessedData(data []map[string]string, outputPath string) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to write")
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var headers []string
	for key := range data[0] {
		headers = append(headers, key)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for _, record := range data {
		var row []string
		for _, header := range headers {
			row = append(row, record[header])
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}
	return nil
}
package main

import (
	"errors"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return errors.New("invalid user ID")
	}
	if strings.TrimSpace(data.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func TransformUserName(data *UserData) {
	data.Name = strings.ToUpper(strings.TrimSpace(data.Name))
}

func ProcessUserInput(data UserData) (UserData, error) {
	if err := ValidateUserData(data); err != nil {
		return UserData{}, err
	}
	TransformUserName(&data)
	return data, nil
}