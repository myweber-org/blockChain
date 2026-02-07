
package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Value   int    `json:"value"`
    Active  bool   `json:"active"`
}

func parseCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []Record{}
    lineNumber := 0

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if lineNumber == 0 {
            lineNumber++
            continue
        }

        id, _ := strconv.Atoi(row[0])
        value, _ := strconv.Atoi(row[2])
        active, _ := strconv.ParseBool(row[3])

        record := Record{
            ID:     id,
            Name:   row[1],
            Value:  value,
            Active: active,
        }
        records = append(records, record)
        lineNumber++
    }
    return records, nil
}

func convertToJSON(records []Record) (string, error) {
    jsonData, err := json.MarshalIndent(records, "", "  ")
    if err != nil {
        return "", err
    }
    return string(jsonData), nil
}

func filterActiveRecords(records []Record) []Record {
    var active []Record
    for _, record := range records {
        if record.Active {
            active = append(active, record)
        }
    }
    return active
}

func calculateTotalValue(records []Record) int {
    total := 0
    for _, record := range records {
        total += record.Value
    }
    return total
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        return
    }

    records, err := parseCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error parsing CSV: %v\n", err)
        return
    }

    fmt.Printf("Total records: %d\n", len(records))
    
    activeRecords := filterActiveRecords(records)
    fmt.Printf("Active records: %d\n", len(activeRecords))
    
    totalValue := calculateTotalValue(records)
    fmt.Printf("Total value: %d\n", totalValue)

    jsonOutput, err := convertToJSON(records)
    if err != nil {
        fmt.Printf("Error converting to JSON: %v\n", err)
        return
    }

    outputFile := "output.json"
    err = os.WriteFile(outputFile, []byte(jsonOutput), 0644)
    if err != nil {
        fmt.Printf("Error writing JSON file: %v\n", err)
        return
    }
    
    fmt.Printf("JSON output written to %s\n", outputFile)
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
	alnumRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleaned := dp.CleanString(input)
	return alnumRegex.ReplaceAllString(cleaned, "")
}

func (dp *DataProcessor) ValidateEmail(input string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(strings.TrimSpace(input))
}
package main

import (
	"regexp"
	"strings"
)

// CleanString removes extra whitespace and normalizes input
func CleanString(input string) string {
	// Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	re := regexp.MustCompile(`\s+`)
	normalized := re.ReplaceAllString(trimmed, " ")
	
	return normalized
}

// NormalizeEmail converts email to lowercase and trims spaces
func NormalizeEmail(email string) string {
	cleaned := CleanString(email)
	return strings.ToLower(cleaned)
}

// ValidateUsername checks if username meets requirements
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	
	// Only allow alphanumeric and underscore
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validPattern.MatchString(username)
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// DataPayload represents a generic JSON payload structure.
type DataPayload struct {
	ID      string          `json:"id"`
	Version int             `json:"version"`
	Content json.RawMessage `json:"content"`
	Active  bool            `json:"active"`
}

// ValidateJSON checks if the provided byte slice is valid JSON.
func ValidateJSON(data []byte) bool {
	return json.Valid(data)
}

// ParsePayload attempts to parse JSON data into a DataPayload struct.
func ParsePayload(data []byte) (*DataPayload, error) {
	if !ValidateJSON(data) {
		return nil, fmt.Errorf("invalid JSON format")
	}

	var payload DataPayload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if payload.ID == "" {
		return nil, fmt.Errorf("payload ID cannot be empty")
	}

	return &payload, nil
}

func main() {
	jsonStr := `{"id": "user123", "version": 2, "content": {"name": "John"}, "active": true}`

	payload, err := ParsePayload([]byte(jsonStr))
	if err != nil {
		log.Fatalf("Error parsing payload: %v", err)
	}

	fmt.Printf("Parsed payload: ID=%s, Version=%d, Active=%v\n", payload.ID, payload.Version, payload.Active)
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
	alnumRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return alnumRegex.ReplaceAllString(input, "")
}

func (dp *DataProcessor) ValidateEmail(input string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(input)
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

func TransformUserData(users []UserProfile) ([]map[string]interface{}, error) {
	var transformed []map[string]interface{}
	
	for _, user := range users {
		if !ValidateEmail(user.Email) {
			return nil, fmt.Errorf("invalid email for user %d: %s", user.ID, user.Email)
		}
		
		data := map[string]interface{}{
			"user_id":   user.ID,
			"username":  NormalizeUsername(user.Username),
			"email":     user.Email,
			"age_group": getAgeGroup(user.Age),
			"status":    "active",
			"tag_count": len(user.Tags),
		}
		
		transformed = append(transformed, data)
	}
	
	return transformed, nil
}

func getAgeGroup(age int) string {
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

func main() {
	users := []UserProfile{
		{ID: 1, Username: "  JohnDoe  ", Email: "john@example.com", Age: 25, Active: true, Tags: []string{"golang", "backend"}},
		{ID: 2, Username: "JaneSmith", Email: "jane@example.org", Age: 42, Active: false, Tags: []string{"frontend"}},
		{ID: 3, Username: "Bob", Email: "bob@test.co", Age: 17, Active: true, Tags: []string{}},
	}
	
	activeUsers := FilterInactiveUsers(users)
	fmt.Printf("Active users: %d\n", len(activeUsers))
	
	transformed, err := TransformUserData(activeUsers)
	if err != nil {
		fmt.Printf("Error transforming data: %v\n", err)
		return
	}
	
	jsonData, _ := json.MarshalIndent(transformed, "", "  ")
	fmt.Println("Transformed data:")
	fmt.Println(string(jsonData))
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

func processCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

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

		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func calculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		return
	}

	filename := os.Args[1]
	records, err := processCSV(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	fmt.Printf("Processed %d records\n", len(records))

	if len(records) > 0 {
		avg, max := calculateStats(records)
		fmt.Printf("Average value: %.2f\n", avg)
		fmt.Printf("Maximum value: %.2f\n", max)
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
			return nil, fmt.Errorf("invalid column count at line %d", lineNum)
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

	return records, nil
}

func ValidateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, rec := range records {
		if rec.ID <= 0 {
			return fmt.Errorf("invalid ID %d", rec.ID)
		}
		if seenIDs[rec.ID] {
			return fmt.Errorf("duplicate ID %d", rec.ID)
		}
		seenIDs[rec.ID] = true
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	if err := ValidateRecords(records); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	for _, rec := range records {
		fmt.Printf("ID: %d, Name: %s, Value: %.2f\n", rec.ID, rec.Name, rec.Value)
	}
}
package main

import (
	"errors"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Timestamp time.Time
	Value     float64
	Tags      []string
	Valid     bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value cannot be negative")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("record timestamp cannot be in the future")
	}
	return nil
}

func TransformTags(tags []string) []string {
	transformed := make([]string, 0, len(tags))
	for _, tag := range tags {
		cleanTag := strings.TrimSpace(tag)
		cleanTag = strings.ToLower(cleanTag)
		if cleanTag != "" {
			transformed = append(transformed, cleanTag)
		}
	}
	return transformed
}

func CalculateAverage(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, errors.New("cannot calculate average of empty slice")
	}
	
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values)), nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	validRecords := make([]DataRecord, 0)
	for _, record := range records {
		if record.Valid && ValidateRecord(record) == nil {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}