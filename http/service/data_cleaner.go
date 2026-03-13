package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput cleans and normalizes user-provided string input
func SanitizeInput(input string) string {
	// Remove leading/trailing whitespace
	cleaned := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")
	
	// Remove potentially dangerous characters (basic example)
	dangerousChars := regexp.MustCompile(`[<>{}]`)
	cleaned = dangerousChars.ReplaceAllString(cleaned, "")
	
	return cleaned
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Name  string
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) {
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range dc.records {
		key := fmt.Sprintf("%d|%s|%s", record.ID, strings.ToLower(record.Email), strings.ToLower(record.Name))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}

	dc.records = unique
	return unique
}

func (dc *DataCleaner) ValidateEmails() (valid, invalid []DataRecord) {
	for _, record := range dc.records {
		if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}
	return valid, invalid
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(DataRecord{1, "user@example.com", "John Doe"})
	cleaner.AddRecord(DataRecord{2, "user@example.com", "John Doe"})
	cleaner.AddRecord(DataRecord{3, "invalid-email", "Jane Smith"})
	cleaner.AddRecord(DataRecord{4, "admin@domain.org", "Admin User"})

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	unique := cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", len(unique))

	valid, invalid := cleaner.ValidateEmails()
	fmt.Printf("Valid emails: %d, Invalid emails: %d\n", len(valid), len(invalid))

	for _, record := range valid {
		fmt.Printf("Valid: ID=%d, Email=%s, Name=%s\n", record.ID, record.Email, record.Name)
	}
}
package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID        int
    Name      string
    Email     string
    Age       int
    Active    bool
}

func cleanEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}

func validateAge(age int) bool {
    return age >= 0 && age <= 120
}

func parseCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []DataRecord
    lineNumber := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        lineNumber++
        if lineNumber == 1 {
            continue
        }

        if len(line) != 5 {
            continue
        }

        id, err := strconv.Atoi(line[0])
        if err != nil {
            continue
        }

        name := strings.TrimSpace(line[1])
        email := cleanEmail(line[2])

        age, err := strconv.Atoi(line[3])
        if err != nil || !validateAge(age) {
            continue
        }

        active := strings.ToLower(line[4]) == "true"

        record := DataRecord{
            ID:     id,
            Name:   name,
            Email:  email,
            Age:    age,
            Active: active,
        }
        records = append(records, record)
    }

    return records, nil
}

func removeDuplicates(records []DataRecord) []DataRecord {
    seen := make(map[string]bool)
    var unique []DataRecord

    for _, record := range records {
        key := fmt.Sprintf("%s|%s", record.Email, record.Name)
        if !seen[key] {
            seen[key] = true
            unique = append(unique, record)
        }
    }
    return unique
}

func filterActiveUsers(records []DataRecord) []DataRecord {
    var active []DataRecord
    for _, record := range records {
        if record.Active {
            active = append(active, record)
        }
    }
    return active
}

func calculateAverageAge(records []DataRecord) float64 {
    if len(records) == 0 {
        return 0
    }

    total := 0
    for _, record := range records {
        total += record.Age
    }
    return float64(total) / float64(len(records))
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_cleaner <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, err := parseCSVFile(filename)
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Original records: %d\n", len(records))

    uniqueRecords := removeDuplicates(records)
    fmt.Printf("After deduplication: %d\n", len(uniqueRecords))

    activeUsers := filterActiveUsers(uniqueRecords)
    fmt.Printf("Active users: %d\n", len(activeUsers))

    avgAge := calculateAverageAge(activeUsers)
    fmt.Printf("Average age of active users: %.2f\n", avgAge)

    for i, user := range activeUsers {
        if i < 3 {
            fmt.Printf("User %d: %s (%s) - Age: %d\n", user.ID, user.Name, user.Email, user.Age)
        }
    }
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedRecord[i] = strings.TrimSpace(field)
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("error writing CSV: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully cleaned %s -> %s\n", inputFile, outputFile)
}
package main

import "fmt"

func removeDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}