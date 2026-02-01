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
	err := ValidateUserData(userData)
	if err != nil {
		return UserData{}, err
	}
	return userData, nil
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

func CalculateStatistics(records []DataRecord) (float64, float64, float64) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var min, max float64
	min = records[0].Value
	max = records[0].Value

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
	return min, max, average
}package main

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
	line := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		line++
		if len(row) != 3 {
			return nil, fmt.Errorf("invalid columns at line %d", line)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", line, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", line)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", line, err)
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

	records, err := ProcessCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d records\n", len(records))
	avg, max := CalculateStats(records)
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
}
package main

import (
	"fmt"
	"sort"
	"time"
)

type Record struct {
	ID        int
	Name      string
	Timestamp time.Time
	Value     float64
}

func filterAndSortByDate(records []Record, startDate, endDate time.Time) []Record {
	var filtered []Record
	for _, r := range records {
		if !r.Timestamp.Before(startDate) && !r.Timestamp.After(endDate) {
			filtered = append(filtered, r)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	return filtered
}

func main() {
	records := []Record{
		{1, "Alpha", time.Date(2023, 5, 10, 0, 0, 0, 0, time.UTC), 12.5},
		{2, "Beta", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), 8.3},
		{3, "Gamma", time.Date(2023, 5, 5, 0, 0, 0, 0, time.UTC), 15.7},
		{4, "Delta", time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC), 9.1},
	}

	start := time.Date(2023, 5, 8, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 5, 18, 0, 0, 0, 0, time.UTC)

	result := filterAndSortByDate(records, start, end)

	for _, r := range result {
		fmt.Printf("ID: %d, Name: %s, Date: %s, Value: %.1f\n",
			r.ID, r.Name, r.Timestamp.Format("2006-01-02"), r.Value)
	}
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
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, err
	}

	transformed := DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp,
		Tags:      append([]string{}, record.Tags...),
	}

	transformed.Tags = append(transformed.Tags, fmt.Sprintf("processed_%d", time.Now().Unix()))
	return transformed, nil
}

func ProcessBatch(records []DataRecord, multiplier float64) ([]DataRecord, error) {
	var processed []DataRecord
	var errors []string

	for i, record := range records {
		transformed, err := TransformRecord(record, multiplier)
		if err != nil {
			errors = append(errors, fmt.Sprintf("record %d: %v", i, err))
			continue
		}
		processed = append(processed, transformed)
	}

	if len(errors) > 0 {
		return processed, fmt.Errorf("processing completed with errors: %s", strings.Join(errors, "; "))
	}
	return processed, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records provided for statistics calculation")
	}

	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	mean := sum / float64(len(records))

	var varianceSum float64
	for _, record := range records {
		diff := record.Value - mean
		varianceSum += diff * diff
	}
	variance := varianceSum / float64(len(records))

	return mean, variance, nil
}package main

import (
	"fmt"
)

// DataProcessor filters and transforms a slice of integers.
func DataProcessor(input []int) []int {
	var result []int
	for _, v := range input {
		if v > 0 {
			result = append(result, v*2)
		}
	}
	return result
}

func main() {
	data := []int{-5, 2, 0, 8, -1, 10}
	processed := DataProcessor(data)
	fmt.Println("Processed data:", processed)
}
package main

import (
	"errors"
	"strings"
)

// ValidateEmail checks if the provided string is a valid email address.
func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

// SanitizeInput removes leading and trailing whitespace from a string.
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// ParsePositiveInteger converts a string to a positive integer.
func ParsePositiveInteger(s string) (int, error) {
	sanitized := SanitizeInput(s)
	var result int
	for _, ch := range sanitized {
		if ch < '0' || ch > '9' {
			return 0, errors.New("invalid character in number")
		}
		result = result*10 + int(ch-'0')
	}
	if result <= 0 {
		return 0, errors.New("number must be positive")
	}
	return result, nil
}
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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func processCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record
	lineNumber := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		record := Record{
			ID:    id,
			Name:  row[1],
			Value: row[2],
		}
		records = append(records, record)
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

func saveJSONToFile(data string, filename string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <input_csv_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	records, err := processCSVFile(inputFile)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	jsonOutput, err := convertToJSON(records)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	outputFile := "output.json"
	err = saveJSONToFile(jsonOutput, outputFile)
	if err != nil {
		fmt.Printf("Error saving JSON file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	fmt.Printf("JSON output saved to %s\n", outputFile)
}package main

import "fmt"

func calculateMovingAverage(data []float64, windowSize int) []float64 {
    if len(data) == 0 || windowSize <= 0 || windowSize > len(data) {
        return []float64{}
    }

    result := make([]float64, len(data)-windowSize+1)
    for i := 0; i <= len(data)-windowSize; i++ {
        sum := 0.0
        for j := 0; j < windowSize; j++ {
            sum += data[i+j]
        }
        result[i] = sum / float64(windowSize)
    }
    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    window := 3
    averages := calculateMovingAverage(sampleData, window)
    
    fmt.Printf("Original data: %v\n", sampleData)
    fmt.Printf("Moving averages (window=%d): %v\n", window, averages)
}