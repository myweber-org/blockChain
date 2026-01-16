
package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Email     string
	Timestamp time.Time
	Status    string
}

func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return fmt.Errorf("regex validation failed: %w", err)
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeString(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func ProcessRecord(record DataRecord) (DataRecord, error) {
	if record.ID == "" {
		return record, errors.New("record ID cannot be empty")
	}

	if err := ValidateEmail(record.Email); err != nil {
		return record, fmt.Errorf("email validation failed: %w", err)
	}

	record.Email = NormalizeString(record.Email)
	record.Status = NormalizeString(record.Status)

	if record.Timestamp.IsZero() {
		record.Timestamp = time.Now().UTC()
	}

	return record, nil
}

func TransformRecords(records []DataRecord) ([]DataRecord, []error) {
	var processed []DataRecord
	var errs []error

	for i, record := range records {
		processedRecord, err := ProcessRecord(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("record %d: %w", i, err))
			continue
		}
		processed = append(processed, processedRecord)
	}

	return processed, errs
}

func main() {
	records := []DataRecord{
		{ID: "001", Email: "USER@EXAMPLE.COM", Timestamp: time.Now(), Status: "ACTIVE"},
		{ID: "002", Email: "invalid-email", Status: "PENDING"},
		{ID: "", Email: "test@domain.com", Status: "inactive"},
	}

	processed, errs := TransformRecords(records)

	fmt.Printf("Processed %d records successfully\n", len(processed))
	if len(errs) > 0 {
		fmt.Printf("Encountered %d errors:\n", len(errs))
		for _, err := range errs {
			fmt.Printf("  - %v\n", err)
		}
	}

	for _, record := range processed {
		fmt.Printf("ID: %s, Email: %s, Status: %s\n",
			record.ID, record.Email, record.Status)
	}
}
package main

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
			Tags:      []string{"test", "sample"},
		},
		{
			ID:        "rec002",
			Value:     100.0,
			Timestamp: time.Now().Add(-time.Hour),
			Tags:      []string{"production"},
		},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, rec := range processed {
		fmt.Printf("Processed: %+v\n", rec)
	}
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
		fmt.Printf("Error processing data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data processing completed successfully")
}