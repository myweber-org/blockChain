
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

func (dc *DataCleaner) Clean(input []string) []string {
	var result []string
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !dc.seen[trimmed] {
			dc.seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func (dc *DataCleaner) Validate(item string) bool {
	return len(item) > 0 && len(item) <= 100
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"apple", " banana ", "apple", "", "cherry", "  "}
	cleaned := cleaner.Clean(data)
	fmt.Println("Cleaned data:", cleaned)
	for _, item := range cleaned {
		if cleaner.Validate(item) {
			fmt.Printf("'%s' is valid\n", item)
		}
	}
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
	slice := []int{1, 2, 2, 3, 4, 4, 5}
	unique := RemoveDuplicates(slice)
	fmt.Println("Original:", slice)
	fmt.Println("Unique:", unique)
}
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

func (dc *DataCleaner) Clean(input string) string {
	normalized := strings.ToLower(strings.TrimSpace(input))
	if dc.seen[normalized] {
		return ""
	}
	dc.seen[normalized] = true
	return normalized
}

func (dc *DataCleaner) ProcessList(items []string) []string {
	var cleaned []string
	for _, item := range items {
		if result := dc.Clean(item); result != "" {
			cleaned = append(cleaned, result)
		}
	}
	return cleaned
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"  Apple", "apple", "BANANA", "banana ", "Cherry", "cherry"}
	result := cleaner.ProcessList(data)
	fmt.Println("Cleaned data:", result)
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	for _, record := range records {
		if ValidateEmail(record.Email) {
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return RemoveDuplicates(cleaned)
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob", "invalid-email", false},
	}

	cleaned := CleanData(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	var unique []string
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] && dc.validateRecord(normalized) {
			dc.processedRecords[normalized] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func (dc *DataCleaner) validateRecord(record string) bool {
	if len(record) == 0 {
		return false
	}
	if strings.Contains(record, "test") {
		return false
	}
	return true
}

func (dc *DataCleaner) Reset() {
	dc.processedRecords = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	sampleData := []string{
		"Customer A",
		"customer a",
		"Customer B",
		"test record",
		"",
		"Customer C",
		"  Customer A  ",
	}
	
	fmt.Println("Original records:", len(sampleData))
	cleaned := cleaner.RemoveDuplicates(sampleData)
	fmt.Println("Cleaned records:", len(cleaned))
	
	for i, record := range cleaned {
		fmt.Printf("%d: %s\n", i+1, record)
	}
}