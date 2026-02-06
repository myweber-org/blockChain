
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
package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}package main

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

func ProcessCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid row length: %d", len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %w", err)
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value format: %w", err)
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
			return fmt.Errorf("invalid ID: %d", record.ID)
		}
		if record.Name == "" {
			return fmt.Errorf("empty name for ID: %d", record.ID)
		}
		if record.Value < 0 {
			return fmt.Errorf("negative value for ID: %d", record.ID)
		}
		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID: %d", record.ID)
		}
		seenIDs[record.ID] = true
	}
	return nil
}

func CalculateTotalValue(records []Record) float64 {
	var total float64
	for _, record := range records {
		total += record.Value
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
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []DataRecord{}
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

func ValidateRecord(record DataRecord) error {
    if record.ID <= 0 {
        return errors.New("ID must be positive")
    }
    if record.Name == "" {
        return errors.New("name cannot be empty")
    }
    if record.Value < 0 {
        return errors.New("value cannot be negative")
    }
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        os.Exit(1)
    }

    records, err := ProcessCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    avg, max, count := CalculateStatistics(records)
    fmt.Printf("Processed %d records\n", count)
    fmt.Printf("Average value: %.2f\n", avg)
    fmt.Printf("Maximum value: %.2f\n", max)

    fmt.Println("\nRecord validation:")
    for _, record := range records {
        if err := ValidateRecord(record); err != nil {
            fmt.Printf("Record %d invalid: %v\n", record.ID, err)
        }
    }
}package data

import (
	"regexp"
	"strings"
)

var (
	alphaNumericRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	whitespaceRegex   = regexp.MustCompile(`\s+`)
)

func SanitizeInput(input string, maxLength int) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	if len(trimmed) > maxLength {
		trimmed = trimmed[:maxLength]
	}

	if !alphaNumericRegex.MatchString(trimmed) {
		return "", false
	}

	sanitized := whitespaceRegex.ReplaceAllString(trimmed, " ")
	return sanitized, true
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
	BirthDate string
	CreatedAt time.Time
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func CalculateAge(birthDate string) (int, error) {
	parsedDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	now := time.Now()
	age := now.Year() - parsedDate.Year()

	if now.YearDay() < parsedDate.YearDay() {
		age--
	}

	if age < 0 {
		return 0, errors.New("birth date cannot be in the future")
	}

	return age, nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateEmail(profile.Email); err != nil {
		return profile, err
	}

	profile.Username = SanitizeUsername(profile.Username)

	age, err := CalculateAge(profile.BirthDate)
	if err != nil {
		return profile, err
	}

	if age < 13 {
		return profile, errors.New("user must be at least 13 years old")
	}

	return profile, nil
}package main

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
			Valid: validateRecord(row[0], row[1], row[2]),
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(id, name, email string) bool {
	if id == "" || name == "" || email == "" {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	return true
}

func GenerateReport(records []DataRecord) {
	validCount := 0
	invalidCount := 0

	fmt.Println("Data Processing Report")
	fmt.Println("======================")

	for _, record := range records {
		if record.Valid {
			validCount++
			fmt.Printf("✓ Valid: %s - %s\n", record.ID, record.Name)
		} else {
			invalidCount++
			fmt.Printf("✗ Invalid: %s - %s (Email: %s)\n", record.ID, record.Name, record.Email)
		}
	}

	fmt.Printf("\nSummary: %d valid, %d invalid records\n", validCount, invalidCount)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	GenerateReport(records)
}
package main

import (
    "encoding/csv"
    "errors"
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

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
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

        if err := validateRecord(record); err != nil {
            return nil, fmt.Errorf("validation failed at line %d: %w", lineNumber, err)
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func validateRecord(record DataRecord) error {
    if record.ID == "" {
        return errors.New("ID cannot be empty")
    }
    if record.Name == "" {
        return errors.New("name cannot be empty")
    }
    if !strings.Contains(record.Email, "@") {
        return errors.New("invalid email format")
    }
    if record.Active != "true" && record.Active != "false" {
        return errors.New("active status must be 'true' or 'false'")
    }
    return nil
}

func GenerateReport(records []DataRecord) {
    activeCount := 0
    for _, record := range records {
        if record.Active == "true" {
            activeCount++
        }
    }

    fmt.Printf("Total Records Processed: %d\n", len(records))
    fmt.Printf("Active Records: %d\n", activeCount)
    fmt.Printf("Inactive Records: %d\n", len(records)-activeCount)
}