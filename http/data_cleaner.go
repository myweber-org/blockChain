package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type Cleaner struct {
    inputPath  string
    outputPath string
    seen       map[string]bool
}

func NewCleaner(input, output string) *Cleaner {
    return &Cleaner{
        inputPath:  input,
        outputPath: output,
        seen:       make(map[string]bool),
    }
}

func (c *Cleaner) Process() error {
    inFile, err := os.Open(c.inputPath)
    if err != nil {
        return fmt.Errorf("open input file: %w", err)
    }
    defer inFile.Close()

    outFile, err := os.Create(c.outputPath)
    if err != nil {
        return fmt.Errorf("create output file: %w", err)
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
            return fmt.Errorf("read csv record: %w", err)
        }

        cleaned := c.cleanRecord(record)
        if cleaned == nil {
            continue
        }

        if err := writer.Write(cleaned); err != nil {
            return fmt.Errorf("write csv record: %w", err)
        }
    }

    return nil
}

func (c *Cleaner) cleanRecord(record []string) []string {
    var cleaned []string
    for _, field := range record {
        trimmed := strings.TrimSpace(field)
        cleaned = append(cleaned, trimmed)
    }

    key := strings.Join(cleaned, "|")
    if c.seen[key] {
        return nil
    }
    c.seen[key] = true

    return cleaned
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    cleaner := NewCleaner(os.Args[1], os.Args[2])
    if err := cleaner.Process(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
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
	var valid []DataRecord
	for _, record := range records {
		record.Valid = validateEmail(record.Email)
		if record.Valid {
			valid = append(valid, record)
		}
	}
	return valid
}

func cleanData(records []DataRecord) []DataRecord {
	deduped := deduplicateRecords(records)
	validated := validateRecords(deduped)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob Wilson", "invalid-email", false},
	}

	cleaned := cleanData(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
}