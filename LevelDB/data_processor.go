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
	Tags      []string `json:"tags"`
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
	if profile.Age < 0 || profile.Age > 120 {
		return profile, fmt.Errorf("invalid age: %d", profile.Age)
	}

	if !ValidateEmail(profile.Email) {
		return profile, fmt.Errorf("invalid email format: %s", profile.Email)
	}

	profile.Username = SanitizeUsername(profile.Username)

	if len(profile.Tags) > 10 {
		profile.Tags = profile.Tags[:10]
	}

	return profile, nil
}

func ProcessUserData(jsonData []byte) ([]byte, error) {
	var profiles []UserProfile
	if err := json.Unmarshal(jsonData, &profiles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var validProfiles []UserProfile
	for _, profile := range profiles {
		transformed, err := TransformProfile(profile)
		if err != nil {
			fmt.Printf("Skipping profile ID %d: %v\n", profile.ID, err)
			continue
		}
		validProfiles = append(validProfiles, transformed)
	}

	result, err := json.MarshalIndent(validProfiles, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return result, nil
}

func main() {
	sampleData := `[
		{"id":1,"username":"  JohnDoe  ","email":"john@example.com","age":25,"active":true,"tags":["golang","backend"]},
		{"id":2,"username":"JaneSmith","email":"invalid-email","age":150,"active":false,"tags":["frontend","design","test","extra"]},
		{"id":3,"username":"Bob","email":"bob@test.org","age":30,"active":true,"tags":[]}
	]`

	processed, err := ProcessUserData([]byte(sampleData))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Println("Processed user profiles:")
	fmt.Println(string(processed))
}