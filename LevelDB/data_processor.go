package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
	Timestamp string `json:"timestamp"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func TransformProfile(profile UserProfile) (UserProfile, error) {
	if profile.Age < 0 || profile.Age > 150 {
		return profile, fmt.Errorf("invalid age: %d", profile.Age)
	}

	if !ValidateEmail(profile.Email) {
		return profile, fmt.Errorf("invalid email format: %s", profile.Email)
	}

	profile.Username = SanitizeUsername(profile.Username)

	if profile.Timestamp == "" {
		profile.Timestamp = "2024-01-01T00:00:00Z"
	}

	return profile, nil
}

func ProcessUserData(inputJSON string) (string, error) {
	var profile UserProfile
	err := json.Unmarshal([]byte(inputJSON), &profile)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	transformedProfile, err := TransformProfile(profile)
	if err != nil {
		return "", fmt.Errorf("validation failed: %v", err)
	}

	outputJSON, err := json.MarshalIndent(transformedProfile, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return string(outputJSON), nil
}

func main() {
	sampleInput := `{
		"id": 1,
		"username": "  JohnDoe  ",
		"email": "john@example.com",
		"age": 30,
		"active": true,
		"timestamp": ""
	}`

	result, err := ProcessUserData(sampleInput)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Processed user profile:")
	fmt.Println(result)
}