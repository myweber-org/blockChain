
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
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func calculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	for _, r := range records {
		sum += r.Value
	}
	average := sum / float64(len(records))

	var variance float64
	for _, r := range records {
		diff := r.Value - average
		variance += diff * diff
	}
	variance = variance / float64(len(records))

	return average, variance
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := processCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d records\n", len(records))

	avg, variance := calculateStats(records)
	fmt.Printf("Average: %.2f, Variance: %.2f\n", avg, variance)
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	Metadata  string    `json:"metadata,omitempty"`
}

func ValidateUserProfile(profile UserProfile) error {
	if profile.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if len(profile.Username) < 3 || len(profile.Username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return fmt.Errorf("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}

	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(transformed.Username)
	transformed.Email = strings.ToLower(transformed.Email)
	transformed.Metadata = strings.TrimSpace(transformed.Metadata)

	if transformed.Metadata == "" {
		transformed.Metadata = "default_metadata"
	}

	return transformed
}

func ProcessUserProfile(data []byte) (UserProfile, error) {
	var profile UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return UserProfile{}, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	if err := ValidateUserProfile(profile); err != nil {
		return UserProfile{}, fmt.Errorf("validation failed: %w", err)
	}

	transformedProfile := TransformProfile(profile)
	return transformedProfile, nil
}

func main() {
	jsonData := []byte(`{
		"id": 1,
		"username": "TestUser",
		"email": "TEST@EXAMPLE.COM",
		"age": 25,
		"created_at": "2023-10-01T12:00:00Z",
		"metadata": "   additional info   "
	}`)

	processedProfile, err := ProcessUserProfile(jsonData)
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}

	fmt.Printf("Processed Profile: %+v\n", processedProfile)
}