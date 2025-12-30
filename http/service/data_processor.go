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
}