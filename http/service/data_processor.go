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