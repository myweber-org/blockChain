
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
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) SanitizeInput(input string) string {
	dangerous := []string{"<", ">", "'", "\"", "&"}
	sanitized := input
	for _, char := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, char, "")
	}
	return strings.TrimSpace(sanitized)
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{"user@example.com", "USER@EXAMPLE.COM", "test@domain", "  user@example.com  "}
	unique := cleaner.RemoveDuplicates(records)
	fmt.Println("Unique records:", unique)
	
	emails := []string{"valid@email.com", "invalid", "user@com", "another@valid.org"}
	for _, email := range emails {
		fmt.Printf("%s validation: %v\n", email, cleaner.ValidateEmail(email))
	}
	
	input := "<script>alert('test')</script>   "
	sanitized := cleaner.SanitizeInput(input)
	fmt.Printf("Sanitized input: '%s'\n", sanitized)
}