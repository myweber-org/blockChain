
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	duplicates map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		duplicates: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	unique := []string{}
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.duplicates[normalized] {
			dc.duplicates[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) Reset() {
	dc.duplicates = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	emails := []string{
		"test@example.com",
		"TEST@example.com",
		"user@domain.org",
		"test@example.com",
		"invalid-email",
	}
	
	uniqueEmails := cleaner.RemoveDuplicates(emails)
	fmt.Println("Unique emails:", uniqueEmails)
	
	for _, email := range uniqueEmails {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("%s is valid\n", email)
		} else {
			fmt.Printf("%s is invalid\n", email)
		}
	}
	
	cleaner.Reset()
	fmt.Println("Cleaner has been reset")
}