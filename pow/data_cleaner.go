package datautils

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
package main

import "fmt"

func removeDuplicates(nums []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    
    for _, num := range nums {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    return result
}

func main() {
    data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
    cleaned := removeDuplicates(data)
    fmt.Printf("Original: %v\n", data)
    fmt.Printf("Cleaned: %v\n", cleaned)
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
    ID      int
    Name    string
    Email   string
    Active  bool
    Score   float64
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

    headers, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read headers: %w", err)
    }

    if err := writer.Write(headers); err != nil {
        return fmt.Errorf("failed to write headers: %w", err)
    }

    lineNum := 1
    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("Warning: line %d: %v\n", lineNum, err)
            continue
        }

        if len(row) != 5 {
            fmt.Printf("Warning: line %d: invalid column count\n", lineNum)
            continue
        }

        record, err := parseRecord(row)
        if err != nil {
            fmt.Printf("Warning: line %d: %v\n", lineNum, err)
            continue
        }

        cleaned := cleanRecord(record)
        outputRow := []string{
            strconv.Itoa(cleaned.ID),
            cleaned.Name,
            cleaned.Email,
            strconv.FormatBool(cleaned.Active),
            fmt.Sprintf("%.2f", cleaned.Score),
        }

        if err := writer.Write(outputRow); err != nil {
            return fmt.Errorf("failed to write row: %w", err)
        }
    }

    return nil
}

func parseRecord(row []string) (Record, error) {
    var rec Record
    var err error

    if rec.ID, err = strconv.Atoi(row[0]); err != nil {
        return rec, fmt.Errorf("invalid ID: %w", err)
    }

    rec.Name = strings.TrimSpace(row[1])
    if rec.Name == "" {
        return rec, fmt.Errorf("empty name")
    }

    rec.Email = strings.TrimSpace(row[2])
    if !strings.Contains(rec.Email, "@") {
        return rec, fmt.Errorf("invalid email format")
    }

    if rec.Active, err = strconv.ParseBool(row[3]); err != nil {
        return rec, fmt.Errorf("invalid active flag: %w", err)
    }

    if rec.Score, err = strconv.ParseFloat(row[4], 64); err != nil {
        return rec, fmt.Errorf("invalid score: %w", err)
    }

    return rec, nil
}

func cleanRecord(rec Record) Record {
    rec.Name = strings.Title(strings.ToLower(rec.Name))
    rec.Email = strings.ToLower(rec.Email)
    if rec.Score < 0 {
        rec.Score = 0
    } else if rec.Score > 100 {
        rec.Score = 100
    }
    return rec
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
}