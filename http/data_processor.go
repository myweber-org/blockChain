
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataProcessor struct {
    InputPath  string
    OutputPath string
    Delimiter  rune
}

func NewDataProcessor(input, output string) *DataProcessor {
    return &DataProcessor{
        InputPath:  input,
        OutputPath: output,
        Delimiter:  ',',
    }
}

func (dp *DataProcessor) ValidateRow(row []string) bool {
    if len(row) == 0 {
        return false
    }
    for _, field := range row {
        if strings.TrimSpace(field) == "" {
            return false
        }
    }
    return true
}

func (dp *DataProcessor) CleanField(field string) string {
    cleaned := strings.TrimSpace(field)
    cleaned = strings.ToUpper(cleaned)
    return cleaned
}

func (dp *DataProcessor) Process() error {
    inputFile, err := os.Open(dp.InputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(dp.OutputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    reader.Comma = dp.Delimiter

    writer := csv.NewWriter(outputFile)
    writer.Comma = dp.Delimiter
    defer writer.Flush()

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV: %w", err)
        }

        if !dp.ValidateRow(record) {
            continue
        }

        cleanedRecord := make([]string, len(record))
        for i, field := range record {
            cleanedRecord[i] = dp.CleanField(field)
        }

        if err := writer.Write(cleanedRecord); err != nil {
            return fmt.Errorf("error writing CSV: %w", err)
        }
    }

    return nil
}

func main() {
    processor := NewDataProcessor("input.csv", "output.csv")
    if err := processor.Process(); err != nil {
        fmt.Printf("Processing failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Data processing completed successfully")
}package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
}

func ValidateUser(data UserData) error {
	if data.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", data.ID)
	}
	if strings.TrimSpace(data.Name) == "" {
		return fmt.Errorf("user name cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return fmt.Errorf("invalid email format: %s", data.Email)
	}
	return nil
}

func TransformUserName(data *UserData) {
	data.Name = strings.ToUpper(strings.TrimSpace(data.Name))
}

func ProcessUserInput(data UserData) (UserData, error) {
	if err := ValidateUser(data); err != nil {
		return UserData{}, err
	}
	TransformUserName(&data)
	return data, nil
}

func main() {
	user := UserData{
		ID:    101,
		Name:  "  john doe  ",
		Email: "john@example.com",
	}
	
	processed, err := ProcessUserInput(user)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Processed user: %+v\n", processed)
}