package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ProcessRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("invalid record ID")
	}

	if strings.TrimSpace(record.Name) == "" {
		return errors.New("record name cannot be empty")
	}

	if record.Value < 0 {
		return errors.New("record value cannot be negative")
	}

	fmt.Printf("Processing record %d: %s (%.2f)\n", record.ID, record.Name, record.Value)
	return nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, error) {
	var validRecords []DataRecord
	var validationErrors []string

	for _, record := range records {
		err := ProcessRecord(record)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("Record %d: %v", record.ID, err))
			continue
		}
		validRecords = append(validRecords, record)
	}

	if len(validationErrors) > 0 {
		return validRecords, fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
	}

	return validRecords, nil
}

func CalculateTotal(records []DataRecord) float64 {
	var total float64
	for _, record := range records {
		total += record.Value
	}
	return total
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
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(trimmed, " ")
	return cleaned
}

func ValidateInput(input string) bool {
	if len(input) == 0 {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9\s.,!?-]+$`)
	return re.MatchString(input)
}

func ProcessData(input string) (string, bool) {
	cleaned := CleanInput(input)
	valid := ValidateInput(cleaned)
	if !valid {
		return "", false
	}
	return cleaned, true
}
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "regexp"
    "strings"
)

type UserProfile struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Age       int    `json:"age"`
    IsActive  bool   `json:"is_active"`
}

func ValidateUserProfile(profile UserProfile) error {
    if profile.ID <= 0 {
        return errors.New("invalid user ID")
    }

    usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
    if !usernameRegex.MatchString(profile.Username) {
        return errors.New("username must be 3-20 alphanumeric characters or underscores")
    }

    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(profile.Email) {
        return errors.New("invalid email format")
    }

    if profile.Age < 0 || profile.Age > 120 {
        return errors.New("age must be between 0 and 120")
    }

    return nil
}

func TransformProfile(profile UserProfile) UserProfile {
    transformed := profile
    transformed.Username = strings.ToLower(transformed.Username)
    transformed.Email = strings.ToLower(transformed.Email)
    transformed.IsActive = true

    if transformed.Age < 18 {
        transformed.IsActive = false
    }

    return transformed
}

func ProcessUserData(inputJSON string) (string, error) {
    var profile UserProfile
    err := json.Unmarshal([]byte(inputJSON), &profile)
    if err != nil {
        return "", fmt.Errorf("failed to parse JSON: %v", err)
    }

    err = ValidateUserProfile(profile)
    if err != nil {
        return "", fmt.Errorf("validation failed: %v", err)
    }

    transformedProfile := TransformProfile(profile)

    outputJSON, err := json.MarshalIndent(transformedProfile, "", "  ")
    if err != nil {
        return "", fmt.Errorf("failed to marshal JSON: %v", err)
    }

    return string(outputJSON), nil
}

func main() {
    sampleJSON := `{
        "id": 123,
        "username": "TestUser_123",
        "email": "TEST@EXAMPLE.COM",
        "age": 25,
        "is_active": false
    }`

    result, err := ProcessUserData(sampleJSON)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println("Processed user profile:")
    fmt.Println(result)
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ExtractDomain(email string) string {
	if !dp.ValidateEmail(email) {
		return ""
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func (dp *DataProcessor) NormalizeSpaces(input string) string {
	return dp.whitespaceRegex.ReplaceAllString(input, " ")
}package main

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
	ID        int
	Name      string
	Value     float64
	Timestamp string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)
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

func parseRow(row []string, line int) (DataRecord, error) {
	var record DataRecord
	var err error

	record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", line, err)
	}

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return DataRecord{}, fmt.Errorf("empty name at line %d", line)
	}

	record.Value, err = strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", line, err)
	}

	record.Timestamp = strings.TrimSpace(row[3])
	if record.Timestamp == "" {
		return DataRecord{}, fmt.Errorf("empty timestamp at line %d", line)
	}

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
			return fmt.Errorf("negative value %f for record ID %d", record.Value, record.ID)
		}
	}
	return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, float64) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var min, max float64
	first := true

	for _, record := range records {
		sum += record.Value
		if first {
			min = record.Value
			max = record.Value
			first = false
		} else {
			if record.Value < min {
				min = record.Value
			}
			if record.Value > max {
				max = record.Value
			}
		}
	}

	average := sum / float64(len(records))
	return average, min, max
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
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func parseCSVFile(filename string) ([]Record, error) {
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

func parseRow(row []string, lineNumber int) (Record, error) {
	var record Record

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

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

func validateRecords(records []Record) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d must be positive", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if len(record.Name) == 0 {
			return fmt.Errorf("empty name for record ID: %d", record.ID)
		}

		if record.Value < 0 {
			return fmt.Errorf("negative value for record ID: %d", record.ID)
		}
	}

	return nil
}

func calculateStatistics(records []Record) (float64, float64, int) {
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
}package main

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