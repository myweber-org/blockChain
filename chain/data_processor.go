package main

import (
    "encoding/csv"
    "errors"
    "io"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
    Valid bool
}

func ParseCSVData(reader io.Reader) ([]DataRecord, error) {
    csvReader := csv.NewReader(reader)
    records, err := csvReader.ReadAll()
    if err != nil {
        return nil, err
    }

    if len(records) == 0 {
        return []DataRecord{}, nil
    }

    var data []DataRecord
    for i, row := range records {
        if len(row) < 4 {
            continue
        }

        id, err := strconv.Atoi(strings.TrimSpace(row[0]))
        if err != nil {
            continue
        }

        name := strings.TrimSpace(row[1])
        if name == "" {
            continue
        }

        value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
        if err != nil {
            continue
        }

        valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

        data = append(data, DataRecord{
            ID:    id,
            Name:  name,
            Value: value,
            Valid: valid,
        })
    }

    return data, nil
}

func ValidateData(records []DataRecord) error {
    if len(records) == 0 {
        return errors.New("no data records provided")
    }

    idSet := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return errors.New("invalid ID found: " + strconv.Itoa(record.ID))
        }

        if idSet[record.ID] {
            return errors.New("duplicate ID found: " + strconv.Itoa(record.ID))
        }
        idSet[record.ID] = true

        if record.Value < 0 {
            return errors.New("negative value not allowed for ID: " + strconv.Itoa(record.ID))
        }
    }

    return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var validCount int
    var maxValue float64

    for _, record := range records {
        if record.Valid {
            sum += record.Value
            validCount++
            if record.Value > maxValue {
                maxValue = record.Value
            }
        }
    }

    if validCount == 0 {
        return 0, 0, 0
    }

    average := sum / float64(validCount)
    return average, maxValue, validCount
}package main

import (
	"errors"
	"strings"
	"unicode"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username can only contain letters, digits, underscores, and hyphens")
		}
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TransformUserData(rawUsername, rawEmail string) (string, string, error) {
	if err := ValidateUsername(rawUsername); err != nil {
		return "", "", err
	}
	normalizedEmail := NormalizeEmail(rawEmail)
	return rawUsername, normalizedEmail, nil
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

		if lineNumber == 1 {
			continue
		}

		if len(row) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if record.ID == "" || record.Name == "" {
			return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
		}

		if !strings.Contains(record.Email, "@") {
			return nil, fmt.Errorf("invalid email format at line %d", lineNumber)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) []string {
	var errors []string
	emailSet := make(map[string]bool)

	for i, record := range records {
		if record.Active != "true" && record.Active != "false" {
			errors = append(errors, fmt.Sprintf("record %d: active field must be 'true' or 'false'", i+1))
		}

		if emailSet[record.Email] {
			errors = append(errors, fmt.Sprintf("record %d: duplicate email detected", i+1))
		}
		emailSet[record.Email] = true
	}

	return errors
}

func GenerateReport(records []DataRecord) {
	activeCount := 0
	for _, record := range records {
		if record.Active == "true" {
			activeCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Active records: %d\n", activeCount)
	fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
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

func processCSVFile(filename string) ([]DataRecord, error) {
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

		if isValidRecord(record) {
			records = append(records, record)
		}
	}

	return records, nil
}

func isValidRecord(record DataRecord) bool {
	if record.ID == "" || record.Name == "" {
		return false
	}
	if !strings.Contains(record.Email, "@") {
		return false
	}
	return record.Active == "true" || record.Active == "false"
}

func generateReport(records []DataRecord) {
	activeCount := 0
	for _, record := range records {
		if record.Active == "true" {
			activeCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Active records: %d\n", activeCount)
	fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	generateReport(records)
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
	BirthDate time.Time
	Active    bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if len(strings.TrimSpace(profile.Username)) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if time.Since(profile.BirthDate).Hours()/24/365 < 13 {
		return errors.New("user must be at least 13 years old")
	}

	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func CalculateAge(birthDate time.Time) int {
	now := time.Now()
	years := now.Year() - birthDate.Year()

	if now.YearDay() < birthDate.YearDay() {
		years--
	}

	return years
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateUserProfile(profile); err != nil {
		return UserProfile{}, err
	}

	profile.Username = TransformUsername(profile.Username)
	return profile, nil
}package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	emailRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return &DataProcessor{emailRegex: regex}
}

func (dp *DataProcessor) SanitizeString(input string) string {
	return strings.TrimSpace(input)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, bool) {
	sanitizedName := dp.SanitizeString(name)
	isValidEmail := dp.ValidateEmail(email)

	if sanitizedName == "" || !isValidEmail {
		return "", false
	}

	return sanitizedName, true
}package main

import (
	"errors"
	"strings"
	"unicode"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username can only contain letters, digits, underscores, and hyphens")
		}
	}

	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func SanitizeInput(input string) string {
	replacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
		"&", "&amp;",
	)
	return replacer.Replace(input)
}

func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("password must contain uppercase, lowercase, digit, and special character")
	}

	return nil
}
package data_processor

import (
	"regexp"
	"strings"
)

type DataCleaner struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dc *DataCleaner) NormalizeWhitespace(input string) string {
	trimmed := strings.TrimSpace(input)
	return dc.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	return dc.emailRegex.MatchString(email)
}

func (dc *DataCleaner) ExtractDomain(email string) string {
	if !dc.ValidateEmail(email) {
		return ""
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func (dc *DataCleaner) SanitizeInput(input string, maxLength int) string {
	normalized := dc.NormalizeWhitespace(input)
	if len(normalized) > maxLength {
		return normalized[:maxLength]
	}
	return normalized
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
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(line) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNum, len(line))
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
		}

		name := line[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNum)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
		lineNum++
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func CalculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, r := range records {
		sum += r.Value
		if r.Value > max {
			max = r.Value
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

	avg, max := CalculateStats(records)
	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateAndParseJSON(rawData []byte) (*UserData, error) {
	var user UserData
	if err := json.Unmarshal(rawData, &user); err != nil {
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
	jsonData := []byte(`{"id": 101, "name": "Alice", "email": "alice@example.com"}`)
	user, err := ValidateAndParseJSON(jsonData)
	if err != nil {
		log.Fatal("Error processing data:", err)
	}
	fmt.Printf("Processed user: %+v\n", user)
}package main

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
	if age < 0 {
		return 0, errors.New("birth date cannot be in the future")
	}
	return int(age), nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateEmail(profile.Email); err != nil {
		return profile, err
	}

	profile.Username = SanitizeUsername(profile.Username)

	_, err := CalculateAge(profile.BirthDate)
	if err != nil {
		return profile, err
	}

	return profile, nil
}
package data_processor

import (
	"regexp"
	"strings"
	"unicode"
)

type DataCleaner struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dc *DataCleaner) NormalizeWhitespace(input string) string {
	return dc.whitespaceRegex.ReplaceAllString(strings.TrimSpace(input), " ")
}

func (dc *DataCleaner) RemoveSpecialCharacters(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, input)
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	return dc.emailRegex.MatchString(email)
}

func (dc *DataCleaner) ProcessString(input string) string {
	normalized := dc.NormalizeWhitespace(input)
	cleaned := dc.RemoveSpecialCharacters(normalized)
	return strings.ToLower(cleaned)
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) bool {
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}

// ParseUserData attempts to parse JSON data into a map.
func ParseUserData(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := `{"name": "Alice", "age": 30, "active": true}`

	if ValidateJSON([]byte(sampleJSON)) {
		fmt.Println("JSON is valid.")
		data, err := ParseUserData(sampleJSON)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Parsed data: %v\n", data)
	} else {
		fmt.Println("Invalid JSON.")
	}
}
package main

import "fmt"

func movingAverage(data []float64, windowSize int) []float64 {
    if windowSize <= 0 || windowSize > len(data) {
        return nil
    }

    result := make([]float64, len(data)-windowSize+1)
    var sum float64

    for i := 0; i < windowSize; i++ {
        sum += data[i]
    }
    result[0] = sum / float64(windowSize)

    for i := windowSize; i < len(data); i++ {
        sum = sum - data[i-windowSize] + data[i]
        result[i-windowSize+1] = sum / float64(windowSize)
    }

    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    window := 3
    averages := movingAverage(sampleData, window)

    fmt.Printf("Data: %v\n", sampleData)
    fmt.Printf("Moving average (window=%d): %v\n", window, averages)
}