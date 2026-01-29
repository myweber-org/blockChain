
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
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Value   int    `json:"value"`
    Active  bool   `json:"active"`
}

func parseCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []Record{}
    lineNumber := 0

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if lineNumber == 0 {
            lineNumber++
            continue
        }

        id, _ := strconv.Atoi(row[0])
        value, _ := strconv.Atoi(row[2])
        active, _ := strconv.ParseBool(row[3])

        record := Record{
            ID:     id,
            Name:   row[1],
            Value:  value,
            Active: active,
        }
        records = append(records, record)
        lineNumber++
    }
    return records, nil
}

func convertToJSON(records []Record) (string, error) {
    jsonData, err := json.MarshalIndent(records, "", "  ")
    if err != nil {
        return "", err
    }
    return string(jsonData), nil
}

func filterActiveRecords(records []Record) []Record {
    var active []Record
    for _, record := range records {
        if record.Active {
            active = append(active, record)
        }
    }
    return active
}

func calculateTotalValue(records []Record) int {
    total := 0
    for _, record := range records {
        total += record.Value
    }
    return total
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        return
    }

    records, err := parseCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error parsing CSV: %v\n", err)
        return
    }

    fmt.Printf("Total records: %d\n", len(records))
    
    activeRecords := filterActiveRecords(records)
    fmt.Printf("Active records: %d\n", len(activeRecords))
    
    totalValue := calculateTotalValue(records)
    fmt.Printf("Total value: %d\n", totalValue)

    jsonOutput, err := convertToJSON(records)
    if err != nil {
        fmt.Printf("Error converting to JSON: %v\n", err)
        return
    }

    outputFile := "output.json"
    err = os.WriteFile(outputFile, []byte(jsonOutput), 0644)
    if err != nil {
        fmt.Printf("Error writing JSON file: %v\n", err)
        return
    }
    
    fmt.Printf("JSON output written to %s\n", outputFile)
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return normalized
}

func (dp *DataProcessor) NormalizeCase(input string, toUpper bool) string {
	cleaned := dp.CleanString(input)
	if toUpper {
		return strings.ToUpper(cleaned)
	}
	return strings.ToLower(cleaned)
}

func (dp *DataProcessor) ExtractAlphanumeric(input string) string {
	alnumRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleaned := dp.CleanString(input)
	return alnumRegex.ReplaceAllString(cleaned, "")
}

func (dp *DataProcessor) ValidateEmail(input string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(strings.TrimSpace(input))
}
package main

import (
	"regexp"
	"strings"
)

// CleanString removes extra whitespace and normalizes input
func CleanString(input string) string {
	// Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	re := regexp.MustCompile(`\s+`)
	normalized := re.ReplaceAllString(trimmed, " ")
	
	return normalized
}

// NormalizeEmail converts email to lowercase and trims spaces
func NormalizeEmail(email string) string {
	cleaned := CleanString(email)
	return strings.ToLower(cleaned)
}

// ValidateUsername checks if username meets requirements
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	
	// Only allow alphanumeric and underscore
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validPattern.MatchString(username)
}