
package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Value float64 `json:"value"`
}

func processCSVFile(inputPath string) ([]Record, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []Record
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error on line %d: %w", lineNumber, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count on line %d: expected 3, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID on line %d: %w", lineNumber, err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value on line %d: %w", lineNumber, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    return records, nil
}

func generateJSONOutput(records []Record, outputPath string) error {
    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    encoder := json.NewEncoder(outputFile)
    encoder.SetIndent("", "  ")

    if err := encoder.Encode(records); err != nil {
        return fmt.Errorf("failed to encode JSON: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_processor <input.csv> <output.json>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    records, err := processCSVFile(inputFile)
    if err != nil {
        fmt.Printf("Error processing CSV: %v\n", err)
        os.Exit(1)
    }

    if err := generateJSONOutput(records, outputFile); err != nil {
        fmt.Printf("Error generating JSON: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d records to %s\n", len(records), outputFile)
}package data_processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type ValidationRule func(interface{}) error

type JSONProcessor struct {
	rules map[string]ValidationRule
}

func NewJSONProcessor() *JSONProcessor {
	return &JSONProcessor{
		rules: make(map[string]ValidationRule),
	}
}

func (jp *JSONProcessor) AddRule(field string, rule ValidationRule) {
	jp.rules[field] = rule
}

func (jp *JSONProcessor) Process(rawData []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}
	
	if len(jp.rules) == 0 {
		return data, nil
	}
	
	var validationErrors []string
	
	for field, rule := range jp.rules {
		if value, exists := data[field]; exists {
			if err := rule(value); err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("%s: %v", field, err))
			}
		}
	}
	
	if len(validationErrors) > 0 {
		return nil, errors.New(strings.Join(validationErrors, "; "))
	}
	
	return data, nil
}

func RequiredStringRule(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("must be a string")
	}
	
	if strings.TrimSpace(str) == "" {
		return errors.New("cannot be empty")
	}
	
	return nil
}

func PositiveNumberRule(value interface{}) error {
	num, ok := value.(float64)
	if !ok {
		return errors.New("must be a number")
	}
	
	if num <= 0 {
		return errors.New("must be positive")
	}
	
	return nil
}