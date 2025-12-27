
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

func FilterActiveUsers(records []DataRecord) []DataRecord {
	var activeUsers []DataRecord
	for _, record := range records {
		if record.Active == "true" {
			activeUsers = append(activeUsers, record)
		}
	}
	return activeUsers
}

func GenerateReport(records []DataRecord) {
	fmt.Printf("Total records processed: %d\n", len(records))
	activeUsers := FilterActiveUsers(records)
	fmt.Printf("Active users: %d\n", len(activeUsers))
	fmt.Printf("Inactive users: %d\n", len(records)-len(activeUsers))
}package main

import (
	"fmt"
)

// FilterAndDouble filters out even numbers from the input slice,
// doubles the remaining odd numbers, and returns the new slice.
func FilterAndDouble(numbers []int) []int {
	var result []int
	for _, num := range numbers {
		if num%2 != 0 {
			result = append(result, num*2)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := FilterAndDouble(input)
	fmt.Println("Original:", input)
	fmt.Println("Filtered and Doubled:", output)
}
package data_processor

func FilterAndTransform(numbers []int, threshold int, transformFunc func(int) int) []int {
	var result []int
	for _, num := range numbers {
		if num > threshold {
			transformed := transformFunc(num)
			result = append(result, transformed)
		}
	}
	return result
}

func DoubleValue(x int) int {
	return x * 2
}
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

// ValidateJSONString checks if a string is valid JSON.
func ValidateJSONString(s string) bool {
    var js interface{}
    return json.Unmarshal([]byte(s), &js) == nil
}

// PrettyPrintJSON takes a JSON string and returns a formatted version.
func PrettyPrintJSON(s string) (string, error) {
    var data interface{}
    err := json.Unmarshal([]byte(s), &data)
    if err != nil {
        return "", err
    }
    pretty, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return "", err
    }
    return string(pretty), nil
}

// ExtractJSONField extracts the value of a top-level field from a JSON string.
func ExtractJSONField(s, field string) (string, error) {
    var data map[string]interface{}
    err := json.Unmarshal([]byte(s), &data)
    if err != nil {
        return "", err
    }
    value, exists := data[field]
    if !exists {
        return "", fmt.Errorf("field '%s' not found", field)
    }
    result, err := json.Marshal(value)
    if err != nil {
        return "", err
    }
    return strings.Trim(string(result), `"`), nil
}

func main() {
    // Example usage
    jsonStr := `{"name":"Alice","age":30,"active":true}`
    fmt.Println("Valid JSON?", ValidateJSONString(jsonStr))

    pretty, _ := PrettyPrintJSON(jsonStr)
    fmt.Println("Pretty JSON:\n", pretty)

    name, _ := ExtractJSONField(jsonStr, "name")
    fmt.Println("Extracted name:", name)
}package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ParseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		lineNum++

		if len(line) != 3 {
			return nil, errors.New("invalid column count at line " + strconv.Itoa(lineNum))
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, errors.New("invalid ID at line " + strconv.Itoa(lineNum))
		}

		name := line[1]
		if name == "" {
			return nil, errors.New("empty name at line " + strconv.Itoa(lineNum))
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, errors.New("invalid value at line " + strconv.Itoa(lineNum))
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, rec := range records {
		if rec.ID <= 0 {
			return errors.New("record ID must be positive: " + strconv.Itoa(rec.ID))
		}
		if seenIDs[rec.ID] {
			return errors.New("duplicate ID found: " + strconv.Itoa(rec.ID))
		}
		seenIDs[rec.ID] = true
	}
	return nil
}

func CalculateTotalValue(records []Record) float64 {
	total := 0.0
	for _, rec := range records {
		total += rec.Value
	}
	return total
}