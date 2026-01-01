
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

func (dp *DataProcessor) NormalizeWhitespace(input string) string {
	return dp.whitespaceRegex.ReplaceAllString(input, " ")
}
package main

import (
	"regexp"
	"strings"
)

func ProcessInput(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", nil
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	cleaned := re.ReplaceAllString(trimmed, "")

	reMultipleSpaces := regexp.MustCompile(`\s+`)
	final := reMultipleSpaces.ReplaceAllString(cleaned, " ")

	return strings.ToLower(final), nil
}
package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int
	Email     string
	Username  string
	BirthDate string
	Active    bool
	CreatedAt time.Time
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func CalculateAge(birthDate string) (int, error) {
	parsedDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	age := time.Since(parsedDate).Hours() / 24 / 365.25
	return int(age), nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateEmail(profile.Email); err != nil {
		return profile, err
	}

	profile.Username = SanitizeUsername(profile.Username)

	age, err := CalculateAge(profile.BirthDate)
	if err != nil {
		return profile, err
	}

	if age < 13 {
		return profile, errors.New("user must be at least 13 years old")
	}

	profile.CreatedAt = time.Now()
	return profile, nil
}

func main() {
	profile := UserProfile{
		ID:        1,
		Email:     "user@example.com",
		Username:  "  john_doe  ",
		BirthDate: "1990-05-15",
		Active:    true,
	}

	processedProfile, err := ProcessUserProfile(profile)
	if err != nil {
		println("Error processing profile:", err.Error())
		return
	}

	println("Profile processed successfully")
	println("Username:", processedProfile.Username)
	println("Created at:", processedProfile.CreatedAt.Format("2006-01-02"))
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Tags      []string  `json:"tags"`
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func NormalizeTags(tags []string) []string {
	uniqueTags := make(map[string]bool)
	var normalized []string

	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		lower := strings.ToLower(trimmed)
		if trimmed != "" && !uniqueTags[lower] {
			uniqueTags[lower] = true
			normalized = append(normalized, trimmed)
		}
	}
	return normalized
}

func TransformRecord(record DataRecord) (DataRecord, error) {
	if !ValidateEmail(record.Email) {
		return DataRecord{}, fmt.Errorf("invalid email format: %s", record.Email)
	}

	if record.Value < 0 {
		record.Value = 0
	}

	record.Tags = NormalizeTags(record.Tags)

	if record.ID == "" {
		record.ID = fmt.Sprintf("rec_%d", time.Now().UnixNano())
	}

	return record, nil
}

func ProcessRecords(records []DataRecord) ([]DataRecord, []error) {
	var processed []DataRecord
	var errors []error

	for _, record := range records {
		transformed, err := TransformRecord(record)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		processed = append(processed, transformed)
	}

	return processed, errors
}

func RecordsToJSON(records []DataRecord) (string, error) {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
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
    ID      int
    Name    string
    Value   float64
    Active  bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)
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

        if len(row) != 4 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
        }

        record, err := parseRow(row, lineNumber)
        if err != nil {
            return nil, err
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func parseRow(row []string, lineNumber int) (DataRecord, error) {
    var record DataRecord
    var err error

    record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
    }

    record.Name = strings.TrimSpace(row[1])
    if record.Name == "" {
        return DataRecord{}, fmt.Errorf("empty name at line %d", lineNumber)
    }

    record.Value, err = strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
    }

    record.Active, err = strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
    }

    return record, nil
}

func ValidateRecords(records []DataRecord) []error {
    var errors []error
    seenIDs := make(map[int]bool)

    for i, record := range records {
        if record.ID <= 0 {
            errors = append(errors, fmt.Errorf("record %d: ID must be positive", i+1))
        }

        if seenIDs[record.ID] {
            errors = append(errors, fmt.Errorf("record %d: duplicate ID %d", i+1, record.ID))
        }
        seenIDs[record.ID] = true

        if record.Value < 0 {
            errors = append(errors, fmt.Errorf("record %d: value cannot be negative", i+1))
        }
    }

    return errors
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var activeCount int
    var maxValue float64

    for _, record := range records {
        sum += record.Value
        if record.Value > maxValue {
            maxValue = record.Value
        }
        if record.Active {
            activeCount++
        }
    }

    average := sum / float64(len(records))
    return average, maxValue, activeCount
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		lineNum++
		if lineNum == 1 {
			continue
		}

		if len(line) != 3 {
			return nil, fmt.Errorf("invalid column count on line %d", lineNum)
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID on line %d: %w", lineNum, err)
		}

		name := line[1]
		if name == "" {
			return nil, fmt.Errorf("empty name on line %d", lineNum)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value on line %d: %w", lineNum, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func CalculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64

	for i, record := range records {
		sum += record.Value
		if i == 0 || record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ProcessCSV(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d records\n", len(records))
	
	avg, max := CalculateStats(records)
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
}package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateUserData(data UserData) error {
	if data.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Username == "" {
		return errors.New("username is required")
	}
	if len(data.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func NormalizeUserData(data UserData) UserData {
	return UserData{
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Username: strings.TrimSpace(data.Username),
		Age:      data.Age,
	}
}