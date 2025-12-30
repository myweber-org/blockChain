
package main

import (
	"strings"
	"unicode"
)

// CleanInput removes extra whitespace and normalizes line endings
func CleanInput(input string) string {
	// Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	var result strings.Builder
	result.Grow(len(trimmed))
	spaceFlag := false
	
	for _, r := range trimmed {
		if unicode.IsSpace(r) {
			if !spaceFlag {
				result.WriteRune(' ')
				spaceFlag = true
			}
		} else {
			result.WriteRune(r)
			spaceFlag = false
		}
	}
	
	return result.String()
}

// NormalizeWhitespace converts all whitespace characters to single spaces
func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

// IsValidInput checks if input contains only printable characters
func IsValidInput(input string) bool {
	for _, r := range input {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
package main

import (
	"errors"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) DataRecord {
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp.UTC(),
		Tags:      append([]string{"processed"}, record.Tags...),
	}
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	
	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
}

func FilterByTag(records []DataRecord, tag string) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		for _, t := range record.Tags {
			if t == tag {
				filtered = append(filtered, record)
				break
			}
		}
	}
	return filtered
}