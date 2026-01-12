package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	emailRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) SanitizeString(input string) string {
	return strings.TrimSpace(input)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, string, bool) {
	sanitizedName := dp.SanitizeString(name)
	sanitizedEmail := dp.SanitizeString(email)
	isValidEmail := dp.ValidateEmail(sanitizedEmail)
	return sanitizedName, sanitizedEmail, isValidEmail
}
package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Value int
	Tags  []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if len(record.Tags) == 0 {
		return errors.New("record must have at least one tag")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := DataRecord{
		ID:    strings.ToUpper(record.ID),
		Value: record.Value * 2,
		Tags:  make([]string, len(record.Tags)),
	}

	for i, tag := range record.Tags {
		transformed.Tags[i] = strings.TrimSpace(tag)
	}

	return transformed
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord

	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}

		transformed := TransformRecord(record)
		processed = append(processed, transformed)
	}

	return processed, nil
}

func main() {
	records := []DataRecord{
		{ID: "abc123", Value: 42, Tags: []string{"test", "sample"}},
		{ID: "def456", Value: 100, Tags: []string{"production", "data"}},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Printf("Processed %d records successfully\n", len(processed))
	for _, record := range processed {
		fmt.Printf("ID: %s, Value: %d, Tags: %v\n", record.ID, record.Value, record.Tags)
	}
}package main

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
	return strings.TrimSpace(username)
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
	rawJSON := `{"email":"test@example.com","username":"  john_doe  ","age":25}`
	processedData, err := ProcessUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
}