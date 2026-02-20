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

	headerProcessed := false
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		if !headerProcessed {
			headerProcessed = true
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write header: %w", err)
			}
			continue
		}

		cleaned := make([]string, len(record))
		for i, field := range record {
			cleaned[i] = strings.TrimSpace(field)
			if cleaned[i] == "" {
				cleaned[i] = "N/A"
			}
		}
		if err := writer.Write(cleaned); err != nil {
			return fmt.Errorf("failed to write cleaned record: %w", err)
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
	fmt.Printf("Successfully cleaned data. Output saved to %s\n", outputFile)
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

func DeduplicateEmails(emails []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, email := range emails {
        if !seen[email] {
            seen[email] = true
            result = append(result, email)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
    if len(email) < 3 || !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
        return false
    }
    return true
}

func CleanData(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        if ValidateEmail(record.Email) && !emailSet[record.Email] {
            emailSet[record.Email] = true
            record.Valid = true
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", false},
        {2, "invalid-email", false},
        {3, "user@example.com", false},
        {4, "test@domain.org", false},
    }
    
    cleaned := CleanData(sampleData)
    fmt.Printf("Cleaned records: %d\n", len(cleaned))
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s\n", r.ID, r.Email)
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
}package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    ID      int
    Name    string
    Email   string
    Active  bool
    Score   float64
}

func cleanEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}

func validateRecord(rec Record) error {
    if rec.ID <= 0 {
        return fmt.Errorf("invalid ID: %d", rec.ID)
    }
    if len(rec.Name) == 0 {
        return fmt.Errorf("empty name")
    }
    if !strings.Contains(rec.Email, "@") {
        return fmt.Errorf("invalid email: %s", rec.Email)
    }
    if rec.Score < 0 || rec.Score > 100 {
        return fmt.Errorf("score out of range: %.2f", rec.Score)
    }
    return nil
}

func processCSVFile(inputPath string) ([]Record, []error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, []error{err}
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true
    records := []Record{}
    errors := []error{}
    lineNum := 0

    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: read error: %v", lineNum, err))
            continue
        }
        if len(row) != 5 {
            errors = append(errors, fmt.Errorf("line %d: expected 5 columns, got %d", lineNum, len(row)))
            continue
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: invalid ID: %v", lineNum, err))
            continue
        }

        active, err := strconv.ParseBool(row[3])
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: invalid active flag: %v", lineNum, err))
            continue
        }

        score, err := strconv.ParseFloat(row[4], 64)
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: invalid score: %v", lineNum, err))
            continue
        }

        record := Record{
            ID:     id,
            Name:   strings.TrimSpace(row[1]),
            Email:  cleanEmail(row[2]),
            Active: active,
            Score:  score,
        }

        if valErr := validateRecord(record); valErr != nil {
            errors = append(errors, fmt.Errorf("line %d: validation failed: %v", lineNum, valErr))
            continue
        }

        records = append(records, record)
    }

    return records, errors
}

func writeCleanData(records []Record, outputPath string) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    header := []string{"ID", "Name", "Email", "Active", "Score"}
    if err := writer.Write(header); err != nil {
        return err
    }

    for _, rec := range records {
        row := []string{
            strconv.Itoa(rec.ID),
            rec.Name,
            rec.Email,
            strconv.FormatBool(rec.Active),
            fmt.Sprintf("%.2f", rec.Score),
        }
        if err := writer.Write(row); err != nil {
            return err
        }
    }

    return nil
}

func main() {
    inputFile := "raw_data.csv"
    outputFile := "cleaned_data.csv"

    records, errs := processCSVFile(inputFile)
    if len(errs) > 0 {
        fmt.Printf("Encountered %d errors during processing:\n", len(errs))
        for _, e := range errs {
            fmt.Printf("  - %v\n", e)
        }
    }

    fmt.Printf("Successfully processed %d valid records\n", len(records))

    if len(records) > 0 {
        if err := writeCleanData(records, outputFile); err != nil {
            fmt.Printf("Error writing output file: %v\n", err)
            return
        }
        fmt.Printf("Cleaned data written to %s\n", outputFile)
    }
}