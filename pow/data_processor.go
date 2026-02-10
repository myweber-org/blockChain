
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Username string
	Email    string
	Age      int
}

func NormalizeUsername(username string) (string, error) {
	trimmed := strings.TrimSpace(username)
	if trimmed == "" {
		return "", errors.New("username cannot be empty")
	}
	if len(trimmed) < 3 {
		return "", errors.New("username must be at least 3 characters")
	}
	return strings.ToLower(trimmed), nil
}

func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	normalizedUsername, err := NormalizeUsername(profile.Username)
	if err != nil {
		return UserProfile{}, err
	}
	profile.Username = normalizedUsername

	if err := ValidateEmail(profile.Email); err != nil {
		return UserProfile{}, err
	}

	if profile.Age < 0 || profile.Age > 150 {
		return UserProfile{}, errors.New("age must be between 0 and 150")
	}

	return profile, nil
}
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
    var total float64
    for _, record := range records {
        total += record.Value
    }
    return total
}