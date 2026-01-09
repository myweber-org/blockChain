package main

import (
    "encoding/csv"
    "errors"
    "io"
    "os"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
    Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if len(row) != 4 {
            continue
        }

        record, err := validateAndCreateRecord(row)
        if err != nil {
            continue
        }

        records = append(records, record)
    }

    return records, nil
}

func validateAndCreateRecord(row []string) (DataRecord, error) {
    var record DataRecord

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, errors.New("invalid id")
    }
    record.ID = id

    name := strings.TrimSpace(row[1])
    if name == "" {
        return record, errors.New("empty name")
    }
    record.Name = name

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return record, errors.New("invalid value")
    }
    record.Value = value

    valid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return record, errors.New("invalid boolean")
    }
    record.Valid = valid

    return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range records {
        if record.Valid {
            filtered = append(filtered, record)
        }
    }
    return filtered
}

func CalculateAverageValue(records []DataRecord) float64 {
    if len(records) == 0 {
        return 0.0
    }

    total := 0.0
    for _, record := range records {
        total += record.Value
    }
    return total / float64(len(records))
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
    inputPath  string
    outputPath string
    delimiter  rune
}

func NewDataProcessor(input, output string) *DataProcessor {
    return &DataProcessor{
        inputPath:  input,
        outputPath: output,
        delimiter:  ',',
    }
}

func (dp *DataProcessor) SetDelimiter(delim rune) {
    dp.delimiter = delim
}

func (dp *DataProcessor) ValidateRow(record []string) bool {
    if len(record) == 0 {
        return false
    }
    for _, field := range record {
        if strings.TrimSpace(field) == "" {
            return false
        }
    }
    return true
}

func (dp *DataProcessor) CleanField(field string) string {
    cleaned := strings.TrimSpace(field)
    cleaned = strings.ToLower(cleaned)
    return cleaned
}

func (dp *DataProcessor) Process() error {
    inputFile, err := os.Open(dp.inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(dp.outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    reader.Comma = dp.delimiter
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    headerWritten := false
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV: %w", err)
        }

        if !headerWritten {
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("error writing header: %w", err)
            }
            headerWritten = true
            continue
        }

        if !dp.ValidateRow(record) {
            continue
        }

        cleanedRecord := make([]string, len(record))
        for i, field := range record {
            cleanedRecord[i] = dp.CleanField(field)
        }

        if err := writer.Write(cleanedRecord); err != nil {
            return fmt.Errorf("error writing record: %w", err)
        }
    }

    return nil
}

func main() {
    processor := NewDataProcessor("input.csv", "output.csv")
    processor.SetDelimiter(',')

    if err := processor.Process(); err != nil {
        fmt.Printf("Processing error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data processing completed successfully")
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

type DataRecord struct {
    ID    int
    Name  string
    Value float64
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    for line := 1; ; line++ {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", line, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", line, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", line, err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", line, err)
        }

        records = append(records, DataRecord{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) error {
    seenIDs := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid record ID: %d (must be positive)", record.ID)
        }
        if record.Name == "" {
            return fmt.Errorf("empty name for record ID: %d", record.ID)
        }
        if record.Value < 0 {
            return fmt.Errorf("negative value for record ID: %d", record.ID)
        }
        if seenIDs[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        seenIDs[record.ID] = true
    }
    return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var max float64
    count := len(records)

    for i, record := range records {
        sum += record.Value
        if i == 0 || record.Value > max {
            max = record.Value
        }
    }

    average := sum / float64(count)
    return average, max, count
}package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) DataRecord {
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp.UTC(),
		Tags:      append(record.Tags, "processed"),
	}
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record, 1.5))
	}
	return processed, nil
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec001",
			Value:     42.5,
			Timestamp: time.Now(),
			Tags:      []string{"test"},
		},
		{
			ID:        "rec002",
			Value:     18.7,
			Timestamp: time.Now().Add(-time.Hour),
			Tags:      []string{"production"},
		},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, record := range processed {
		fmt.Printf("Processed: %+v\n", record)
	}
}