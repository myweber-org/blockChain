package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func TransformUserData(rawData []byte) (*UserData, error) {
	var user UserData
	err := json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	if !ValidateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format: %s", user.Email)
	}

	user.Username = SanitizeUsername(user.Username)

	if user.Age < 0 || user.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", user.Age)
	}

	return &user, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  john_doe  ","age":25}`
	user, err := TransformUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func TransformUserData(rawData []byte) (*UserData, error) {
	var user UserData
	err := json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	user.Username = SanitizeUsername(user.Username)

	if !ValidateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format: %s", user.Email)
	}

	if user.Age < 0 || user.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", user.Age)
	}

	return &user, nil
}

func main() {
	rawJSON := []byte(`{"email":"test@example.com","username":"  john_doe  ","age":25}`)
	user, err := TransformUserData(rawJSON)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
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

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
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
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber+1, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber+1, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID format at line %d: %w", lineNumber+1, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber+1)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value format at line %d: %w", lineNumber+1, err)
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
		lineNumber++
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64 = records[0].Value
	count := len(records)

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(count)
	return average, max, count
}

func ValidateRecords(records []DataRecord) []error {
	var errors []error
	seenIDs := make(map[int]bool)

	for i, record := range records {
		if record.ID <= 0 {
			errors = append(errors, fmt.Errorf("record %d: invalid ID %d", i, record.ID))
		}

		if seenIDs[record.ID] {
			errors = append(errors, fmt.Errorf("record %d: duplicate ID %d", i, record.ID))
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			errors = append(errors, fmt.Errorf("record %d: negative value %f", i, record.Value))
		}
	}

	return errors
}
package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ParseUserData attempts to parse JSON data into a map.
func ParseUserData(jsonData []byte) (map[string]interface{}, error) {
	valid, err := ValidateJSON(jsonData)
	if !valid {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := []byte(`{"name": "Alice", "age": 30, "active": true}`)

	parsed, err := ParseUserData(sampleJSON)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Parsed data: %v\n", parsed)
}
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Processed bool    `json:"processed"`
}

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNumber := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNumber, err)
		}

		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(row[1])

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			continue
		}

		record := Record{
			ID:        id,
			Name:      name,
			Value:     value,
			Processed: false,
		}
		records = append(records, record)
		lineNumber++
	}

	return records, nil
}

func convertToJSON(records []Record) (string, error) {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func processRecords(records []Record) []Record {
	for i := range records {
		if records[i].Value > 100 {
			records[i].Value *= 0.9
			records[i].Processed = true
		}
	}
	return records
}

func writeOutput(filename string, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <input.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	records, err := parseCSVFile(inputFile)
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed %d records\n", len(records))

	processed := processRecords(records)

	jsonOutput, err := convertToJSON(processed)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	outputFile := strings.TrimSuffix(inputFile, ".csv") + ".json"
	err = writeOutput(outputFile, jsonOutput)
	if err != nil {
		fmt.Printf("Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Output written to %s\n", outputFile)
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
	return userData, err
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
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func ParseCSVFile(filename string) ([]Record, error) {
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
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 columns, got %d", lineNumber, len(row))
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

func parseRow(row []string, lineNumber int) (Record, error) {
	var record Record

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("line %d: invalid ID format: %w", lineNumber, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("line %d: invalid value format: %w", lineNumber, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("line %d: invalid active flag: %w", lineNumber, err)
	}
	record.Active = active

	return record, nil
}

func ValidateRecords(records []Record) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("record %s has invalid ID: %d", record.Name, record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if strings.TrimSpace(record.Name) == "" {
			return fmt.Errorf("record ID %d has empty name", record.ID)
		}

		if record.Value < 0 {
			return fmt.Errorf("record %s has negative value: %f", record.Name, record.Value)
		}
	}

	return nil
}

func CalculateTotalValue(records []Record) float64 {
	var total float64
	for _, record := range records {
		if record.Active {
			total += record.Value
		}
	}
	return total
}

func FilterActiveRecords(records []Record) []Record {
	var active []Record
	for _, record := range records {
		if record.Active {
			active = append(active, record)
		}
	}
	return active
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func sanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	data.Username = sanitizeUsername(data.Username)

	if !validateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  john_doe  ","age":25}`
	processedData, err := processUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
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

func processCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
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

func validateRecords(records []DataRecord) []string {
	var errors []string
	emailSet := make(map[string]bool)

	for i, record := range records {
		if record.Active != "true" && record.Active != "false" {
			errors = append(errors, fmt.Sprintf("record %d: invalid active status", i+1))
		}

		if emailSet[record.Email] {
			errors = append(errors, fmt.Sprintf("record %d: duplicate email detected", i+1))
		}
		emailSet[record.Email] = true
	}

	return errors
}

func generateSummary(records []DataRecord) {
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
	normalized := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return normalized
}

func (dp *DataProcessor) NormalizeCase(input string, toUpper bool) string {
	cleaned := dp.CleanString(input)
	if toUpper {
		return strings.ToUpper(cleaned)
	}
	return strings.ToLower(cleaned)
}

func (dp *DataProcessor) ExtractAlphanumeric(input string) string {
	alphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleaned := dp.CleanString(input)
	return alphanumericRegex.ReplaceAllString(cleaned, "")
}

func main() {
	processor := NewDataProcessor()
	
	sampleData := "  Hello   World!  This  is   a  test.  "
	
	cleaned := processor.CleanString(sampleData)
	println("Cleaned:", cleaned)
	
	upper := processor.NormalizeCase(sampleData, true)
	println("Uppercase:", upper)
	
	alphanum := processor.ExtractAlphanumeric(sampleData)
	println("Alphanumeric only:", alphanum)
}package main

import (
	"encoding/csv"
	"errors"
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
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 4 {
			continue
		}

		record, parseErr := parseRow(row)
		if parseErr != nil {
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, errors.New("invalid id")
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, errors.New("empty name")
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, errors.New("invalid value")
	}
	record.Value = value

	valid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		record.Valid = false
	} else {
		record.Valid = valid
	}

	return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	validRecords := make([]DataRecord, 0)
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}

	total := 0.0
	count := 0
	for _, record := range records {
		if record.Valid {
			total += record.Value
			count++
		}
	}

	if count == 0 {
		return 0.0
	}
	return total / float64(count)
}package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TrimAndLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}
package main

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

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func transformTags(tags []string) []string {
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

func processUserProfile(rawData []byte) (*UserProfile, error) {
	var profile UserProfile
	if err := json.Unmarshal(rawData, &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if profile.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	profile.Username = normalizeUsername(profile.Username)

	if !validateEmail(profile.Email) {
		return nil, fmt.Errorf("invalid email format: %s", profile.Email)
	}

	if profile.Age < 0 || profile.Age > 150 {
		return nil, fmt.Errorf("age must be between 0 and 150")
	}

	profile.Tags = transformTags(profile.Tags)

	return &profile, nil
}

func main() {
	jsonData := `{
		"id": 1,
		"username": "  JohnDoe  ",
		"email": "john@example.com",
		"age": 30,
		"active": true,
		"tags": ["Go", "backend", "GO", "  ", "devops"]
	}`

	processedProfile, err := processUserProfile([]byte(jsonData))
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}

	fmt.Printf("Processed Profile: %+v\n", processedProfile)
	output, _ := json.MarshalIndent(processedProfile, "", "  ")
	fmt.Println("JSON Output:", string(output))
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

func (dc *DataCleaner) RemoveSpecialChars(input string, keep string) string {
	pattern := "[^a-zA-Z0-9" + regexp.QuoteMeta(keep) + "]+"
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, "")
}

func (dc *DataCleaner) TruncateString(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	return input[:maxLength-3] + "..."
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ParseUserData attempts to parse JSON data into a generic map.
func ParseUserData(rawData []byte) (map[string]interface{}, error) {
	valid, err := ValidateJSON(rawData)
	if !valid {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rawData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := []byte(`{"name": "alice", "age": 30, "active": true}`)

	parsed, err := ParseUserData(sampleJSON)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Parsed data: %v\n", parsed)
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(username, "")
	return username
}

func TransformUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	rawJSON := []byte(`{"email":"test@example.com","username":"user_123!","age":25}`)
	processedData, err := TransformUserData(rawJSON)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed Data: %+v\n", processedData)
}