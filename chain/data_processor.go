
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	cleaned := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return cleaned
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (dp *DataProcessor) ExtractDomain(email string) (string, bool) {
	if !dp.ValidateEmail(email) {
		return "", false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", false
	}
	return parts[1], true
}
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
	Valid bool
}

func ProcessRecord(record DataRecord) (string, error) {
	if !record.Valid {
		return "", errors.New("invalid record")
	}

	if record.ID <= 0 {
		return "", errors.New("invalid ID")
	}

	if strings.TrimSpace(record.Name) == "" {
		return "", errors.New("name cannot be empty")
	}

	if record.Value < 0 {
		return "", errors.New("value cannot be negative")
	}

	processedName := strings.ToUpper(record.Name)
	result := fmt.Sprintf("ID:%d|NAME:%s|VALUE:%.2f", record.ID, processedName, record.Value)

	return result, nil
}

func ValidateAndProcess(records []DataRecord) ([]string, []error) {
	var results []string
	var errs []error

	for _, record := range records {
		result, err := ProcessRecord(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("record %d: %w", record.ID, err))
			continue
		}
		results = append(results, result)
	}

	return results, errs
}

func main() {
	records := []DataRecord{
		{ID: 1, Name: "record_one", Value: 100.5, Valid: true},
		{ID: 2, Name: "", Value: 50.0, Valid: true},
		{ID: 3, Name: "record_three", Value: -10.0, Valid: true},
		{ID: 0, Name: "record_four", Value: 75.3, Valid: true},
		{ID: 5, Name: "record_five", Value: 200.8, Valid: false},
	}

	results, errs := ValidateAndProcess(records)

	fmt.Println("Processing Results:")
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\nErrors:")
	for _, err := range errs {
		fmt.Println(err)
	}
}