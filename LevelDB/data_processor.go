
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

func processCSVFile(inputPath, outputPath string) error {
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

    headerProcessed := false
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV record: %w", err)
        }

        if !headerProcessed {
            headerProcessed = true
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("error writing header: %w", err)
            }
            continue
        }

        cleanedRecord := make([]string, len(record))
        for i, field := range record {
            cleanedRecord[i] = strings.TrimSpace(field)
            if cleanedRecord[i] == "" {
                cleanedRecord[i] = "N/A"
            }
        }

        if err := writer.Write(cleanedRecord); err != nil {
            return fmt.Errorf("error writing record: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run data_processor.go <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := processCSVFile(inputFile, outputFile); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %s -> %s\n", inputFile, outputFile)
}package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if len(strings.TrimSpace(data.Username)) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	userData := UserData{
		Email:    strings.TrimSpace(email),
		Username: TransformUsername(username),
		Age:      age,
	}
	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}
	return userData, nil
}