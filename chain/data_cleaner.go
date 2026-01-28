
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(input)
	fmt.Printf("Original: %v\n", input)
	fmt.Printf("Cleaned: %v\n", cleaned)
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSV(inputPath, outputPath string) error {
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

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		cleaned := make([]string, len(record))
		for i, field := range record {
			cleaned[i] = strings.TrimSpace(field)
		}
		if err := writer.Write(cleaned); err != nil {
			return fmt.Errorf("error writing CSV: %w", err)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}
	if err := cleanCSV(os.Args[1], os.Args[2]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("CSV cleaning completed successfully")
}package main

import "fmt"

func removeDuplicates(input []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    
    for _, value := range input {
        if !seen[value] {
            seen[value] = true
            result = append(result, value)
        }
    }
    
    return result
}

func main() {
    data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
    cleaned := removeDuplicates(data)
    fmt.Println("Original:", data)
    fmt.Println("Cleaned:", cleaned)
}
package main

import "fmt"

func removeDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"fmt"
	"sort"
)

func CleanData(data []string) []string {
	seen := make(map[string]struct{})
	unique := []string{}

	for _, item := range data {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			unique = append(unique, item)
		}
	}

	sort.Strings(unique)
	return unique
}

func main() {
	sample := []string{"zebra", "apple", "banana", "apple", "cherry", "banana"}
	cleaned := CleanData(sample)
	fmt.Println("Original:", sample)
	fmt.Println("Cleaned:", cleaned)
}
package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}
package main

import (
	"crypto/md5"
	"fmt"
	"sort"
)

type Record struct {
	ID   int
	Data string
}

func deduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, rec := range records {
		hash := fmt.Sprintf("%x", md5.Sum([]byte(rec.Data)))
		if !seen[hash] {
			seen[hash] = true
			unique = append(unique, rec)
		}
	}
	return unique
}

func validateRecords(records []Record) []Record {
	var valid []Record
	for _, rec := range records {
		if rec.ID > 0 && rec.Data != "" {
			valid = append(valid, rec)
		}
	}
	return valid
}

func sortRecordsByID(records []Record) []Record {
	sorted := make([]Record, len(records))
	copy(sorted, records)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID < sorted[j].ID
	})
	return sorted
}

func processDataPipeline(input []Record) []Record {
	validated := validateRecords(input)
	deduped := deduplicateRecords(validated)
	return sortRecordsByID(deduped)
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}