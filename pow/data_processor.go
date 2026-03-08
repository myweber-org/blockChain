
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

func (dp *DataProcessor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return normalized
}

func (dp *DataProcessor) NormalizeCase(input string, toUpper bool) string {
	cleaned := dp.CleanInput(input)
	if toUpper {
		return strings.ToUpper(cleaned)
	}
	return strings.ToLower(cleaned)
}

func (dp *DataProcessor) ExtractAlphanumeric(input string) string {
	alnumRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return alnumRegex.ReplaceAllString(input, "")
}

func (dp *DataProcessor) ValidateEmail(input string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(input)
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ReadCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
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

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		record := DataRecord{
			ID:    id,
			Name:  row[1],
			Value: value,
		}
		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d must be positive", record.ID)
		}
		if record.Name == "" {
			return fmt.Errorf("record %d has empty name", record.ID)
		}
		if record.Value < 0 {
			return fmt.Errorf("record %d has negative value: %f", record.ID, record.Value)
		}
		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true
	}
	return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64
	count := len(records)

	for i, record := range records {
		sum += record.Value
		if i == 0 || record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(count)
	return average, max, count
}
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
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, err
	}

	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
	transformed.Tags = normalizeTags(transformed.Tags)
	return transformed, nil
}

func normalizeTags(tags []string) []string {
	uniqueTags := make(map[string]bool)
	var result []string

	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !uniqueTags[normalized] {
			uniqueTags[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	var errors []string

	for i, record := range records {
		transformed, err := TransformRecord(record)
		if err != nil {
			errors = append(errors, fmt.Sprintf("record %d: %v", i, err))
			continue
		}
		processed = append(processed, transformed)
	}

	if len(errors) > 0 {
		return processed, fmt.Errorf("processing errors: %s", strings.Join(errors, "; "))
	}
	return processed, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records provided")
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
	return average, max - min, nil
}package main

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

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	return input
}

func (dp *DataProcessor) NormalizeUsername(username string) string {
	username = strings.ToLower(username)
	username = strings.TrimSpace(username)
	return strings.ReplaceAll(username, " ", "_")
}

func (dp *DataProcessor) ProcessUserData(email, username, comment string) (bool, string, string, string) {
	isValidEmail := dp.ValidateEmail(email)
	sanitizedComment := dp.SanitizeInput(comment)
	normalizedUsername := dp.NormalizeUsername(username)
	
	return isValidEmail, normalizedUsername, sanitizedComment, email
}
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
	headerSkipped := false

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		if !headerSkipped {
			headerSkipped = true
			continue
		}

		if len(row) < 4 {
			continue
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if record.ID == "" || record.Email == "" {
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterActiveUsers(records []DataRecord) []DataRecord {
	var active []DataRecord
	for _, record := range records {
		if strings.ToLower(record.Active) == "true" && ValidateEmail(record.Email) {
			active = append(active, record)
		}
	}
	return active
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run data_processor.go <csv_file>")
		return
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	activeUsers := FilterActiveUsers(records)
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Active users: %d\n", len(activeUsers))

	for _, user := range activeUsers {
		fmt.Printf("ID: %s, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}
}