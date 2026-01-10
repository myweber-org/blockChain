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
    lineNumber := 0

    for {
        row, err := reader.Read()
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

func generateReport(records []DataRecord) {
    validCount := 0
    for _, record := range records {
        if record.Valid {
            validCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Valid records: %d\n", validCount)
    fmt.Printf("Invalid records: %d\n", len(records)-validCount)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_cleaner <csv_file>")
        return
    }

    records, err := processCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    generateReport(records)
}package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateEmails(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, DataRecord{
				ID:    record.ID,
				Email: email,
				Valid: validateEmail(email),
			})
		}
	}
	return unique
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func PrintRecords(records []DataRecord) {
	for _, r := range records {
		status := "INVALID"
		if r.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", r.ID, r.Email, status)
	}
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
		{5, "user@example.com", false},
		{6, "", false},
	}

	cleaned := DeduplicateEmails(sampleData)
	fmt.Println("Cleaned Data Records:")
	PrintRecords(cleaned)
}