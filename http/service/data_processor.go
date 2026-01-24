
package main

import (
    "fmt"
    "strings"
)

type UserData struct {
    Username string
    Email    string
    Age      int
}

func ValidateUserData(data UserData) error {
    if strings.TrimSpace(data.Username) == "" {
        return fmt.Errorf("username cannot be empty")
    }
    if !strings.Contains(data.Email, "@") {
        return fmt.Errorf("invalid email format")
    }
    if data.Age < 0 || data.Age > 150 {
        return fmt.Errorf("age must be between 0 and 150")
    }
    return nil
}

func TransformUsername(data *UserData) {
    data.Username = strings.ToLower(strings.TrimSpace(data.Username))
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
    user := UserData{
        Username: username,
        Email:    email,
        Age:      age,
    }

    TransformUsername(&user)

    if err := ValidateUserData(user); err != nil {
        return UserData{}, err
    }

    return user, nil
}

func main() {
    result, err := ProcessUserInput("  JohnDoe  ", "john@example.com", 30)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Processed user: %+v\n", result)
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	InputPath  string
	OutputPath string
}

func NewDataProcessor(input, output string) *DataProcessor {
	return &DataProcessor{
		InputPath:  input,
		OutputPath: output,
	}
}

func (dp *DataProcessor) Process() error {
	inputFile, err := os.Open(dp.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dp.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	cleanedHeaders := dp.cleanHeaders(headers)
	if err := writer.Write(cleanedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		cleanedRecord := dp.cleanRecord(record)
		if dp.isValidRecord(cleanedRecord) {
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			recordCount++
		}
	}

	fmt.Printf("Processed %d valid records\n", recordCount)
	return nil
}

func (dp *DataProcessor) cleanHeaders(headers []string) []string {
	cleaned := make([]string, len(headers))
	for i, header := range headers {
		cleaned[i] = strings.TrimSpace(strings.ToLower(header))
	}
	return cleaned
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		cleaned[i] = strings.TrimSpace(field)
		if cleaned[i] == "" {
			cleaned[i] = "N/A"
		}
	}
	return cleaned
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
	for _, field := range record {
		if field == "" {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1], os.Args[2])
	if err := processor.Process(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
package main

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

func (dp *DataProcessor) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToLower(trimmed)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, bool) {
	sanitizedName := dp.SanitizeInput(name)
	sanitizedEmail := dp.SanitizeInput(email)

	if !dp.ValidateEmail(sanitizedEmail) {
		return "", false
	}

	result := sanitizedName + " <" + sanitizedEmail + ">"
	return result, true
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
			return errors.New("record ID must be positive")
		}
		if seenIDs[rec.ID] {
			return errors.New("duplicate ID found: " + strconv.Itoa(rec.ID))
		}
		seenIDs[rec.ID] = true
	}
	return nil
}

func CalculateTotal(records []Record) float64 {
	total := 0.0
	for _, rec := range records {
		total += rec.Value
	}
	return total
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

type DataRecord struct {
    ID        int
    Name      string
    Value     float64
    Timestamp string
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

        if len(row) != 4 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
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

func parseRow(row []string, lineNumber int) (DataRecord, error) {
    var record DataRecord

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
    }
    record.ID = id

    name := strings.TrimSpace(row[1])
    if name == "" {
        return record, fmt.Errorf("empty name at line %d", lineNumber)
    }
    record.Name = name

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
    }
    record.Value = value

    timestamp := strings.TrimSpace(row[3])
    if timestamp == "" {
        return record, fmt.Errorf("empty timestamp at line %d", lineNumber)
    }
    record.Timestamp = timestamp

    return record, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, float64) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var min, max float64

    for i, record := range records {
        sum += record.Value

        if i == 0 {
            min = record.Value
            max = record.Value
        } else {
            if record.Value < min {
                min = record.Value
            }
            if record.Value > max {
                max = record.Value
            }
        }
    }

    average := sum / float64(len(records))
    return average, min, max
}

func ValidateRecords(records []DataRecord) []error {
    var errors []error

    seenIDs := make(map[int]bool)
    for _, record := range records {
        if seenIDs[record.ID] {
            errors = append(errors, fmt.Errorf("duplicate ID found: %d", record.ID))
        }
        seenIDs[record.ID] = true

        if record.Value < 0 {
            errors = append(errors, fmt.Errorf("negative value for ID %d: %f", record.ID, record.Value))
        }
    }

    return errors
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

// ParseUserData attempts to parse JSON into a User struct.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ParseUserData(rawData []byte) (*User, error) {
	valid, err := ValidateJSON(rawData)
	if !valid {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return &user, nil
}

func main() {
	jsonStr := `{"id": 101, "name": "Alice", "email": "alice@example.com"}`
	user, err := ParseUserData([]byte(jsonStr))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Parsed User: %+v\n", user)
}