
package main

import (
    "regexp"
    "strings"
)

type UserData struct {
    Username string
    Email    string
    Age      int
}

func ValidateUsername(username string) bool {
    if len(username) < 3 || len(username) > 20 {
        return false
    }
    validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    return validUsername.MatchString(username)
}

func SanitizeEmail(email string) string {
    trimmed := strings.TrimSpace(email)
    return strings.ToLower(trimmed)
}

func ValidateUserAge(age int) bool {
    return age >= 13 && age <= 120
}

func ProcessUserInput(data UserData) (UserData, error) {
    if !ValidateUsername(data.Username) {
        return UserData{}, ErrInvalidUsername
    }

    sanitizedEmail := SanitizeEmail(data.Email)

    if !ValidateUserAge(data.Age) {
        return UserData{}, ErrInvalidAge
    }

    return UserData{
        Username: data.Username,
        Email:    sanitizedEmail,
        Age:      data.Age,
    }, nil
}

var (
    ErrInvalidUsername = errors.New("invalid username format")
    ErrInvalidAge      = errors.New("age must be between 13 and 120")
)package main

import (
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Username string
	Email    string
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

func ValidateEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

func ProcessUserInput(username, email string) (*User, error) {
	if !ValidateUsername(username) {
		return nil, &ValidationError{Field: "username", Message: "invalid username format"}
	}

	sanitizedEmail := SanitizeEmail(email)
	if !ValidateEmail(sanitizedEmail) {
		return nil, &ValidationError{Field: "email", Message: "invalid email format"}
	}

	return &User{
		ID:       generateID(),
		Username: username,
		Email:    sanitizedEmail,
	}, nil
}

func generateID() int {
	return 1000
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}