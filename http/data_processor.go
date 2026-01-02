
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID      string
	Name    string
	Email   string
	Active  string
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
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

		if lineNumber == 1 {
			continue
		}

		if len(row) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if record.ID == "" || record.Name == "" {
			return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
		}

		if !strings.Contains(record.Email, "@") {
			return nil, fmt.Errorf("invalid email format at line %d", lineNumber)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) []string {
	var errors []string
	emailSet := make(map[string]bool)

	for i, record := range records {
		if record.Active != "true" && record.Active != "false" {
			errors = append(errors, fmt.Sprintf("record %d: invalid active status '%s'", i+1, record.Active))
		}

		if emailSet[record.Email] {
			errors = append(errors, fmt.Sprintf("record %d: duplicate email '%s'", i+1, record.Email))
		}
		emailSet[record.Email] = true
	}

	return errors
}

func GenerateReport(records []DataRecord) {
	activeCount := 0
	for _, record := range records {
		if record.Active == "true" {
			activeCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Active records: %d\n", activeCount)
	fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
}
package main

import (
    "encoding/csv"
    "errors"
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

type DataProcessor struct {
    records []DataRecord
}

func NewDataProcessor() *DataProcessor {
    return &DataProcessor{
        records: make([]DataRecord, 0),
    }
}

func (dp *DataProcessor) LoadFromCSV(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    lineNumber := 0
    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if lineNumber == 1 {
            continue
        }

        record, err := parseCSVRow(row, lineNumber)
        if err != nil {
            return err
        }

        dp.records = append(dp.records, record)
    }

    return nil
}

func parseCSVRow(row []string, lineNumber int) (DataRecord, error) {
    if len(row) != 4 {
        return DataRecord{}, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
    }

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
    }

    name := strings.TrimSpace(row[1])
    if name == "" {
        return DataRecord{}, fmt.Errorf("empty name at line %d", lineNumber)
    }

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
    }

    active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
    }

    return DataRecord{
        ID:     id,
        Name:   name,
        Value:  value,
        Active: active,
    }, nil
}

func (dp *DataProcessor) FilterActive() []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range dp.records {
        if record.Active {
            filtered = append(filtered, record)
        }
    }
    return filtered
}

func (dp *DataProcessor) CalculateTotal() float64 {
    total := 0.0
    for _, record := range dp.records {
        total += record.Value
    }
    return total
}

func (dp *DataProcessor) FindByName(name string) (DataRecord, error) {
    for _, record := range dp.records {
        if strings.EqualFold(record.Name, name) {
            return record, nil
        }
    }
    return DataRecord{}, errors.New("record not found")
}

func (dp *DataProcessor) Count() int {
    return len(dp.records)
}

func (dp *DataProcessor) ExportToCSV(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    header := []string{"ID", "Name", "Value", "Active"}
    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    for _, record := range dp.records {
        row := []string{
            strconv.Itoa(record.ID),
            record.Name,
            strconv.FormatFloat(record.Value, 'f', 2, 64),
            strconv.FormatBool(record.Active),
        }
        if err := writer.Write(row); err != nil {
            return fmt.Errorf("failed to write row: %w", err)
        }
    }

    return nil
}