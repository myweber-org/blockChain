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
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
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
	jsonData := []byte(`{"email":"test@example.com","username":"  JohnDoe  ","age":25}`)
	processedData, err := ProcessUserData(jsonData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
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

type DataProcessor struct {
	records []Record
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		records: make([]Record, 0),
	}
}

func (dp *DataProcessor) LoadCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	lineNumber := 0
	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) != 4 {
			return fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		record, err := parseRecord(row)
		if err != nil {
			return fmt.Errorf("parse error at line %d: %w", lineNumber, err)
		}

		dp.records = append(dp.records, record)
	}

	return nil
}

func parseRecord(row []string) (Record, error) {
	var record Record

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID: %w", err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, errors.New("name cannot be empty")
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value: %w", err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid active flag: %w", err)
	}
	record.Active = active

	return record, nil
}

func (dp *DataProcessor) FilterActive() []Record {
	var active []Record
	for _, record := range dp.records {
		if record.Active {
			active = append(active, record)
		}
	}
	return active
}

func (dp *DataProcessor) CalculateTotal() float64 {
	var total float64
	for _, record := range dp.records {
		total += record.Value
	}
	return total
}

func (dp *DataProcessor) FindByName(name string) *Record {
	for _, record := range dp.records {
		if strings.EqualFold(record.Name, name) {
			return &record
		}
	}
	return nil
}

func (dp *DataProcessor) Statistics() map[string]interface{} {
	if len(dp.records) == 0 {
		return map[string]interface{}{
			"count":   0,
			"average": 0.0,
			"max":     0.0,
			"min":     0.0,
		}
	}

	var total float64
	maxValue := dp.records[0].Value
	minValue := dp.records[0].Value
	activeCount := 0

	for _, record := range dp.records {
		total += record.Value
		if record.Value > maxValue {
			maxValue = record.Value
		}
		if record.Value < minValue {
			minValue = record.Value
		}
		if record.Active {
			activeCount++
		}
	}

	return map[string]interface{}{
		"count":       len(dp.records),
		"average":     total / float64(len(dp.records)),
		"max":         maxValue,
		"min":         minValue,
		"activeCount": activeCount,
	}
}

func main() {
	processor := NewDataProcessor()

	err := processor.LoadCSV("data.csv")
	if err != nil {
		fmt.Printf("Error loading CSV: %v\n", err)
		return
	}

	stats := processor.Statistics()
	fmt.Printf("Data Statistics: %+v\n", stats)

	activeRecords := processor.FilterActive()
	fmt.Printf("Active records: %d\n", len(activeRecords))

	total := processor.CalculateTotal()
	fmt.Printf("Total value: %.2f\n", total)

	found := processor.FindByName("example")
	if found != nil {
		fmt.Printf("Found record: %+v\n", found)
	}
}