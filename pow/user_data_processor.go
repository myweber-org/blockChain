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
}package main

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

func ProcessUserInput(username, email string) (User, error) {
	if !ValidateUsername(username) {
		return User{}, ErrInvalidUsername
	}

	cleanEmail := SanitizeEmail(email)
	user := User{
		Username: username,
		Email:    cleanEmail,
	}
	return user, nil
}

var ErrInvalidUsername = errors.New("invalid username format")package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Password string
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

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	return hasUpper && hasLower && hasDigit
}

func ProcessUserData(data UserData) (UserData, error) {
	if !ValidateUsername(data.Username) {
		return UserData{}, ErrInvalidUsername
	}
	data.Email = SanitizeEmail(data.Email)
	if !ValidatePassword(data.Password) {
		return UserData{}, ErrWeakPassword
	}
	return data, nil
}

var (
	ErrInvalidUsername = errors.New("invalid username format")
	ErrWeakPassword    = errors.New("password does not meet security requirements")
)