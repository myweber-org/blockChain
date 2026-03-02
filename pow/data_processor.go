package main

import "fmt"

func calculateAverage(numbers []float64) float64 {
    if len(numbers) == 0 {
        return 0
    }
    
    var sum float64
    for _, num := range numbers {
        sum += num
    }
    
    return sum / float64(len(numbers))
}

func main() {
    data := []float64{10.5, 20.3, 30.7, 40.1, 50.9}
    avg := calculateAverage(data)
    fmt.Printf("Average: %.2f\n", avg)
}package main

import (
	"errors"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	for _, r := range data.Username {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' {
			return errors.New("username can only contain letters, numbers, and underscores")
		}
	}

	if !strings.Contains(data.Email, "@") || !strings.Contains(data.Email, ".") {
		return errors.New("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ProcessUserInput(rawUsername, rawEmail string, age int) (UserData, error) {
	userData := UserData{
		Username: TransformUsername(rawUsername),
		Email:    NormalizeEmail(rawEmail),
		Age:      age,
	}

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}
package main

import "fmt"

func movingAverage(data []float64, windowSize int) []float64 {
    if windowSize <= 0 || len(data) == 0 {
        return []float64{}
    }

    result := make([]float64, 0, len(data)-windowSize+1)
    for i := 0; i <= len(data)-windowSize; i++ {
        sum := 0.0
        for j := i; j < i+windowSize; j++ {
            sum += data[j]
        }
        result = append(result, sum/float64(windowSize))
    }
    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    averaged := movingAverage(sampleData, 3)
    fmt.Printf("Moving average result: %v\n", averaged)
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateAndParseJSON(rawData []byte) (*UserData, error) {
	var user UserData
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if user.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	return &user, nil
}

func main() {
	jsonStr := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`
	user, err := ValidateAndParseJSON([]byte(jsonStr))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Parsed user: %+v\n", user)
}
package main

import (
	"fmt"
)

// FilterAndDouble filters out even numbers from the input slice,
// doubles the remaining odd numbers, and returns the new slice.
func FilterAndDouble(numbers []int) []int {
	var result []int
	for _, num := range numbers {
		if num%2 != 0 {
			result = append(result, num*2)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := FilterAndDouble(input)
	fmt.Printf("Input: %v\n", input)
	fmt.Printf("Output: %v\n", output)
}
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

func (dp *DataProcessor) ValidateAndClean() error {
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

    lineNumber := 0
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV at line %d: %w", lineNumber, err)
        }

        cleanedRecord := dp.cleanRecord(record)
        if len(cleanedRecord) > 0 {
            if err := writer.Write(cleanedRecord); err != nil {
                return fmt.Errorf("error writing record at line %d: %w", lineNumber, err)
            }
        }

        lineNumber++
    }

    writer.Flush()
    if err := writer.Error(); err != nil {
        return fmt.Errorf("error flushing writer: %w", err)
    }

    return nil
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
    cleaned := make([]string, 0, len(record))
    for _, field := range record {
        cleanedField := strings.TrimSpace(field)
        cleanedField = strings.ToValidUTF8(cleanedField, "")
        if cleanedField != "" {
            cleaned = append(cleaned, cleanedField)
        }
    }
    return cleaned
}

func (dp *DataProcessor) SetDelimiter(delim rune) {
    dp.Delimiter = delim
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: data_processor <input.csv> <output.csv>")
        os.Exit(1)
    }

    processor := NewDataProcessor(os.Args[1], os.Args[2])
    if err := processor.ValidateAndClean(); err != nil {
        fmt.Printf("Processing failed: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data processing completed successfully")
}