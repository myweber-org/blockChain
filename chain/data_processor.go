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

		if !strings.Contains(record.Email, "@") {
			return nil, fmt.Errorf("invalid email format at line %d", lineNumber)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) []string {
	var errors []string
	emailSet := make(map[string]bool)

	for i, record := range records {
		if record.Active != "true" && record.Active != "false" {
			errors = append(errors, fmt.Sprintf("record %d: active field must be 'true' or 'false'", i+1))
		}

		if emailSet[record.Email] {
			errors = append(errors, fmt.Sprintf("record %d: duplicate email detected", i+1))
		}
		emailSet[record.Email] = true
	}

	return errors
}

func GenerateReport(records []DataRecord) {
	activeCount := 0
	for _, record := range records {
		if record.Active == "true" {
			activeCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Active records: %d\n", activeCount)
	fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
}
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

func processCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
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

		if len(row) < 4 {
			continue
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if isValidRecord(record) {
			records = append(records, record)
		}
	}

	return records, nil
}

func isValidRecord(record DataRecord) bool {
	if record.ID == "" || record.Name == "" {
		return false
	}
	if !strings.Contains(record.Email, "@") {
		return false
	}
	return record.Active == "true" || record.Active == "false"
}

func generateReport(records []DataRecord) {
	activeCount := 0
	for _, record := range records {
		if record.Active == "true" {
			activeCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Active records: %d\n", activeCount)
	fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	generateReport(records)
}
package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int
	Email     string
	Username  string
	BirthDate time.Time
	Active    bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if len(strings.TrimSpace(profile.Username)) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if time.Since(profile.BirthDate).Hours()/24/365 < 13 {
		return errors.New("user must be at least 13 years old")
	}

	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func CalculateAge(birthDate time.Time) int {
	now := time.Now()
	years := now.Year() - birthDate.Year()

	if now.YearDay() < birthDate.YearDay() {
		years--
	}

	return years
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateUserProfile(profile); err != nil {
		return UserProfile{}, err
	}

	profile.Username = TransformUsername(profile.Username)
	return profile, nil
}