
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

func RemoveDuplicates(records []DataRecord) []DataRecord {
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

func ValidateRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		record.Valid = strings.Contains(record.Email, "@") && len(record.Name) > 0
		if record.Valid {
			valid = append(valid, record)
		}
	}
	return valid
}

func CleanData(records []DataRecord) []DataRecord {
	unique := RemoveDuplicates(records)
	valid := ValidateRecords(unique)
	return valid
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Invalid User", "invalid-email", false},
		{5, "", "empty@example.com", false},
	}

	cleaned := CleanData(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
	
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func removeDuplicates(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	seen := make(map[string]bool)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}
		if len(record) == 0 {
			continue
		}
		key := record[0]
		if !seen[key] {
			seen[key] = true
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	outputPath := os.Args[2]
	if err := removeDuplicates(inputPath, outputPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Duplicate removal completed successfully")
}
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

    seen := make(map[string]bool)
    headerWritten := false

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV: %w", err)
        }

        for i, field := range record {
            record[i] = strings.TrimSpace(field)
        }

        key := strings.Join(record, "|")
        if seen[key] {
            continue
        }
        seen[key] = true

        if !headerWritten {
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("error writing header: %w", err)
            }
            headerWritten = true
            continue
        }

        if err := writer.Write(record); err != nil {
            return fmt.Errorf("error writing record: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run data_cleaner.go <input.csv> <output.csv>")
        os.Exit(1)
    }

    err := cleanCSV(os.Args[1], os.Args[2])
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("CSV cleaning completed successfully")
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
	numbers := []int{1, 2, 2, 3, 4, 4, 5, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange", "banana"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}package main

import (
	"fmt"
	"sort"
)

type DataRecord struct {
	ID   int
	Name string
}

func cleanData(records []DataRecord) []DataRecord {
	seen := make(map[int]bool)
	var unique []DataRecord

	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			unique = append(unique, record)
		}
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].ID < unique[j].ID
	})

	return unique
}

func main() {
	records := []DataRecord{
		{ID: 3, Name: "Charlie"},
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 1, Name: "Alice"},
		{ID: 4, Name: "David"},
	}

	cleaned := cleanData(records)
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
	}
}
package main

import (
	"fmt"
)

// RemoveDuplicates removes duplicate strings from a slice.
// It preserves the order of the first occurrence of each element.
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
	data := []string{"apple", "banana", "apple", "cherry", "banana", "date"}
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

	seen := make(map[string]bool)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV record: %w", err)
		}

		for i, field := range record {
			record[i] = strings.TrimSpace(field)
		}

		key := strings.Join(record, "|")
		if seen[key] {
			continue
		}
		seen[key] = true

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %w", err)
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

	fmt.Printf("Successfully cleaned data written to %s\n", outputFile)
}