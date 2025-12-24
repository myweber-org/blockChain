package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func processCSVFile(inputPath, outputPath string) error {
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
			return fmt.Errorf("error writing cleaned record: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := processCSVFile(inputFile, outputFile); err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %s to %s\n", inputFile, outputFile)
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
	Tags      []string `json:"tags"`
}

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", profile.ID)
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return fmt.Errorf("invalid email format: %s", profile.Email)
	}

	if profile.Age < 0 || profile.Age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}

	return nil
}

func TransformUserProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(transformed.Username)
	transformed.Email = strings.ToLower(transformed.Email)
	transformed.Tags = append(transformed.Tags, "processed")

	if transformed.Age < 18 {
		transformed.Active = false
	}

	return transformed
}

func ProcessUserProfile(data []byte) ([]byte, error) {
	var profile UserProfile
	err := json.Unmarshal(data, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user profile: %w", err)
	}

	err = ValidateUserProfile(profile)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	transformedProfile := TransformUserProfile(profile)
	result, err := json.MarshalIndent(transformedProfile, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed profile: %w", err)
	}

	return result, nil
}

func main() {
	jsonData := []byte(`{
		"id": 123,
		"username": "JohnDoe",
		"email": "JOHN@EXAMPLE.COM",
		"age": 25,
		"active": true,
		"tags": ["golang", "backend"]
	}`)

	processed, err := ProcessUserProfile(jsonData)
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}

	fmt.Println("Processed user profile:")
	fmt.Println(string(processed))
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
}

func NewDataProcessor(input, output string) *DataProcessor {
	return &DataProcessor{
		InputPath:  input,
		OutputPath: output,
	}
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
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	cleanedCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		recordCount++
		cleanedRecord := dp.cleanRecord(record)
		if dp.isValidRecord(cleanedRecord) {
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			cleanedCount++
		}
	}

	fmt.Printf("Processed %d records, wrote %d valid records\n", recordCount, cleanedCount)
	return nil
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		cleaned[i] = strings.TrimSpace(field)
	}
	return cleaned
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
	for _, field := range record {
		if field == "" {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1], os.Args[2])
	if err := processor.Process(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
package main

import "fmt"

func MovingAverage(data []float64, windowSize int) []float64 {
    if windowSize <= 0 || len(data) == 0 {
        return []float64{}
    }

    result := make([]float64, 0, len(data)-windowSize+1)
    for i := 0; i <= len(data)-windowSize; i++ {
        sum := 0.0
        for j := i; j < i+windowSize; j++ {
            sum += data[j]
        }
        average := sum / float64(windowSize)
        result = append(result, average)
    }
    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    window := 3
    averages := MovingAverage(sampleData, window)
    fmt.Printf("Moving averages (window=%d): %v\n", window, averages)
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TransformString(input string) string {
	return strings.ToUpper(strings.TrimSpace(input))
}

func ConvertToJSON(data map[string]interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func main() {
	email := "test@example.com"
	fmt.Printf("Email validation result: %v\n", ValidateEmail(email))

	str := "  hello world  "
	fmt.Printf("Transformed string: '%s'\n", TransformString(str))

	sampleData := map[string]interface{}{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}
	jsonStr, err := ConvertToJSON(sampleData)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
	} else {
		fmt.Printf("JSON output:\n%s\n", jsonStr)
	}
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

type DataRecord struct {
    ID      int
    Name    string
    Value   float64
    Active  bool
}

func ProcessCSVFile(inputPath string, outputPath string) error {
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

    csvReader := csv.NewReader(inputFile)
    csvWriter := csv.NewWriter(outputFile)
    defer csvWriter.Flush()

    header, err := csvReader.Read()
    if err != nil {
        return fmt.Errorf("failed to read CSV header: %w", err)
    }

    if len(header) != 4 {
        return fmt.Errorf("invalid CSV format: expected 4 columns, got %d", len(header))
    }

    if err := csvWriter.Write([]string{"ID", "Name", "Value", "Active", "Status"}); err != nil {
        return fmt.Errorf("failed to write output header: %w", err)
    }

    lineNumber := 1
    for {
        record, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading line %d: %w", lineNumber, err)
        }

        if len(record) != 4 {
            return fmt.Errorf("invalid data at line %d: expected 4 fields, got %d", lineNumber, len(record))
        }

        dataRec, err := parseRecord(record, lineNumber)
        if err != nil {
            return err
        }

        status := "VALID"
        if dataRec.Value < 0 || dataRec.Value > 1000 {
            status = "VALUE_OUT_OF_RANGE"
        }
        if !dataRec.Active && dataRec.Value > 0 {
            status = "INACTIVE_WITH_POSITIVE_VALUE"
        }

        outputRecord := []string{
            strconv.Itoa(dataRec.ID),
            strings.ToUpper(dataRec.Name),
            fmt.Sprintf("%.2f", dataRec.Value),
            strconv.FormatBool(dataRec.Active),
            status,
        }

        if err := csvWriter.Write(outputRecord); err != nil {
            return fmt.Errorf("failed to write line %d: %w", lineNumber, err)
        }

        lineNumber++
    }

    return nil
}

func parseRecord(fields []string, lineNumber int) (DataRecord, error) {
    var rec DataRecord
    var err error

    rec.ID, err = strconv.Atoi(fields[0])
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
    }

    rec.Name = strings.TrimSpace(fields[1])
    if rec.Name == "" {
        return DataRecord{}, fmt.Errorf("empty name at line %d", lineNumber)
    }

    rec.Value, err = strconv.ParseFloat(fields[2], 64)
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
    }

    rec.Active, err = strconv.ParseBool(fields[3])
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
    }

    return rec, nil
}

func ValidateCSVStructure(filePath string) (int, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return 0, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return 0, fmt.Errorf("failed to read CSV: %w", err)
    }

    if len(records) == 0 {
        return 0, fmt.Errorf("empty CSV file")
    }

    expectedColumns := 4
    if len(records[0]) != expectedColumns {
        return 0, fmt.Errorf("invalid header: expected %d columns, got %d", expectedColumns, len(records[0]))
    }

    for i := 1; i < len(records); i++ {
        if len(records[i]) != expectedColumns {
            return 0, fmt.Errorf("line %d: expected %d columns, got %d", i+1, expectedColumns, len(records[i]))
        }
    }

    return len(records) - 1, nil
}