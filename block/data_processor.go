
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
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Value float64 `json:"value"`
}

func processCSVFile(inputPath string) ([]Record, error) {
    file, err := os.Open(inputPath)
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
            return nil, fmt.Errorf("csv read error on line %d: %w", lineNumber, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count on line %d: expected 3, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID on line %d: %w", lineNumber, err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value on line %d: %w", lineNumber, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    return records, nil
}

func generateJSONOutput(records []Record, outputPath string) error {
    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    encoder := json.NewEncoder(outputFile)
    encoder.SetIndent("", "  ")

    if err := encoder.Encode(records); err != nil {
        return fmt.Errorf("failed to encode JSON: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_processor <input.csv> <output.json>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    records, err := processCSVFile(inputFile)
    if err != nil {
        fmt.Printf("Error processing CSV: %v\n", err)
        os.Exit(1)
    }

    if err := generateJSONOutput(records, outputFile); err != nil {
        fmt.Printf("Error generating JSON: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d records to %s\n", len(records), outputFile)
}package data_processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type ValidationRule func(interface{}) error

type JSONProcessor struct {
	rules map[string]ValidationRule
}

func NewJSONProcessor() *JSONProcessor {
	return &JSONProcessor{
		rules: make(map[string]ValidationRule),
	}
}

func (jp *JSONProcessor) AddRule(field string, rule ValidationRule) {
	jp.rules[field] = rule
}

func (jp *JSONProcessor) Process(rawData []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}
	
	if len(jp.rules) == 0 {
		return data, nil
	}
	
	var validationErrors []string
	
	for field, rule := range jp.rules {
		if value, exists := data[field]; exists {
			if err := rule(value); err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("%s: %v", field, err))
			}
		}
	}
	
	if len(validationErrors) > 0 {
		return nil, errors.New(strings.Join(validationErrors, "; "))
	}
	
	return data, nil
}

func RequiredStringRule(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("must be a string")
	}
	
	if strings.TrimSpace(str) == "" {
		return errors.New("cannot be empty")
	}
	
	return nil
}

func PositiveNumberRule(value interface{}) error {
	num, ok := value.(float64)
	if !ok {
		return errors.New("must be a number")
	}
	
	if num <= 0 {
		return errors.New("must be positive")
	}
	
	return nil
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

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
	transformed.Tags = normalizeTags(transformed.Tags)
	return transformed
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
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}
	var sum float64
	var min, max float64
	if len(records) > 0 {
		min = records[0].Value
		max = records[0].Value
	}
	for _, record := range records {
		sum += record.Value
		if record.Value < min {
			min = record.Value
		}
		if record.Value > max {
			max = record.Value
		}
	}
	average := sum / float64(len(records))
	return average, max - min
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

func ProcessUserData(rawData []byte) (*UserData, error) {
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
	jsonData := []byte(`{"email":"test@example.com","username":"  john_doe  ","age":25}`)
	processed, err := ProcessUserData(jsonData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", processed)
}