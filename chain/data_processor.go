package data_processor

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID        int
	Name      string
	Value     float64
	Validated bool
}

func ParseCSVData(reader io.Reader) ([]DataRecord, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRecord
	for i, row := range records {
		if len(row) < 3 {
			return nil, errors.New("insufficient columns in row " + strconv.Itoa(i))
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, errors.New("invalid ID in row " + strconv.Itoa(i))
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, errors.New("empty name in row " + strconv.Itoa(i))
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, errors.New("invalid value in row " + strconv.Itoa(i))
		}

		validated := false
		if len(row) > 3 {
			validated = strings.ToLower(strings.TrimSpace(row[3])) == "true"
		}

		data = append(data, DataRecord{
			ID:        id,
			Name:      name,
			Value:     value,
			Validated: validated,
		})
	}

	return data, nil
}

func ValidateRecords(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		if record.ID > 0 && record.Value >= 0 {
			record.Validated = true
			validated = append(validated, record)
		}
	}
	return validated
}

func CalculateTotal(records []DataRecord) float64 {
	var total float64
	for _, record := range records {
		if record.Validated {
			total += record.Value
		}
	}
	return total
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
	err := json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
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
	jsonStr := `{"id": 101, "name": "Alice", "email": "alice@example.com"}`
	user, err := ValidateAndParseJSON([]byte(jsonStr))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Parsed user: %+v\n", user)
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

func NormalizeUserData(data UserData) UserData {
	return UserData{
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Username: strings.TrimSpace(data.Username),
		Age:      data.Age,
	}
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	data := UserData{
		Email:    email,
		Username: username,
		Age:      age,
	}
	
	normalizedData := NormalizeUserData(data)
	
	if err := ValidateUserData(normalizedData); err != nil {
		return UserData{}, err
	}
	
	return normalizedData, nil
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type CSVProcessor struct {
	FilePath string
	Headers  []string
	Data     [][]string
}

func NewCSVProcessor(filePath string) *CSVProcessor {
	return &CSVProcessor{
		FilePath: filePath,
	}
}

func (cp *CSVProcessor) Load() error {
	file, err := os.Open(cp.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 1 {
		return fmt.Errorf("empty CSV file")
	}

	cp.Headers = records[0]
	cp.Data = records[1:]
	return nil
}

func (cp *CSVProcessor) Validate() []string {
	var errors []string
	for i, row := range cp.Data {
		if len(row) != len(cp.Headers) {
			errors = append(errors, fmt.Sprintf("row %d: column count mismatch", i+1))
		}
		for j, cell := range row {
			if strings.TrimSpace(cell) == "" {
				errors = append(errors, fmt.Sprintf("row %d, column %d: empty value", i+1, j+1))
			}
		}
	}
	return errors
}

func (cp *CSVProcessor) Clean() {
	var cleanedData [][]string
	for _, row := range cp.Data {
		cleanedRow := make([]string, len(row))
		for j, cell := range row {
			cleanedRow[j] = strings.TrimSpace(cell)
		}
		cleanedData = append(cleanedData, cleanedRow)
	}
	cp.Data = cleanedData
}

func (cp *CSVProcessor) Save(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(cp.Headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for _, row := range cp.Data {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}
	return nil
}

func (cp *CSVProcessor) PrintSummary() {
	fmt.Printf("File: %s\n", cp.FilePath)
	fmt.Printf("Headers: %d\n", len(cp.Headers))
	fmt.Printf("Data rows: %d\n", len(cp.Data))
	fmt.Printf("Total columns: %d\n", len(cp.Headers))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <input.csv> [output.csv]")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := "cleaned_" + inputFile
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}

	processor := NewCSVProcessor(inputFile)
	if err := processor.Load(); err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		os.Exit(1)
	}

	processor.PrintSummary()

	validationErrors := processor.Validate()
	if len(validationErrors) > 0 {
		fmt.Println("Validation errors found:")
		for _, err := range validationErrors {
			fmt.Printf("  - %s\n", err)
		}
	} else {
		fmt.Println("No validation errors found")
	}

	processor.Clean()
	if err := processor.Save(outputFile); err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Cleaned data saved to: %s\n", outputFile)
}package main

import (
	"encoding/csv"
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
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNumber, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 columns, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid ID: %v", lineNumber, err)
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, fmt.Errorf("line %d: name cannot be empty", lineNumber)
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid value: %v", lineNumber, err)
		}

		active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid active flag: %v", lineNumber, err)
		}

		records = append(records, Record{
			ID:     id,
			Name:   name,
			Value:  value,
			Active: active,
		})
	}

	return records, nil
}

func filterActiveRecords(records []Record) []Record {
	var active []Record
	for _, r := range records {
		if r.Active {
			active = append(active, r)
		}
	}
	return active
}

func calculateTotalValue(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}

func generateReport(records []Record) {
	active := filterActiveRecords(records)
	total := calculateTotalValue(active)

	fmt.Printf("Total Records: %d\n", len(records))
	fmt.Printf("Active Records: %d\n", len(active))
	fmt.Printf("Total Value of Active Records: %.2f\n", total)
	fmt.Println("\nActive Record Details:")
	for _, r := range active {
		fmt.Printf("  ID: %d, Name: %s, Value: %.2f\n", r.ID, r.Name, r.Value)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := parseCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	generateReport(records)
}
package main

import (
    "encoding/csv"
    "errors"
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

func ProcessCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []Record{}
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

        name := row[1]
        if name == "" {
            return nil, fmt.Errorf("empty name at line %d", lineNumber)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func CalculateStats(records []Record) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    min := records[0].Value
    max := records[0].Value

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
}

func ValidateRecords(records []Record) error {
    seenIDs := make(map[int]bool)

    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid ID %d: must be positive", record.ID)
        }

        if seenIDs[record.ID] {
            return fmt.Errorf("duplicate ID %d found", record.ID)
        }
        seenIDs[record.ID] = true

        if record.Value < 0 {
            return fmt.Errorf("negative value %f for record ID %d", record.Value, record.ID)
        }
    }

    return nil
}package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_\-@. ]+$`)
	return &DataProcessor{allowedPattern: pattern}
}

func (dp *DataProcessor) SanitizeInput(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	if !dp.allowedPattern.MatchString(trimmed) {
		return "", false
	}

	return trimmed, true
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (map[string]string, bool) {
	sanitizedName, ok := dp.SanitizeInput(name)
	if !ok {
		return nil, false
	}

	if !dp.ValidateEmail(email) {
		return nil, false
	}

	result := map[string]string{
		"name":  sanitizedName,
		"email": strings.ToLower(email),
	}

	return result, true
}
package main

import (
	"errors"
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

func TransformRecord(record DataRecord, multiplier float64) DataRecord {
	if multiplier <= 0 {
		multiplier = 1.0
	}
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp,
		Tags:      append(record.Tags, "processed"),
	}
}

func FilterRecords(records []DataRecord, minValue float64) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Value >= minValue {
			filtered = append(filtered, record)
		}
	}
	return filtered
}