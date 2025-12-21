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