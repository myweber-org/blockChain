
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
    Delimiter  rune
}

func NewDataProcessor(input, output string) *DataProcessor {
    return &DataProcessor{
        InputPath:  input,
        OutputPath: output,
        Delimiter:  ',',
    }
}

func (dp *DataProcessor) ValidateRow(row []string) bool {
    if len(row) == 0 {
        return false
    }
    for _, field := range row {
        if strings.TrimSpace(field) == "" {
            return false
        }
    }
    return true
}

func (dp *DataProcessor) CleanField(field string) string {
    cleaned := strings.TrimSpace(field)
    cleaned = strings.ToUpper(cleaned)
    return cleaned
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
    reader.Comma = dp.Delimiter

    writer := csv.NewWriter(outputFile)
    writer.Comma = dp.Delimiter
    defer writer.Flush()

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV: %w", err)
        }

        if !dp.ValidateRow(record) {
            continue
        }

        cleanedRecord := make([]string, len(record))
        for i, field := range record {
            cleanedRecord[i] = dp.CleanField(field)
        }

        if err := writer.Write(cleanedRecord); err != nil {
            return fmt.Errorf("error writing CSV: %w", err)
        }
    }

    return nil
}

func main() {
    processor := NewDataProcessor("input.csv", "output.csv")
    if err := processor.Process(); err != nil {
        fmt.Printf("Processing failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Data processing completed successfully")
}package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
}

func ValidateUser(data UserData) error {
	if data.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", data.ID)
	}
	if strings.TrimSpace(data.Name) == "" {
		return fmt.Errorf("user name cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return fmt.Errorf("invalid email format: %s", data.Email)
	}
	return nil
}

func TransformUserName(data *UserData) {
	data.Name = strings.ToUpper(strings.TrimSpace(data.Name))
}

func ProcessUserInput(data UserData) (UserData, error) {
	if err := ValidateUser(data); err != nil {
		return UserData{}, err
	}
	TransformUserName(&data)
	return data, nil
}

func main() {
	user := UserData{
		ID:    101,
		Name:  "  john doe  ",
		Email: "john@example.com",
	}
	
	processed, err := ProcessUserInput(user)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Processed user: %+v\n", processed)
}
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TransformToTitleCase(input string) string {
	if len(input) == 0 {
		return input
	}
	return strings.ToUpper(input[:1]) + strings.ToLower(input[1:])
}

func PrettyPrintJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func main() {
	email := "test@example.com"
	fmt.Printf("Email validation result: %v\n", ValidateEmail(email))

	name := "john doe"
	fmt.Printf("Transformed name: %s\n", TransformToTitleCase(name))

	sampleData := map[string]interface{}{
		"id":    1,
		"name":  "Sample Item",
		"price": 29.99,
	}
	prettyJSON, err := PrettyPrintJSON(sampleData)
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
	} else {
		fmt.Printf("Formatted JSON:\n%s\n", prettyJSON)
	}
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
	Timestamp string `json:"timestamp"`
}

func ValidateUsername(username string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,20}$`, username)
	return matched
}

func ValidateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}

func TransformProfile(profile *UserProfile) error {
	profile.Username = strings.TrimSpace(profile.Username)
	profile.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	
	if profile.Age < 0 {
		profile.Age = 0
	}
	
	if !ValidateUsername(profile.Username) {
		return fmt.Errorf("invalid username format")
	}
	
	if !ValidateEmail(profile.Email) {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

func ProcessUserData(inputJSON string) (string, error) {
	var profile UserProfile
	
	err := json.Unmarshal([]byte(inputJSON), &profile)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}
	
	err = TransformProfile(&profile)
	if err != nil {
		return "", fmt.Errorf("validation failed: %v", err)
	}
	
	profile.Active = true
	
	outputJSON, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to generate JSON: %v", err)
	}
	
	return string(outputJSON), nil
}

func main() {
	sampleInput := `{
		"id": 123,
		"username": "  John_Doe  ",
		"email": "  JOHN@EXAMPLE.COM  ",
		"age": 25,
		"active": false,
		"timestamp": "2024-01-15T10:30:00Z"
	}`
	
	result, err := ProcessUserData(sampleInput)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Processed User Profile:")
	fmt.Println(result)
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
	var records []Record

	for i := 0; ; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		if i == 0 {
			continue
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid row length at line %d", i+1)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", i+1, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", i+1)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", i+1, err)
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
	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID %d", record.ID)
		}
		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID %d", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value for record %d", record.ID)
		}
	}
	return nil
}

func CalculateStatistics(records []Record) (float64, float64, float64) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var min, max float64
	first := true

	for _, record := range records {
		sum += record.Value
		if first {
			min = record.Value
			max = record.Value
			first = false
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