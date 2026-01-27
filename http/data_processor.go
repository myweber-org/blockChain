
package main

import (
    "encoding/csv"
    "errors"
    "io"
    "os"
    "strconv"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
}

func ReadCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []DataRecord

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if len(row) != 3 {
            return nil, errors.New("invalid csv format")
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, err
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, err
        }

        record := DataRecord{
            ID:    id,
            Name:  row[1],
            Value: value,
        }
        records = append(records, record)
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) error {
    for _, record := range records {
        if record.ID <= 0 {
            return errors.New("invalid id")
        }
        if record.Name == "" {
            return errors.New("empty name")
        }
        if record.Value < 0 {
            return errors.New("negative value")
        }
    }
    return nil
}

func CalculateTotal(records []DataRecord) float64 {
    total := 0.0
    for _, record := range records {
        total += record.Value
    }
    return total
}package main

import (
	"fmt"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return len(email) > 5
}

func ValidateAge(age int) bool {
	return age >= 0 && age <= 120
}

func ProcessUserInput(username, email string, age int) (*UserData, error) {
	normalizedUsername := NormalizeUsername(username)
	if normalizedUsername == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	for _, r := range normalizedUsername {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' && r != '-' {
			return nil, fmt.Errorf("username contains invalid characters")
		}
	}

	if !ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	if !ValidateAge(age) {
		return nil, fmt.Errorf("age must be between 0 and 120")
	}

	return &UserData{
		Username: normalizedUsername,
		Email:    strings.ToLower(email),
		Age:      age,
	}, nil
}

func main() {
	user, err := ProcessUserInput("  john_doe123  ", "John@Example.COM", 30)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}