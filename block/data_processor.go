
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

func ProcessCSVFile(filepath string) ([]DataRecord, error) {
    file, err := os.Open(filepath)
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

func ValidateRecords(records []DataRecord) error {
    seenIDs := make(map[int]bool)

    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid ID: %d must be positive", record.ID)
        }

        if seenIDs[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        seenIDs[record.ID] = true

        if record.Value < 0 {
            return fmt.Errorf("negative value not allowed for record ID %d", record.ID)
        }
    }

    return nil
}
package main

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
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}

	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Email:    strings.TrimSpace(email),
		Username: transformedUsername,
		Age:      age,
	}

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}