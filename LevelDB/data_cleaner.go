package utils

import "strings"

func TrimWhitespaceFromSlice(slice []string) []string {
    trimmed := make([]string, len(slice))
    for i, s := range slice {
        trimmed[i] = strings.TrimSpace(s)
    }
    return trimmed
}
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Record struct {
	ID   string
	Data string
}

type Cleaner struct {
	seen map[string]bool
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		seen: make(map[string]bool),
	}
}

func (c *Cleaner) GenerateHash(data string) string {
	hash := sha256.Sum256([]byte(strings.TrimSpace(data)))
	return hex.EncodeToString(hash[:])
}

func (c *Cleaner) IsDuplicate(data string) bool {
	hash := c.GenerateHash(data)
	if c.seen[hash] {
		return true
	}
	c.seen[hash] = true
	return false
}

func (c *Cleaner) ValidateRecord(record Record) error {
	if record.ID == "" {
		return fmt.Errorf("record ID cannot be empty")
	}
	if len(record.Data) > 1000 {
		return fmt.Errorf("data exceeds maximum length")
	}
	return nil
}

func (c *Cleaner) ProcessRecords(records []Record) []Record {
	var cleaned []Record
	for _, record := range records {
		if err := c.ValidateRecord(record); err != nil {
			fmt.Printf("Validation failed for %s: %v\n", record.ID, err)
			continue
		}
		if c.IsDuplicate(record.Data) {
			fmt.Printf("Duplicate detected for %s\n", record.ID)
			continue
		}
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func main() {
	cleaner := NewCleaner()
	records := []Record{
		{ID: "001", Data: "Sample data"},
		{ID: "002", Data: "Sample data"},
		{ID: "003", Data: ""},
		{ID: "004", Data: strings.Repeat("x", 1001)},
		{ID: "005", Data: "Unique content"},
	}
	result := cleaner.ProcessRecords(records)
	fmt.Printf("Cleaned records: %d\n", len(result))
}