
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func processCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) != 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		record := Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		}
		records = append(records, record)
	}

	return records, nil
}

func calculateTotal(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <filename.csv>")
		return
	}

	records, err := processCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Total value: %.2f\n", calculateTotal(records))
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ParseJSONToMap parses JSON data into a map[string]interface{}.
func ParseJSONToMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := `{"name": "test", "value": 42, "active": true}`

	valid, err := ValidateJSON([]byte(sampleJSON))
	if err != nil {
		log.Printf("Validation error: %v", err)
	} else {
		fmt.Println("JSON is valid:", valid)
	}

	parsed, err := ParseJSONToMap([]byte(sampleJSON))
	if err != nil {
		log.Printf("Parse error: %v", err)
	} else {
		fmt.Printf("Parsed data: %+v\n", parsed)
	}
}