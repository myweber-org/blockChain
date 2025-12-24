package main

import "fmt"

func RemoveDuplicates(input []int) []int {
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
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
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
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    ID    int
    Name  string
    Email string
    Score float64
}

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

    header, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }

    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    lineNum := 1
    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("line %d: read error: %w", lineNum, err)
        }

        if len(row) != 4 {
            continue
        }

        id, err := strconv.Atoi(strings.TrimSpace(row[0]))
        if err != nil || id <= 0 {
            continue
        }

        name := strings.TrimSpace(row[1])
        if name == "" {
            continue
        }

        email := strings.TrimSpace(row[2])
        if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
            continue
        }

        score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
        if err != nil || score < 0 || score > 100 {
            continue
        }

        cleanRow := []string{
            strconv.Itoa(id),
            name,
            strings.ToLower(email),
            fmt.Sprintf("%.2f", score),
        }

        if err := writer.Write(cleanRow); err != nil {
            return fmt.Errorf("line %d: write error: %w", lineNum, err)
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

    fmt.Println("Data cleaning completed successfully")
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func CleanCSVData(reader io.Reader, writer io.Writer) error {
	csvReader := csv.NewReader(reader)
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedField := strings.TrimSpace(field)
			cleanedField = strings.ToLower(cleanedField)
			cleanedRecord[i] = cleanedField
		}

		if err := csvWriter.Write(cleanedRecord); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}
	return nil
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

func cleanString(s string) string {
    return strings.TrimSpace(strings.ToLower(s))
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processCSVFile(inputPath string) ([]DataRecord, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []DataRecord{}
    lineNum := 0

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        lineNum++
        if lineNum == 1 {
            continue
        }

        if len(row) < 3 {
            continue
        }

        record := DataRecord{
            ID:    cleanString(row[0]),
            Name:  cleanString(row[1]),
            Email: cleanString(row[2]),
        }
        record.Valid = validateEmail(record.Email)

        records = append(records, record)
    }

    return records, nil
}

func writeCleanData(records []DataRecord, outputPath string) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    header := []string{"id", "name", "email", "valid"}
    if err := writer.Write(header); err != nil {
        return err
    }

    for _, record := range records {
        row := []string{
            record.ID,
            record.Name,
            record.Email,
            fmt.Sprintf("%t", record.Valid),
        }
        if err := writer.Write(row); err != nil {
            return err
        }
    }

    return nil
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    records, err := processCSVFile(inputFile)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    validCount := 0
    for _, record := range records {
        if record.Valid {
            validCount++
        }
    }

    err = writeCleanData(records, outputFile)
    if err != nil {
        fmt.Printf("Error writing output: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Processed %d records, %d valid entries\n", len(records), validCount)
    fmt.Printf("Clean data written to %s\n", outputFile)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !dc.seen[trimmed] {
			dc.seen[trimmed] = true
			unique = append(unique, trimmed)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"  apple  ", "banana", "apple", "", "orange", "banana"}
	unique := cleaner.RemoveDuplicates(data)
	fmt.Printf("Unique items: %v\n", unique)
	
	emails := []string{"test@example.com", "invalid-email", "user@domain"}
	for _, email := range emails {
		valid := cleaner.ValidateEmail(email)
		fmt.Printf("Email %s valid: %v\n", email, valid)
	}
	
	cleaner.Reset()
	fmt.Println("Cleaner reset completed")
}package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	
	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	
	for _, record := range records {
		if ValidateEmail(record.Email) {
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	
	return DeduplicateRecords(cleaned)
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob Wilson", "invalid-email", false},
		{5, "Alice Brown", "alice@company.org", false},
	}
	
	cleaned := CleanData(sampleData)
	
	fmt.Printf("Original records: %d\n", len(sampleData))
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Valid: %v\n", 
			record.ID, record.Name, record.Email, record.Valid)
	}
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
	ID    string
	Email string
	Phone string
}

type DataCleaner struct {
	records map[string]DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make(map[string]DataRecord),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) bool {
	if !dc.isValidRecord(record) {
		return false
	}

	key := dc.generateKey(record)
	if _, exists := dc.records[key]; exists {
		return false
	}

	dc.records[key] = record
	return true
}

func (dc *DataCleaner) isValidRecord(record DataRecord) bool {
	if record.ID == "" || record.Email == "" {
		return false
	}
	if !strings.Contains(record.Email, "@") {
		return false
	}
	return true
}

func (dc *DataCleaner) generateKey(record DataRecord) string {
	return fmt.Sprintf("%s|%s", record.ID, strings.ToLower(record.Email))
}

func (dc *DataCleaner) LoadFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Email: strings.TrimSpace(row[1]),
			Phone: strings.TrimSpace(row[2]),
		}

		dc.AddRecord(record)
	}

	return nil
}

func (dc *DataCleaner) GetUniqueRecords() []DataRecord {
	var result []DataRecord
	for _, record := range dc.records {
		result = append(result, record)
	}
	return result
}

func (dc *DataCleaner) Count() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	sampleData := []DataRecord{
		{"001", "user@example.com", "1234567890"},
		{"002", "admin@test.org", "0987654321"},
		{"001", "user@example.com", "5555555555"},
		{"003", "invalid-email", "1111111111"},
		{"004", "another@domain.com", "2222222222"},
	}

	for _, record := range sampleData {
		if added := cleaner.AddRecord(record); added {
			fmt.Printf("Added: %+v\n", record)
		} else {
			fmt.Printf("Skipped duplicate/invalid: %+v\n", record)
		}
	}

	fmt.Printf("Total unique records: %d\n", cleaner.Count())
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processRecords(records []DataRecord) []DataRecord {
	emailMap := make(map[string]bool)
	var validRecords []DataRecord

	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		if validateEmail(cleanEmail) && !emailMap[cleanEmail] {
			emailMap[cleanEmail] = true
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
	}

	cleaned := processRecords(records)
	fmt.Printf("Processed %d records, %d valid\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}package main

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
	"fmt"
	"strings"
)

func CleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	data := []string{" apple ", "banana", " apple", "banana ", "", "  cherry  "}
	cleaned := CleanData(data)
	fmt.Println("Cleaned data:", cleaned)
}
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	sample := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	unique := RemoveDuplicates(sample)
	fmt.Println("Original:", sample)
	fmt.Println("Unique:", unique)
}