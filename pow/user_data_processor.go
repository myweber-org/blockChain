package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Password string
}

func ValidateUserData(data UserData) (bool, []string) {
	var errors []string

	if len(data.Username) < 3 || len(data.Username) > 20 {
		errors = append(errors, "Username must be between 3 and 20 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(data.Email) {
		errors = append(errors, "Invalid email format")
	}

	if len(data.Password) < 8 {
		errors = append(errors, "Password must be at least 8 characters")
	}

	if len(errors) > 0 {
		return false, errors
	}
	return true, nil
}

func SanitizeUsername(username string) string {
	sanitized := strings.TrimSpace(username)
	sanitized = regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(sanitized, "")
	return sanitized
}

func ProcessUserInput(username, email, password string) (UserData, bool, []string) {
	userData := UserData{
		Username: SanitizeUsername(username),
		Email:    strings.TrimSpace(email),
		Password: password,
	}

	valid, errors := ValidateUserData(userData)
	return userData, valid, errors
}