
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ExtractDomain(email string) (string, bool) {
	if !dp.ValidateEmail(email) {
		return "", false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", false
	}
	return parts[1], true
}

func (dp *DataProcessor) NormalizeWhitespace(input string) string {
	return dp.whitespaceRegex.ReplaceAllString(input, " ")
}
package main

import (
	"regexp"
	"strings"
)

func ProcessInput(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", nil
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	cleaned := re.ReplaceAllString(trimmed, "")

	reMultipleSpaces := regexp.MustCompile(`\s+`)
	final := reMultipleSpaces.ReplaceAllString(cleaned, " ")

	return strings.ToLower(final), nil
}
package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int
	Email     string
	Username  string
	BirthDate string
	Active    bool
	CreatedAt time.Time
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func CalculateAge(birthDate string) (int, error) {
	parsedDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	age := time.Since(parsedDate).Hours() / 24 / 365.25
	return int(age), nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateEmail(profile.Email); err != nil {
		return profile, err
	}

	profile.Username = SanitizeUsername(profile.Username)

	age, err := CalculateAge(profile.BirthDate)
	if err != nil {
		return profile, err
	}

	if age < 13 {
		return profile, errors.New("user must be at least 13 years old")
	}

	profile.CreatedAt = time.Now()
	return profile, nil
}

func main() {
	profile := UserProfile{
		ID:        1,
		Email:     "user@example.com",
		Username:  "  john_doe  ",
		BirthDate: "1990-05-15",
		Active:    true,
	}

	processedProfile, err := ProcessUserProfile(profile)
	if err != nil {
		println("Error processing profile:", err.Error())
		return
	}

	println("Profile processed successfully")
	println("Username:", processedProfile.Username)
	println("Created at:", processedProfile.CreatedAt.Format("2006-01-02"))
}package main

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

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func NormalizeTags(tags []string) []string {
	uniqueTags := make(map[string]bool)
	var normalized []string

	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		lower := strings.ToLower(trimmed)
		if trimmed != "" && !uniqueTags[lower] {
			uniqueTags[lower] = true
			normalized = append(normalized, trimmed)
		}
	}
	return normalized
}

func TransformRecord(record DataRecord) (DataRecord, error) {
	if !ValidateEmail(record.Email) {
		return DataRecord{}, fmt.Errorf("invalid email format: %s", record.Email)
	}

	if record.Value < 0 {
		record.Value = 0
	}

	record.Tags = NormalizeTags(record.Tags)

	if record.ID == "" {
		record.ID = fmt.Sprintf("rec_%d", time.Now().UnixNano())
	}

	return record, nil
}

func ProcessRecords(records []DataRecord) ([]DataRecord, []error) {
	var processed []DataRecord
	var errors []error

	for _, record := range records {
		transformed, err := TransformRecord(record)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		processed = append(processed, transformed)
	}

	return processed, errors
}

func RecordsToJSON(records []DataRecord) (string, error) {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}