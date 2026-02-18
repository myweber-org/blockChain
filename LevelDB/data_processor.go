package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
	Tags      []string `json:"tags"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func FilterInactiveUsers(users []UserProfile) []UserProfile {
	var activeUsers []UserProfile
	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers
}

func TransformUserData(users []UserProfile) ([]map[string]interface{}, error) {
	var transformed []map[string]interface{}
	
	for _, user := range users {
		if !ValidateEmail(user.Email) {
			return nil, fmt.Errorf("invalid email for user %d", user.ID)
		}
		
		data := map[string]interface{}{
			"user_id":   user.ID,
			"username":  NormalizeUsername(user.Username),
			"email":     strings.ToLower(user.Email),
			"age_group": categorizeAge(user.Age),
			"tag_count": len(user.Tags),
			"status":    "active",
		}
		
		if !user.Active {
			data["status"] = "inactive"
		}
		
		transformed = append(transformed, data)
	}
	
	return transformed, nil
}

func categorizeAge(age int) string {
	switch {
	case age < 18:
		return "minor"
	case age >= 18 && age <= 35:
		return "young_adult"
	case age > 35 && age <= 60:
		return "adult"
	default:
		return "senior"
	}
}

func ProcessUserJSON(jsonData []byte) ([]map[string]interface{}, error) {
	var users []UserProfile
	
	if err := json.Unmarshal(jsonData, &users); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}
	
	activeUsers := FilterInactiveUsers(users)
	return TransformUserData(activeUsers)
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
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
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

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNum, len(row))
		}

		record, err := parseRow(row, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNum int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNum)
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}
	record.Value = value

	validStr := strings.ToLower(strings.TrimSpace(row[3]))
	if validStr != "true" && validStr != "false" {
		return record, fmt.Errorf("invalid boolean at line %d: must be 'true' or 'false'", lineNum)
	}
	record.Valid = validStr == "true"

	return record, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, []error) {
	validRecords := []DataRecord{}
	errors := []error{}

	for _, record := range records {
		if record.ID <= 0 {
			errors = append(errors, fmt.Errorf("record ID %d: invalid ID (must be positive)", record.ID))
			continue
		}

		if record.Value < 0 {
			errors = append(errors, fmt.Errorf("record ID %d: negative value not allowed", record.ID))
			continue
		}

		if len(record.Name) > 100 {
			errors = append(errors, fmt.Errorf("record ID %d: name exceeds maximum length", record.ID))
			continue
		}

		validRecords = append(validRecords, record)
	}

	return validRecords, errors
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64
	validCount := 0

	for _, record := range records {
		if record.Valid {
			sum += record.Value
			if record.Value > max {
				max = record.Value
			}
			validCount++
		}
	}

	average := 0.0
	if validCount > 0 {
		average = sum / float64(validCount)
	}

	return average, max, validCount
}

func ProcessDataFile(filename string) error {
	records, err := ParseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	validRecords, validationErrors := ValidateRecords(records)

	if len(validationErrors) > 0 {
		fmt.Printf("Validation completed with %d error(s):\n", len(validationErrors))
		for _, err := range validationErrors {
			fmt.Printf("  - %v\n", err)
		}
	}

	average, max, validCount := CalculateStatistics(validRecords)

	fmt.Printf("\nProcessing Results:\n")
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Maximum value: %.2f\n", max)

	if len(validationErrors) > 0 {
		return errors.New("data validation completed with errors")
	}

	return nil
}