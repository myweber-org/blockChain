
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
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) CleanPhoneNumber(phone string) string {
	var builder strings.Builder
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

func main() {
	cleaner := NewDataCleaner()

	emails := []string{
		"test@example.com",
		"TEST@example.com",
		"user@domain.org",
		"test@example.com",
		"  invalid  ",
	}

	uniqueEmails := cleaner.RemoveDuplicates(emails)
	fmt.Println("Unique emails:", uniqueEmails)

	for _, email := range uniqueEmails {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("Valid email: %s\n", email)
		} else {
			fmt.Printf("Invalid email: %s\n", email)
		}
	}

	phone := "+1 (234) 567-8900"
	cleanedPhone := cleaner.CleanPhoneNumber(phone)
	fmt.Printf("Original: %s -> Cleaned: %s\n", phone, cleanedPhone)
}