package main

import (
    "encoding/csv"
    "errors"
    "io"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
    Valid bool
}

func ParseCSVData(reader io.Reader) ([]DataRecord, error) {
    csvReader := csv.NewReader(reader)
    records, err := csvReader.ReadAll()
    if err != nil {
        return nil, err
    }

    if len(records) == 0 {
        return []DataRecord{}, nil
    }

    var data []DataRecord
    for i, row := range records {
        if len(row) < 4 {
            continue
        }

        id, err := strconv.Atoi(strings.TrimSpace(row[0]))
        if err != nil {
            continue
        }

        name := strings.TrimSpace(row[1])
        if name == "" {
            continue
        }

        value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
        if err != nil {
            continue
        }

        valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

        data = append(data, DataRecord{
            ID:    id,
            Name:  name,
            Value: value,
            Valid: valid,
        })
    }

    return data, nil
}

func ValidateData(records []DataRecord) error {
    if len(records) == 0 {
        return errors.New("no data records provided")
    }

    idSet := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return errors.New("invalid ID found: " + strconv.Itoa(record.ID))
        }

        if idSet[record.ID] {
            return errors.New("duplicate ID found: " + strconv.Itoa(record.ID))
        }
        idSet[record.ID] = true

        if record.Value < 0 {
            return errors.New("negative value not allowed for ID: " + strconv.Itoa(record.ID))
        }
    }

    return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var validCount int
    var maxValue float64

    for _, record := range records {
        if record.Valid {
            sum += record.Value
            validCount++
            if record.Value > maxValue {
                maxValue = record.Value
            }
        }
    }

    if validCount == 0 {
        return 0, 0, 0
    }

    average := sum / float64(validCount)
    return average, maxValue, validCount
}package main

import (
	"errors"
	"strings"
	"unicode"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username can only contain letters, digits, underscores, and hyphens")
		}
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TransformUserData(rawUsername, rawEmail string) (string, string, error) {
	if err := ValidateUsername(rawUsername); err != nil {
		return "", "", err
	}
	normalizedEmail := NormalizeEmail(rawEmail)
	return rawUsername, normalizedEmail, nil
}