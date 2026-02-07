
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

		records = append(records, record)
	}

	return records, nil
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
	var active []DataRecord
	for _, record := range records {
		if strings.ToLower(record.Active) == "true" || record.Active == "1" {
			active = append(active, record)
		}
	}
	return active
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	activeRecords := FilterActiveRecords(records)
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Active records: %d\n", len(activeRecords))

	for i, record := range activeRecords {
		if !ValidateEmail(record.Email) {
			fmt.Printf("Warning: Invalid email for record %d (ID: %s)\n", i+1, record.ID)
		}
	}
}package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Value string
	Valid bool
}

func ProcessRecords(records []DataRecord) ([]string, error) {
	if len(records) == 0 {
		return nil, errors.New("no records to process")
	}

	var processed []string
	for _, record := range records {
		if !record.Valid {
			continue
		}

		trimmed := strings.TrimSpace(record.Value)
		if trimmed == "" {
			continue
		}

		processed = append(processed, fmt.Sprintf("ID:%d|Value:%s", record.ID, trimmed))
	}

	if len(processed) == 0 {
		return nil, errors.New("no valid records found")
	}

	return processed, nil
}

func ValidateRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("invalid ID")
	}
	if strings.TrimSpace(record.Value) == "" {
		return errors.New("empty value")
	}
	return nil
}

func main() {
	records := []DataRecord{
		{ID: 1, Value: "alpha", Valid: true},
		{ID: 2, Value: "  ", Valid: true},
		{ID: 3, Value: "beta", Valid: false},
		{ID: 4, Value: "gamma", Valid: true},
	}

	result, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Println("Processed records:")
	for _, item := range result {
		fmt.Println(item)
	}
}
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

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func sanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func transformUserData(rawData []byte) (*UserData, error) {
	var user UserData
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	user.Username = sanitizeUsername(user.Username)

	if !validateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format: %s", user.Email)
	}

	if user.Age < 0 || user.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", user.Age)
	}

	return &user, nil
}

func processInput(data []byte) {
	user, err := transformUserData(data)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Printf("Processed user: %+v\n", user)
}

func main() {
	testData := []byte(`{"email":"test@example.com","username":"  john_doe  ","age":25}`)
	processInput(testData)
}
package main

import "fmt"

func FilterAndDoubleEvenNumbers(numbers []int) []int {
    var result []int
    for _, num := range numbers {
        if num%2 == 0 {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    output := FilterAndDoubleEvenNumbers(input)
    fmt.Printf("Input: %v\n", input)
    fmt.Printf("Output: %v\n", output)
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

	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func NormalizeUserData(data *UserData) {
	data.Email = strings.ToLower(strings.TrimSpace(data.Email))
	data.Username = strings.TrimSpace(data.Username)
}