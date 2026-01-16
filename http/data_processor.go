
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, err
	}

	transformed := DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp.UTC(),
		Tags:      append([]string{}, record.Tags...),
	}

	if len(transformed.Tags) == 0 {
		transformed.Tags = append(transformed.Tags, "default")
	}

	return transformed, nil
}

func ProcessBatch(records []DataRecord, multiplier float64) ([]DataRecord, []error) {
	var processed []DataRecord
	var errs []error

	for i, record := range records {
		transformed, err := TransformRecord(record, multiplier)
		if err != nil {
			errs = append(errs, fmt.Errorf("record %d: %w", i, err))
			continue
		}
		processed = append(processed, transformed)
	}

	return processed, errs
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records provided")
	}

	var sum float64
	var count int

	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			continue
		}
		sum += record.Value
		count++
	}

	if count == 0 {
		return 0, 0, errors.New("no valid records found")
	}

	average := sum / float64(count)

	var varianceSum float64
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			continue
		}
		diff := record.Value - average
		varianceSum += diff * diff
	}

	variance := varianceSum / float64(count)

	return average, variance, nil
}package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return fmt.Errorf("invalid user ID: must be positive integer")
	}
	if strings.TrimSpace(data.Name) == "" {
		return fmt.Errorf("user name cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func TransformUserName(data UserData) UserData {
	data.Name = strings.ToUpper(strings.TrimSpace(data.Name))
	return data
}

func ProcessUserInput(data UserData) (UserData, error) {
	if err := ValidateUserData(data); err != nil {
		return UserData{}, err
	}
	transformedData := TransformUserName(data)
	return transformedData, nil
}

func main() {
	testData := UserData{
		ID:    1001,
		Name:  "  john doe  ",
		Email: "john@example.com",
	}

	result, err := ProcessUserInput(testData)
	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		return
	}

	fmt.Printf("Processed user: ID=%d, Name='%s', Email='%s'\n",
		result.ID, result.Name, result.Email)
}
package data_processor

import (
	"errors"
	"regexp"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Value float64
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !emailRegex.MatchString(record.Email) {
		return errors.New("invalid email format")
	}
	if record.Value < 0 || record.Value > 10000 {
		return errors.New("value must be between 0 and 10000")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TransformRecords(records []DataRecord) ([]DataRecord, error) {
	var validRecords []DataRecord
	for _, record := range records {
		record.Email = NormalizeEmail(record.Email)
		if err := ValidateRecord(record); err != nil {
			return nil, err
		}
		validRecords = append(validRecords, record)
	}
	return validRecords, nil
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
}
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateAndParseUser(data []byte) (*User, error) {
	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if user.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	return &user, nil
}

func main() {
	jsonData := []byte(`{"id": 123, "name": "John Doe", "email": "john@example.com"}`)

	user, err := ValidateAndParseUser(jsonData)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Parsed user: %+v\n", user)
}