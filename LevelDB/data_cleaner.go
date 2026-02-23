
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