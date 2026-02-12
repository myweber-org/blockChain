package main

import (
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Username string
	Email    string
	Age      int
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validPattern.MatchString(username)
}

func SanitizeEmail(email string) string {
	trimmed := strings.TrimSpace(email)
	return strings.ToLower(trimmed)
}

func ValidateUserAge(age int) bool {
	return age >= 13 && age <= 120
}

func ProcessUserInput(username, email string, age int) (User, error) {
	if !ValidateUsername(username) {
		return User{}, &ValidationError{Field: "username", Message: "invalid username format"}
	}

	sanitizedEmail := SanitizeEmail(email)
	if !strings.Contains(sanitizedEmail, "@") {
		return User{}, &ValidationError{Field: "email", Message: "invalid email address"}
	}

	if !ValidateUserAge(age) {
		return User{}, &ValidationError{Field: "age", Message: "age must be between 13 and 120"}
	}

	return User{
		Username: username,
		Email:    sanitizedEmail,
		Age:      age,
	}, nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}