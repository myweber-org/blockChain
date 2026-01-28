
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