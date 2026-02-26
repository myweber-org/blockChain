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
    Age     int
    Active  bool
}

func cleanCSV(inputPath, outputPath string) error {
    inFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("cannot open input file: %w", err)
    }
    defer inFile.Close()

    outFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("cannot create output file: %w", err)
    }
    defer outFile.Close()

    reader := csv.NewReader(inFile)
    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    headers, err := reader.Read()
    if err != nil {
        return fmt.Errorf("cannot read headers: %w", err)
    }

    if err := writer.Write(headers); err != nil {
        return fmt.Errorf("cannot write headers: %w", err)
    }

    lineNum := 1
    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("line %d: skip due to read error: %v\n", lineNum, err)
            continue
        }

        if len(row) != 5 {
            fmt.Printf("line %d: skip due to column count mismatch\n", lineNum)
            continue
        }

        record, err := parseRecord(row)
        if err != nil {
            fmt.Printf("line %d: skip due to parse error: %v\n", lineNum, err)
            continue
        }

        if !validateRecord(record) {
            fmt.Printf("line %d: skip due to validation failure\n", lineNum)
            continue
        }

        cleaned := []string{
            strconv.Itoa(record.ID),
            strings.TrimSpace(record.Name),
            strings.ToLower(strings.TrimSpace(record.Email)),
            strconv.Itoa(record.Age),
            strconv.FormatBool(record.Active),
        }

        if err := writer.Write(cleaned); err != nil {
            return fmt.Errorf("line %d: cannot write row: %w", lineNum, err)
        }
    }

    return nil
}

func parseRecord(row []string) (Record, error) {
    var r Record
    var err error

    r.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return r, fmt.Errorf("invalid ID: %w", err)
    }

    r.Name = row[1]
    r.Email = row[2]

    r.Age, err = strconv.Atoi(strings.TrimSpace(row[3]))
    if err != nil {
        return r, fmt.Errorf("invalid age: %w", err)
    }

    r.Active, err = strconv.ParseBool(strings.TrimSpace(row[4]))
    if err != nil {
        return r, fmt.Errorf("invalid active flag: %w", err)
    }

    return r, nil
}

func validateRecord(r Record) bool {
    if r.ID <= 0 {
        return false
    }
    if len(r.Name) == 0 || len(r.Name) > 100 {
        return false
    }
    if !strings.Contains(r.Email, "@") {
        return false
    }
    if r.Age < 0 || r.Age > 150 {
        return false
    }
    return true
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

import "fmt"

func RemoveDuplicates(nums []int) []int {
    if len(nums) == 0 {
        return nums
    }

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
    data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
    cleaned := RemoveDuplicates(data)
    fmt.Printf("Original: %v\n", data)
    fmt.Printf("Cleaned: %v\n", cleaned)
}