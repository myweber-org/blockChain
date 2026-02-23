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

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []Record) []Record {
	var valid []Record
	for _, r := range records {
		if r.ID > 0 && r.Name != "" && r.Value >= 0 {
			valid = append(valid, r)
		}
	}
	return valid
}

func CalculateTotal(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		return
	}

	records, err := ProcessCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	validRecords := ValidateRecords(records)
	total := CalculateTotal(validRecords)

	fmt.Printf("Processed %d records, %d valid\n", len(records), len(validRecords))
	fmt.Printf("Total value: %.2f\n", total)
}
package main

import (
	"fmt"
	"math"
)

func processData(values []float64) (float64, float64, error) {
	if len(values) == 0 {
		return 0, 0, fmt.Errorf("empty data set")
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	var varianceSum float64
	for _, v := range values {
		diff := v - mean
		varianceSum += diff * diff
	}
	variance := varianceSum / float64(len(values))
	stdDev := math.Sqrt(variance)

	return mean, stdDev, nil
}

func main() {
	data := []float64{2.5, 3.7, 1.8, 4.2, 3.1}
	avg, dev, err := processData(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Average: %.2f\n", avg)
	fmt.Printf("Standard Deviation: %.2f\n", dev)
}