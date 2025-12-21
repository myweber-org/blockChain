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

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func TransformProfile(profile UserProfile) (UserProfile, error) {
	if profile.Age < 0 || profile.Age > 120 {
		return profile, fmt.Errorf("invalid age: %d", profile.Age)
	}

	if !ValidateEmail(profile.Email) {
		return profile, fmt.Errorf("invalid email format: %s", profile.Email)
	}

	profile.Username = SanitizeUsername(profile.Username)

	if len(profile.Tags) > 10 {
		profile.Tags = profile.Tags[:10]
	}

	return profile, nil
}

func ProcessUserData(jsonData []byte) ([]byte, error) {
	var profiles []UserProfile
	if err := json.Unmarshal(jsonData, &profiles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var validProfiles []UserProfile
	for _, profile := range profiles {
		transformed, err := TransformProfile(profile)
		if err != nil {
			fmt.Printf("Skipping profile ID %d: %v\n", profile.ID, err)
			continue
		}
		validProfiles = append(validProfiles, transformed)
	}

	result, err := json.MarshalIndent(validProfiles, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return result, nil
}

func main() {
	sampleData := `[
		{"id":1,"username":"  JohnDoe  ","email":"john@example.com","age":25,"active":true,"tags":["golang","backend"]},
		{"id":2,"username":"JaneSmith","email":"invalid-email","age":150,"active":false,"tags":["frontend","design","test","extra"]},
		{"id":3,"username":"Bob","email":"bob@test.org","age":30,"active":true,"tags":[]}
	]`

	processed, err := ProcessUserData([]byte(sampleData))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Println("Processed user profiles:")
	fmt.Println(string(processed))
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
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVFile(filepath string) ([]DataRecord, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
		lineNumber++
	}

	if len(records) == 0 {
		return nil, errors.New("no valid data records found")
	}

	return records, nil
}

func parseRow(row []string, lineNum int) (DataRecord, error) {
	if len(row) < 4 {
		return DataRecord{}, fmt.Errorf("invalid column count at line %d", lineNum)
	}

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}

	name := strings.TrimSpace(row[1])
	if name == "" {
		return DataRecord{}, fmt.Errorf("empty name at line %d", lineNum)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}

	valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

	return DataRecord{
		ID:    id,
		Name:  name,
		Value: value,
		Valid: valid,
	}, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, []DataRecord) {
	valid := []DataRecord{}
	invalid := []DataRecord{}

	for _, record := range records {
		if record.Valid && record.Value >= 0 {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}

	return valid, invalid
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}

	total := 0.0
	for _, record := range records {
		total += record.Value
	}

	return total / float64(len(records))
}

func GenerateReport(valid, invalid []DataRecord) string {
	report := strings.Builder{}
	report.WriteString(fmt.Sprintf("Data Processing Report\n"))
	report.WriteString(fmt.Sprintf("======================\n"))
	report.WriteString(fmt.Sprintf("Valid records: %d\n", len(valid)))
	report.WriteString(fmt.Sprintf("Invalid records: %d\n", len(invalid)))
	report.WriteString(fmt.Sprintf("Average value: %.2f\n", CalculateAverage(valid)))

	if len(invalid) > 0 {
		report.WriteString("\nInvalid Records:\n")
		for _, record := range invalid {
			report.WriteString(fmt.Sprintf("  ID: %d, Name: %s, Value: %.2f\n", 
				record.ID, record.Name, record.Value))
		}
	}

	return report.String()
}
package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(trimmed, " ")
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func ContainsOnlyAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}
package main

import (
	"regexp"
	"strings"
)

// DataProcessor handles cleaning and normalization of string data
type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

// NewDataProcessor creates a new DataProcessor instance
func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

// CleanString removes extra whitespace and trims the input
func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	cleaned := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return cleaned
}

// NormalizeString converts string to lowercase and cleans it
func (dp *DataProcessor) NormalizeString(input string) string {
	cleaned := dp.CleanString(input)
	return strings.ToLower(cleaned)
}

// ProcessBatch applies normalization to a slice of strings
func (dp *DataProcessor) ProcessBatch(inputs []string) []string {
	results := make([]string, len(inputs))
	for i, input := range inputs {
		results[i] = dp.NormalizeString(input)
	}
	return results
}
package main

import "fmt"

func MovingAverage(data []float64, window int) []float64 {
    if window <= 0 || window > len(data) {
        return nil
    }

    result := make([]float64, len(data)-window+1)
    var sum float64

    for i := 0; i < window; i++ {
        sum += data[i]
    }
    result[0] = sum / float64(window)

    for i := window; i < len(data); i++ {
        sum = sum - data[i-window] + data[i]
        result[i-window+1] = sum / float64(window)
    }

    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    averaged := MovingAverage(sampleData, 3)
    fmt.Println("Moving average with window 3:", averaged)
}