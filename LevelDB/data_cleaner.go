
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(item) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && !strings.ContainsAny(item, "!@#$%")
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"apple", "Apple", "banana", "", "cherry!", "banana", "date"}
	
	cleaned := cleaner.RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
	
	cleaner.Reset()
	moreData := []string{"grape", "Grape", "kiwi"}
	cleaned2 := cleaner.RemoveDuplicates(moreData)
	fmt.Println("Second batch:", cleaned2)
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
	ID    int
	Name  string
	Email string
	Score float64
}

func cleanString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []DataRecord
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNum++
		if lineNum == 1 {
			continue
		}

		if len(line) < 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(line[0]))
		if err != nil {
			continue
		}

		name := cleanString(line[1])
		email := cleanString(line[2])
		score, err := strconv.ParseFloat(strings.TrimSpace(line[3]), 64)
		if err != nil {
			continue
		}

		if !validateEmail(email) {
			continue
		}

		record := DataRecord{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
		}
		records = append(records, record)
	}

	return records, nil
}

func calculateAverageScore(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}

	total := 0.0
	for _, record := range records {
		total += record.Score
	}
	return total / float64(len(records))
}

func filterHighScorers(records []DataRecord, threshold float64) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Score >= threshold {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_cleaner <csv_file>")
		return
	}

	filename := os.Args[1]
	records, err := parseCSVFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("Successfully parsed %d valid records\n", len(records))
	
	avgScore := calculateAverageScore(records)
	fmt.Printf("Average score: %.2f\n", avgScore)

	highScorers := filterHighScorers(records, 80.0)
	fmt.Printf("High scorers (>=80): %d\n", len(highScorers))

	for _, record := range highScorers {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Score: %.1f\n",
			record.ID, record.Name, record.Email, record.Score)
	}
}