
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Record struct {
	ID    string
	Email string
	Phone string
}

type Cleaner struct {
	seenHashes map[string]bool
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		seenHashes: make(map[string]bool),
	}
}

func (c *Cleaner) NormalizeEmail(email string) string {
	parts := strings.Split(strings.ToLower(email), "@")
	if len(parts) != 2 {
		return ""
	}
	local := strings.Split(parts[0], "+")[0]
	local = strings.ReplaceAll(local, ".", "")
	return local + "@" + parts[1]
}

func (c *Cleaner) GenerateHash(record Record) string {
	normalizedEmail := c.NormalizeEmail(record.Email)
	if normalizedEmail == "" {
		return ""
	}
	data := fmt.Sprintf("%s|%s", normalizedEmail, strings.TrimSpace(record.Phone))
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (c *Cleaner) IsDuplicate(record Record) bool {
	hash := c.GenerateHash(record)
	if hash == "" {
		return false
	}
	if c.seenHashes[hash] {
		return true
	}
	c.seenHashes[hash] = true
	return false
}

func (c *Cleaner) ValidateRecord(record Record) bool {
	if len(record.ID) == 0 || len(record.Email) == 0 {
		return false
	}
	if !strings.Contains(record.Email, "@") {
		return false
	}
	if len(record.Phone) > 0 && !strings.HasPrefix(record.Phone, "+") {
		return false
	}
	return true
}

func (c *Cleaner) ProcessRecords(records []Record) []Record {
	var cleaned []Record
	for _, rec := range records {
		if !c.ValidateRecord(rec) {
			continue
		}
		if c.IsDuplicate(rec) {
			continue
		}
		cleaned = append(cleaned, rec)
	}
	return cleaned
}

func main() {
	cleaner := NewCleaner()
	records := []Record{
		{ID: "1", Email: "test@example.com", Phone: "+1234567890"},
		{ID: "2", Email: "TEST@example.com", Phone: "+1234567890"},
		{ID: "3", Email: "test+tag@example.com", Phone: "+1234567890"},
		{ID: "4", Email: "invalid-email", Phone: "+0987654321"},
		{ID: "5", Email: "another@test.com", Phone: ""},
	}

	result := cleaner.ProcessRecords(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(result))
	for _, rec := range result {
		fmt.Printf("ID: %s, Email: %s\n", rec.ID, rec.Email)
	}
}