package main

import (
	"fmt"
	"strings"
	"unicode"
)

type UserProfile struct {
	Username string
	Email    string
	Age      int
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateAge(age int) bool {
	return age >= 0 && age <= 120
}

func processUserData(username, email string, age int) (*UserProfile, error) {
	normalizedUsername := normalizeUsername(username)
	if normalizedUsername == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	if !validateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}
	normalizedEmail := normalizeEmail(email)

	if !validateAge(age) {
		return nil, fmt.Errorf("age must be between 0 and 120")
	}

	return &UserProfile{
		Username: normalizedUsername,
		Email:    normalizedEmail,
		Age:      age,
	}, nil
}

func main() {
	user, err := processUserData("  JohnDoe  ", "JOHN@EXAMPLE.COM", 30)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}