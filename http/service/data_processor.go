
package main

import (
	"fmt"
)

// FilterAndDouble filters out even numbers from a slice and doubles the remaining odd numbers.
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
	fmt.Println("Original:", input)
	fmt.Println("Filtered and Doubled:", output)
}
package main

import (
    "encoding/csv"
    "errors"
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

func ProcessCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]Record, 0)
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
        }

        name := row[1]
        if name == "" {
            return nil, fmt.Errorf("empty name at line %d", lineNumber)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func CalculateTotalValue(records []Record) float64 {
    total := 0.0
    for _, record := range records {
        total += record.Value
    }
    return total
}

func FindMaxValueRecord(records []Record) *Record {
    if len(records) == 0 {
        return nil
    }

    maxRecord := &records[0]
    for i := 1; i < len(records); i++ {
        if records[i].Value > maxRecord.Value {
            maxRecord = &records[i]
        }
    }
    return maxRecord
}