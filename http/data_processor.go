
package main

import "fmt"

func MovingAverage(data []float64, windowSize int) []float64 {
    if len(data) == 0 || windowSize <= 0 || windowSize > len(data) {
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
    averages := MovingAverage(sampleData, window)
    fmt.Printf("Moving averages (window=%d): %v\n", window, averages)
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
	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Username: transformedUsername,
		Email:    strings.TrimSpace(email),
		Age:      age,
	}
	
	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}
	
	return userData, nil
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

type Record struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ProcessCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
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

		record, err := parseRecord(row, lineNumber)
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

func parseRecord(row []string, lineNumber int) (Record, error) {
	var record Record

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return Record{}, fmt.Errorf("empty name at line %d", lineNumber)
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	validStr := strings.ToLower(strings.TrimSpace(row[3]))
	if validStr != "true" && validStr != "false" {
		return Record{}, fmt.Errorf("invalid boolean at line %d: must be 'true' or 'false'", lineNumber)
	}
	record.Valid = validStr == "true"

	return record, nil
}

func CalculateStats(records []Record) (float64, float64, int) {
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

	var average float64
	if validCount > 0 {
		average = sum / float64(validCount)
	}

	return average, maxValue, validCount
}

func FilterValidRecords(records []Record) []Record {
	var validRecords []Record
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ProcessCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d records\n", len(records))

	validRecords := FilterValidRecords(records)
	fmt.Printf("Valid records: %d\n", len(validRecords))

	average, maxValue, validCount := CalculateStats(records)
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Maximum value: %.2f\n", maxValue)
	fmt.Printf("Valid record count: %d\n", validCount)
}package main

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

func TransformUserData(user UserProfile) map[string]interface{} {
	return map[string]interface{}{
		"user_id":   user.ID,
		"name":      NormalizeUsername(user.Username),
		"contact":   user.Email,
		"age_group": categorizeAge(user.Age),
		"status":    user.Active,
		"metadata":  len(user.Tags),
	}
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
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	var validUsers []UserProfile
	for _, user := range users {
		if ValidateEmail(user.Email) && user.Age > 0 {
			user.Username = NormalizeUsername(user.Username)
			validUsers = append(validUsers, user)
		}
	}

	activeUsers := FilterInactiveUsers(validUsers)
	var result []map[string]interface{}
	for _, user := range activeUsers {
		result = append(result, TransformUserData(user))
	}

	return result, nil
}

func main() {
	jsonInput := `[
		{"id":1,"username":"JohnDoe ","email":"john@example.com","age":25,"active":true,"tags":["gamer","developer"]},
		{"id":2,"username":"JaneSmith","email":"invalid-email","age":30,"active":false,"tags":["designer"]},
		{"id":3,"username":" Bob ","email":"bob@test.org","age":17,"active":true,"tags":[]}
	]`

	processed, err := ProcessUserJSON([]byte(jsonInput))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	output, _ := json.MarshalIndent(processed, "", "  ")
	fmt.Println("Processed users:", string(output))
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

func (dp *DataProcessor) SanitizeString(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToValidUTF8(trimmed, "")
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, string, bool) {
	sanitizedName := dp.SanitizeString(name)
	sanitizedEmail := dp.SanitizeString(email)
	isValidEmail := dp.ValidateEmail(sanitizedEmail)
	return sanitizedName, sanitizedEmail, isValidEmail
}package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Email:    strings.TrimSpace(email),
		Username: transformedUsername,
		Age:      age,
	}
	err := ValidateUserData(userData)
	if err != nil {
		return UserData{}, err
	}
	return userData, nil
}
package main

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

func (dc *DataCleaner) ProcessUserInput(rawInput string) (string, bool) {
	cleaned := dc.NormalizeWhitespace(rawInput)
	if cleaned == "" {
		return "", false
	}
	return cleaned, true
}