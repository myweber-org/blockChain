
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID      string
	Name    string
	Email   string
	Active  string
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

		if lineNumber == 1 {
			continue
		}

		if len(row) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if record.ID == "" || record.Name == "" {
			return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
	var active []DataRecord
	for _, record := range records {
		if strings.ToLower(record.Active) == "true" || record.Active == "1" {
			active = append(active, record)
		}
	}
	return active
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	activeRecords := FilterActiveRecords(records)
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Active records: %d\n", len(activeRecords))

	for i, record := range activeRecords {
		if !ValidateEmail(record.Email) {
			fmt.Printf("Warning: Invalid email for record %d (ID: %s)\n", i+1, record.ID)
		}
	}
}package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Value string
	Valid bool
}

func ProcessRecords(records []DataRecord) ([]string, error) {
	if len(records) == 0 {
		return nil, errors.New("no records to process")
	}

	var processed []string
	for _, record := range records {
		if !record.Valid {
			continue
		}

		trimmed := strings.TrimSpace(record.Value)
		if trimmed == "" {
			continue
		}

		processed = append(processed, fmt.Sprintf("ID:%d|Value:%s", record.ID, trimmed))
	}

	if len(processed) == 0 {
		return nil, errors.New("no valid records found")
	}

	return processed, nil
}

func ValidateRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("invalid ID")
	}
	if strings.TrimSpace(record.Value) == "" {
		return errors.New("empty value")
	}
	return nil
}

func main() {
	records := []DataRecord{
		{ID: 1, Value: "alpha", Valid: true},
		{ID: 2, Value: "  ", Valid: true},
		{ID: 3, Value: "beta", Valid: false},
		{ID: 4, Value: "gamma", Valid: true},
	}

	result, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Println("Processed records:")
	for _, item := range result {
		fmt.Println(item)
	}
}
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func sanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func transformUserData(rawData []byte) (*UserData, error) {
	var user UserData
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	user.Username = sanitizeUsername(user.Username)

	if !validateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format: %s", user.Email)
	}

	if user.Age < 0 || user.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", user.Age)
	}

	return &user, nil
}

func processInput(data []byte) {
	user, err := transformUserData(data)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Printf("Processed user: %+v\n", user)
}

func main() {
	testData := []byte(`{"email":"test@example.com","username":"  john_doe  ","age":25}`)
	processInput(testData)
}
package main

import "fmt"

func FilterAndDoubleEvenNumbers(numbers []int) []int {
    var result []int
    for _, num := range numbers {
        if num%2 == 0 {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    output := FilterAndDoubleEvenNumbers(input)
    fmt.Printf("Input: %v\n", input)
    fmt.Printf("Output: %v\n", output)
}package main

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

func NormalizeUserData(data *UserData) {
	data.Email = strings.ToLower(strings.TrimSpace(data.Email))
	data.Username = strings.TrimSpace(data.Username)
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

func TransformUsername(data UserData) UserData {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
	return data
}

func ProcessUserInput(rawUsername string, rawEmail string, rawAge int) (UserData, error) {
	userData := UserData{
		Username: rawUsername,
		Email:    rawEmail,
		Age:      rawAge,
	}

	userData = TransformUsername(userData)

	if err := ValidateUserData(userData); err != nil {
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

type DataRecord struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
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

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNumber int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return DataRecord{}, fmt.Errorf("empty name at line %d", lineNumber)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
	}
	record.Active = active

	return record, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	var minValue, maxValue float64

	for i, record := range records {
		sum += record.Value

		if record.Active {
			activeCount++
		}

		if i == 0 {
			minValue = record.Value
			maxValue = record.Value
		} else {
			if record.Value < minValue {
				minValue = record.Value
			}
			if record.Value > maxValue {
				maxValue = record.Value
			}
		}
	}

	average := sum / float64(len(records))
	return average, maxValue - minValue, activeCount
}

func FilterRecords(records []DataRecord, minValue float64, activeOnly bool) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Value < minValue {
			continue
		}
		if activeOnly && !record.Active {
			continue
		}
		filtered = append(filtered, record)
	}
	return filtered
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

	fmt.Printf("Successfully processed %d records\n", len(records))

	average, rangeValue, activeCount := CalculateStatistics(records)
	fmt.Printf("Statistics - Average: %.2f, Range: %.2f, Active records: %d\n", average, rangeValue, activeCount)

	filtered := FilterRecords(records, 50.0, true)
	fmt.Printf("Filtered records (value >= 50.0 and active): %d\n", len(filtered))

	for _, record := range filtered {
		fmt.Printf("ID: %d, Name: %s, Value: %.2f, Active: %v\n", record.ID, record.Name, record.Value, record.Active)
	}
}