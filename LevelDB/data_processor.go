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

type Record struct {
	ID      int
	Name    string
	Email   string
	Active  bool
	Score   float64
}

func ParseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(row) != 5 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 5, got %d", lineNum, len(row))
		}

		record, err := parseRow(row, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNum int) (Record, error) {
	var record Record
	var err error

	record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return Record{}, fmt.Errorf("empty name at line %d", lineNum)
	}

	record.Email = strings.TrimSpace(row[2])
	if !strings.Contains(record.Email, "@") {
		return Record{}, fmt.Errorf("invalid email format at line %d", lineNum)
	}

	record.Active, err = strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid active flag at line %d: %w", lineNum, err)
	}

	record.Score, err = strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid score at line %d: %w", lineNum, err)
	}

	if record.Score < 0 || record.Score > 100 {
		return Record{}, fmt.Errorf("score out of range at line %d: must be between 0 and 100", lineNum)
	}

	return record, nil
}

func ValidateRecords(records []Record) error {
	if len(records) == 0 {
		return errors.New("no records to validate")
	}

	emailSet := make(map[string]bool)
	idSet := make(map[int]bool)

	for _, record := range records {
		if emailSet[record.Email] {
			return fmt.Errorf("duplicate email found: %s", record.Email)
		}
		emailSet[record.Email] = true

		if idSet[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		idSet[record.ID] = true
	}

	return nil
}

func CalculateStatistics(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var maxScore float64
	activeCount := 0

	for _, record := range records {
		sum += record.Score
		if record.Score > maxScore {
			maxScore = record.Score
		}
		if record.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return average, maxScore, activeCount
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
	IsValid bool
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
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

		if len(row) < 4 {
			continue
		}

		record, err := parseRecord(row)
		if err != nil {
			fmt.Printf("Skipping line %d: %v\n", lineNumber, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRecord(row []string) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID format: %w", err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, fmt.Errorf("name cannot be empty")
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value format: %w", err)
	}
	record.Value = value

	isValid := strings.ToLower(strings.TrimSpace(row[3])) == "true"
	record.IsValid = isValid

	return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if record.IsValid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateAverageValue(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}

	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
}
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataProcessor struct {
	workers   int
	inputChan chan int
	outputChan chan int
	wg        sync.WaitGroup
	errChan   chan error
}

func NewDataProcessor(workers int) *DataProcessor {
	return &DataProcessor{
		workers:   workers,
		inputChan: make(chan int, 100),
		outputChan: make(chan int, 100),
		errChan:   make(chan error, workers),
	}
}

func (dp *DataProcessor) Start() {
	for i := 0; i < dp.workers; i++ {
		dp.wg.Add(1)
		go dp.worker(i)
	}
}

func (dp *DataProcessor) worker(id int) {
	defer dp.wg.Done()
	
	for data := range dp.inputChan {
		if data < 0 {
			dp.errChan <- fmt.Errorf("worker %d: negative value %d", id, data)
			continue
		}
		
		processed := dp.processData(data)
		dp.outputChan <- processed
		
		time.Sleep(10 * time.Millisecond)
	}
}

func (dp *DataProcessor) processData(value int) int {
	return value * 2
}

func (dp *DataProcessor) Submit(data []int) {
	for _, d := range data {
		dp.inputChan <- d
	}
	close(dp.inputChan)
}

func (dp *DataProcessor) Collect() ([]int, []error) {
	dp.wg.Wait()
	close(dp.outputChan)
	close(dp.errChan)
	
	results := []int{}
	for res := range dp.outputChan {
		results = append(results, res)
	}
	
	errors := []error{}
	for err := range dp.errChan {
		errors = append(errors, err)
	}
	
	return results, errors
}

func main() {
	processor := NewDataProcessor(3)
	processor.Start()
	
	testData := []int{1, 5, -3, 8, -1, 10}
	processor.Submit(testData)
	
	results, errs := processor.Collect()
	
	fmt.Println("Processing results:")
	for _, r := range results {
		fmt.Printf("Result: %d\n", r)
	}
	
	if len(errs) > 0 {
		fmt.Println("Errors encountered:")
		for _, e := range errs {
			fmt.Printf("Error: %v\n", e)
		}
	}
}