
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID      string
	Content string
	Hash    string
}

func generateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		if !seen[record.Hash] {
			seen[record.Hash] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateRecord(record DataRecord) bool {
	if strings.TrimSpace(record.ID) == "" {
		return false
	}
	if strings.TrimSpace(record.Content) == "" {
		return false
	}
	if record.Hash != generateHash(record.Content) {
		return false
	}
	return true
}

func cleanData(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if validateRecord(record) {
			validRecords = append(validRecords, record)
		}
	}
	return deduplicateRecords(validRecords)
}

func main() {
	records := []DataRecord{
		{ID: "001", Content: "Sample data", Hash: generateHash("Sample data")},
		{ID: "002", Content: "Duplicate data", Hash: generateHash("Duplicate data")},
		{ID: "003", Content: "Sample data", Hash: generateHash("Sample data")},
		{ID: "", Content: "Invalid record", Hash: generateHash("Invalid record")},
		{ID: "004", Content: "", Hash: generateHash("")},
	}

	cleaned := cleanData(records)
	fmt.Printf("Original: %d records\n", len(records))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processRecords(records []DataRecord) []DataRecord {
	emailMap := make(map[string]bool)
	var cleanRecords []DataRecord

	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		if validateEmail(cleanEmail) && !emailMap[cleanEmail] {
			emailMap[cleanEmail] = true
			record.Email = cleanEmail
			record.Valid = true
			cleanRecords = append(cleanRecords, record)
		}
	}
	return cleanRecords
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
		{5, "test@domain.org", false},
	}

	cleanRecords := processRecords(records)
	fmt.Printf("Processed %d records, %d valid unique records found\n", len(records), len(cleanRecords))
	
	for _, record := range cleanRecords {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
	}
}