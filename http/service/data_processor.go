package main

import (
	"fmt"
	"strings"
	"unicode"
)

type UserProfile struct {
	Username string
	Email    string
	Age      int
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateAge(age int) bool {
	return age >= 0 && age <= 120
}

func processUserData(username, email string, age int) (*UserProfile, error) {
	normalizedUsername := normalizeUsername(username)
	if normalizedUsername == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	if !validateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}
	normalizedEmail := normalizeEmail(email)

	if !validateAge(age) {
		return nil, fmt.Errorf("age must be between 0 and 120")
	}

	return &UserProfile{
		Username: normalizedUsername,
		Email:    normalizedEmail,
		Age:      age,
	}, nil
}

func main() {
	user, err := processUserData("  JohnDoe  ", "JOHN@EXAMPLE.COM", 30)
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

        if value < 0 {
            return nil, fmt.Errorf("negative value at line %d: %f", lineNumber, value)
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

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var max float64
    count := len(records)

    for i, record := range records {
        sum += record.Value
        if i == 0 || record.Value > max {
            max = record.Value
        }
    }

    average := sum / float64(count)
    return average, max, count
}