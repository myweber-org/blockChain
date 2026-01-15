
package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Email     string
	Timestamp time.Time
	Status    string
}

func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return fmt.Errorf("regex validation failed: %w", err)
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeString(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func ProcessRecord(record DataRecord) (DataRecord, error) {
	if record.ID == "" {
		return record, errors.New("record ID cannot be empty")
	}

	if err := ValidateEmail(record.Email); err != nil {
		return record, fmt.Errorf("email validation failed: %w", err)
	}

	record.Email = NormalizeString(record.Email)
	record.Status = NormalizeString(record.Status)

	if record.Timestamp.IsZero() {
		record.Timestamp = time.Now().UTC()
	}

	return record, nil
}

func TransformRecords(records []DataRecord) ([]DataRecord, []error) {
	var processed []DataRecord
	var errs []error

	for i, record := range records {
		processedRecord, err := ProcessRecord(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("record %d: %w", i, err))
			continue
		}
		processed = append(processed, processedRecord)
	}

	return processed, errs
}

func main() {
	records := []DataRecord{
		{ID: "001", Email: "USER@EXAMPLE.COM", Timestamp: time.Now(), Status: "ACTIVE"},
		{ID: "002", Email: "invalid-email", Status: "PENDING"},
		{ID: "", Email: "test@domain.com", Status: "inactive"},
	}

	processed, errs := TransformRecords(records)

	fmt.Printf("Processed %d records successfully\n", len(processed))
	if len(errs) > 0 {
		fmt.Printf("Encountered %d errors:\n", len(errs))
		for _, err := range errs {
			fmt.Printf("  - %v\n", err)
		}
	}

	for _, record := range processed {
		fmt.Printf("ID: %s, Email: %s, Status: %s\n",
			record.ID, record.Email, record.Status)
	}
}