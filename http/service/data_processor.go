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

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  JohnDoe  ","age":25}`
	processedData, err := ProcessUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Processed Data: %+v\n", processedData)
}
package main

import (
	"strings"
	"unicode"
)

// DataProcessor handles cleaning and normalization of string data
type DataProcessor struct{}

// CleanInput removes extra whitespace and trims the input string
func (dp *DataProcessor) CleanInput(input string) string {
	return strings.TrimSpace(input)
}

// NormalizeSpaces collapses multiple consecutive spaces into a single space
func (dp *DataProcessor) NormalizeSpaces(input string) string {
	var result strings.Builder
	prevSpace := false

	for _, r := range input {
		if unicode.IsSpace(r) {
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}

	return result.String()
}

// Process combines cleaning and normalization operations
func (dp *DataProcessor) Process(input string) string {
	cleaned := dp.CleanInput(input)
	return dp.NormalizeSpaces(cleaned)
}

func main() {
	processor := &DataProcessor{}
	
	testInput := "   Hello    World   \t\n"
	processed := processor.Process(testInput)
	
	println("Original:", testInput)
	println("Processed:", processed)
}package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	if len(records) == 0 {
		return nil, errors.New("no records to process")
	}

	var processed []DataRecord
	var errorsFound []string

	for _, record := range records {
		if !record.Valid {
			errorsFound = append(errorsFound, fmt.Sprintf("record %d is invalid", record.ID))
			continue
		}

		if record.Value < 0 {
			errorsFound = append(errorsFound, fmt.Sprintf("record %d has negative value", record.ID))
			continue
		}

		if strings.TrimSpace(record.Name) == "" {
			errorsFound = append(errorsFound, fmt.Sprintf("record %d has empty name", record.ID))
			continue
		}

		processed = append(processed, record)
	}

	if len(errorsFound) > 0 {
		return processed, fmt.Errorf("processing completed with errors: %s", strings.Join(errorsFound, "; "))
	}

	return processed, nil
}

func main() {
	records := []DataRecord{
		{ID: 1, Name: "Record A", Value: 100.5, Valid: true},
		{ID: 2, Name: "", Value: 50.0, Valid: true},
		{ID: 3, Name: "Record C", Value: -10.0, Valid: true},
		{ID: 4, Name: "Record D", Value: 75.3, Valid: false},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Successfully processed %d records\n", len(processed))
	for _, record := range processed {
		fmt.Printf("ID: %d, Name: %s, Value: %.2f\n", record.ID, record.Name, record.Value)
	}
}