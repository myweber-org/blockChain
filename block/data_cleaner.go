
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Phone string
	Valid bool
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := generateHash(record.Email + record.Phone)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateRecords(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		record.Valid = isValidEmail(record.Email) && isValidPhone(record.Phone)
		validated = append(validated, record)
	}
	return validated
}

func generateHash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidPhone(phone string) bool {
	return len(phone) >= 10 && strings.Count(phone, "") > 5
}

func processDataPipeline(records []DataRecord) []DataRecord {
	unique := deduplicateRecords(records)
	validated := validateRecords(unique)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{ID: "1", Email: "test@example.com", Phone: "1234567890"},
		{ID: "2", Email: "duplicate@example.com", Phone: "0987654321"},
		{ID: "3", Email: "test@example.com", Phone: "1234567890"},
		{ID: "4", Email: "invalid-email", Phone: "123"},
	}

	processed := processDataPipeline(sampleData)
	for _, record := range processed {
		fmt.Printf("ID: %s, Valid: %v\n", record.ID, record.Valid)
	}
}package datautils

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