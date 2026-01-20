
package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	
	// Trim spaces from start and end
	cleaned = strings.TrimSpace(cleaned)
	
	// Convert to lowercase for normalization
	cleaned = strings.ToLower(cleaned)
	
	return cleaned
}

func NormalizeString(input string) string {
	cleaned := CleanInput(input)
	
	// Remove special characters except alphanumeric and spaces
	re := regexp.MustCompile(`[^a-z0-9\s]`)
	normalized := re.ReplaceAllString(cleaned, "")
	
	return normalized
}

func ProcessData(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		processed := NormalizeString(input)
		results = append(results, processed)
	}
	return results
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
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func TransformUserData(rawData []byte) (*UserData, error) {
	var user UserData
	err := json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	if !ValidateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format: %s", user.Email)
	}

	user.Username = SanitizeUsername(user.Username)

	if user.Age < 0 || user.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", user.Age)
	}

	return &user, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  JohnDoe  ","age":25}`
	user, err := TransformUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
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

type DataRecord struct {
    ID    int
    Name  string
    Value float64
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

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
            return nil, fmt.Errorf("invalid ID format at line %d: %w", lineNumber, err)
        }

        name := row[1]
        if name == "" {
            return nil, fmt.Errorf("empty name at line %d", lineNumber)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value format at line %d: %w", lineNumber, err)
        }

        records = append(records, DataRecord{
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

func CalculateStatistics(records []DataRecord) (float64, float64) {
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

func ValidateRecords(records []DataRecord) []error {
    var errors []error
    seenIDs := make(map[int]bool)

    for _, record := range records {
        if record.ID <= 0 {
            errors = append(errors, fmt.Errorf("invalid ID %d: must be positive", record.ID))
        }

        if seenIDs[record.ID] {
            errors = append(errors, fmt.Errorf("duplicate ID %d found", record.ID))
        }
        seenIDs[record.ID] = true

        if record.Value < 0 {
            errors = append(errors, fmt.Errorf("record ID %d has negative value: %f", record.ID, record.Value))
        }
    }

    return errors
}