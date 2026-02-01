package main

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

        cleaned := make([]string, len(record))
        for i, field := range record {
            cleaned[i] = strings.TrimSpace(field)
        }

        if err := writer.Write(cleaned); err != nil {
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

    err := cleanCSV(os.Args[1], os.Args[2])
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("CSV cleaning completed successfully")
}
package main

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

func deduplicateRecords(records []DataRecord) []DataRecord {
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

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func validateRecords(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		record.Valid = validateEmail(record.Email)
		validated = append(validated, record)
	}
	return validated
}

func processData(records []DataRecord) []DataRecord {
	unique := deduplicateRecords(records)
	validated := validateRecords(unique)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob Wilson", "bob.wilson", false},
	}

	processed := processData(sampleData)

	for _, record := range processed {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, Status: %s\n",
			record.ID, record.Name, record.Email, status)
	}
}