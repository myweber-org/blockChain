package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID    string
	Name  string
	Email string
	Valid bool
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
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

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Email: strings.TrimSpace(row[2]),
			Valid: validateRecord(strings.TrimSpace(row[0]), strings.TrimSpace(row[2])),
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(id, email string) bool {
	if id == "" || email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file_path>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	validRecords := FilterValidRecords(records)
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", len(validRecords))

	for _, record := range validRecords {
		fmt.Printf("ID: %s, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
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
	return strings.TrimSpace(username)
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
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
	rawJSON := `{"email":"test@example.com","username":"  john_doe  ","age":25}`
	processedData, err := ProcessUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
}
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func ProcessUserData(rawData []byte) (*UserData, error) {
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
	if user.Age < 0 || user.Age > 120 {
		return nil, fmt.Errorf("invalid age value: %d", user.Age)
	}

	user.Email = sanitizeEmail(user.Email)
	return &user, nil
}

func sanitizeEmail(email string) string {
	// Simple email normalization
	return email
}

func main() {
	jsonData := `{"id": 123, "name": "John Doe", "email": "john@example.com", "age": 30}`
	user, err := ProcessUserData([]byte(jsonData))
	if err != nil {
		log.Fatalf("Error processing data: %v", err)
	}
	fmt.Printf("Processed user: %+v\n", user)
}package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Value string
}

type Processor interface {
	Process(DataRecord) (DataRecord, error)
}

type Validator struct{}

func (v Validator) Process(record DataRecord) (DataRecord, error) {
	if record.ID == "" {
		return record, errors.New("empty ID field")
	}
	if len(record.Value) < 3 {
		return record, errors.New("value too short")
	}
	return record, nil
}

type Transformer struct{}

func (t Transformer) Process(record DataRecord) (DataRecord, error) {
	record.Value = strings.ToUpper(record.Value)
	return record, nil
}

type Pipeline struct {
	processors []Processor
}

func (p *Pipeline) AddProcessor(proc Processor) {
	p.processors = append(p.processors, proc)
}

func (p *Pipeline) Execute(record DataRecord) (DataRecord, error) {
	var err error
	for _, proc := range p.processors {
		record, err = proc.Process(record)
		if err != nil {
			return record, err
		}
	}
	return record, nil
}

func main() {
	pipeline := &Pipeline{}
	pipeline.AddProcessor(Validator{})
	pipeline.AddProcessor(Transformer{})

	testRecords := []DataRecord{
		{ID: "001", Value: "test"},
		{ID: "", Value: "data"},
		{ID: "003", Value: "ab"},
	}

	for _, record := range testRecords {
		result, err := pipeline.Execute(record)
		if err != nil {
			fmt.Printf("Error processing %v: %v\n", record.ID, err)
		} else {
			fmt.Printf("Processed: ID=%s, Value=%s\n", result.ID, result.Value)
		}
	}
}