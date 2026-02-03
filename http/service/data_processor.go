package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func processCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
	}

	return records, nil
}

func validateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, r := range records {
		if r.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d", r.ID)
		}
		if seenIDs[r.ID] {
			return fmt.Errorf("duplicate ID found: %d", r.ID)
		}
		seenIDs[r.ID] = true
	}
	return nil
}

func calculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, r := range records {
		sum += r.Value
		if r.Value > max {
			max = r.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        string
	Email     string
	Username  string
	Age       int
	Active    bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(profile UserProfile) error {
	if strings.TrimSpace(profile.ID) == "" {
		return errors.New("ID cannot be empty")
	}

	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	if profile.Age < 0 || profile.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(strings.TrimSpace(profile.Username))
	transformed.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	return transformed
}

func ProcessUserData(profile UserProfile) (UserProfile, error) {
	if err := ValidateProfile(profile); err != nil {
		return UserProfile{}, err
	}

	transformed := TransformProfile(profile)
	return transformed, nil
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func processCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(line) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(line))
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := line[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
		lineNumber++
	}

	return records, nil
}

func validateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d must be positive", record.ID)
		}
		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true
	}
	return nil
}

func calculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}package main

import (
	"errors"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Timestamp time.Time
	Value     float64
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

func TransformTags(tags []string) []string {
	var processed []string
	seen := make(map[string]bool)
	
	for _, tag := range tags {
		cleaned := strings.TrimSpace(strings.ToLower(tag))
		if cleaned != "" && !seen[cleaned] {
			seen[cleaned] = true
			processed = append(processed, cleaned)
		}
	}
	return processed
}

func CalculateAverage(records []DataRecord) (float64, error) {
	if len(records) == 0 {
		return 0, errors.New("cannot calculate average of empty slice")
	}
	
	var sum float64
	validCount := 0
	
	for _, record := range records {
		if err := ValidateRecord(record); err == nil {
			sum += record.Value
			validCount++
		}
	}
	
	if validCount == 0 {
		return 0, errors.New("no valid records found")
	}
	
	return sum / float64(validCount), nil
}

func FilterByTag(records []DataRecord, tagFilter string) []DataRecord {
	var filtered []DataRecord
	targetTag := strings.ToLower(strings.TrimSpace(tagFilter))
	
	for _, record := range records {
		for _, tag := range record.Tags {
			if strings.ToLower(tag) == targetTag {
				filtered = append(filtered, record)
				break
			}
		}
	}
	return filtered
}
package main

import (
	"errors"
	"strings"
	"unicode"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username can only contain letters, digits, underscores, and hyphens")
		}
	}
	return nil
}

func NormalizeEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	localPart := strings.ToLower(strings.TrimSpace(parts[0]))
	domain := strings.ToLower(strings.TrimSpace(parts[1]))
	return localPart + "@" + domain
}

func SanitizeInput(input string) string {
	replacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(input)
}