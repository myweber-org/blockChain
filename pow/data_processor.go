package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Tags      []string  `json:"tags"`
}

type ProcessedRecord struct {
	RecordID  string    `json:"record_id"`
	Valid     bool      `json:"valid"`
	Domain    string    `json:"domain,omitempty"`
	NormalizedValue float64 `json:"normalized_value"`
	TagCount  int       `json:"tag_count"`
	AgeDays   int       `json:"age_days"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ExtractDomain(email string) (string, bool) {
	if !ValidateEmail(email) {
		return "", false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", false
	}
	return parts[1], true
}

func NormalizeValue(value float64, min, max float64) float64 {
	if max <= min {
		return 0.0
	}
	normalized := (value - min) / (max - min)
	if normalized < 0 {
		return 0.0
	}
	if normalized > 1 {
		return 1.0
	}
	return normalized
}

func ProcessRecord(record DataRecord, minValue, maxValue float64) ProcessedRecord {
	domain, hasDomain := ExtractDomain(record.Email)
	ageDays := int(time.Since(record.Timestamp).Hours() / 24)
	if ageDays < 0 {
		ageDays = 0
	}

	return ProcessedRecord{
		RecordID:        record.ID,
		Valid:           hasDomain,
		Domain:          domain,
		NormalizedValue: NormalizeValue(record.Value, minValue, maxValue),
		TagCount:        len(record.Tags),
		AgeDays:         ageDays,
	}
}

func ProcessBatch(records []DataRecord, minValue, maxValue float64) []ProcessedRecord {
	processed := make([]ProcessedRecord, len(records))
	for i, record := range records {
		processed[i] = ProcessRecord(record, minValue, maxValue)
	}
	return processed
}

func GenerateSummary(processed []ProcessedRecord) map[string]interface{} {
	validCount := 0
	totalTags := 0
	domains := make(map[string]int)

	for _, record := range processed {
		if record.Valid {
			validCount++
			domains[record.Domain]++
		}
		totalTags += record.TagCount
	}

	avgTags := 0.0
	if len(processed) > 0 {
		avgTags = float64(totalTags) / float64(len(processed))
	}

	return map[string]interface{}{
		"total_records":   len(processed),
		"valid_records":   validCount,
		"invalid_records": len(processed) - validCount,
		"unique_domains":  len(domains),
		"average_tags":    fmt.Sprintf("%.2f", avgTags),
	}
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec001",
			Email:     "user@example.com",
			Timestamp: time.Now().Add(-48 * time.Hour),
			Value:     75.5,
			Tags:      []string{"premium", "active"},
		},
		{
			ID:        "rec002",
			Email:     "invalid-email",
			Timestamp: time.Now().Add(-24 * time.Hour),
			Value:     120.0,
			Tags:      []string{"trial"},
		},
		{
			ID:        "rec003",
			Email:     "admin@company.org",
			Timestamp: time.Now().Add(-72 * time.Hour),
			Value:     45.0,
			Tags:      []string{"admin", "verified", "staff"},
		},
	}

	processed := ProcessBatch(records, 0.0, 100.0)
	summary := GenerateSummary(processed)

	summaryJSON, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println("Processing Summary:")
	fmt.Println(string(summaryJSON))

	fmt.Println("\nProcessed Records:")
	for _, record := range processed {
		recordJSON, _ := json.MarshalIndent(record, "", "  ")
		fmt.Println(string(recordJSON))
	}
}