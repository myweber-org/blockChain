
package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func readCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func filterRows(records [][]string, filterFunc func([]string) bool) [][]string {
	var filtered [][]string
	for _, row := range records {
		if filterFunc(row) {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func containsKeyword(row []string, keyword string) bool {
	for _, field := range row {
		if strings.Contains(strings.ToLower(field), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func writeCSVFile(filePath string, records [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func processCSVData(inputPath, outputPath, keyword string) error {
	records, err := readCSVFile(inputPath)
	if err != nil {
		return err
	}

	filtered := filterRows(records, func(row []string) bool {
		return containsKeyword(row, keyword)
	})

	if len(filtered) == 0 {
		log.Println("No matching records found")
	}

	return writeCSVFile(outputPath, filtered)
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

    if len(records) == 0 {
        return nil, fmt.Errorf("no valid records found in file")
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
}package main

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
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

func TrimAndLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func RemoveSpecialChars(input string) string {
	re := regexp.MustCompile(`[^\w\s]`)
	return re.ReplaceAllString(input, "")
}