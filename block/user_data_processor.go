
package main

import (
	"regexp"
	"strings"
)

type UserDataProcessor struct {
	allowedChars *regexp.Regexp
	maxLength    int
}

func NewUserDataProcessor() *UserDataProcessor {
	re := regexp.MustCompile(`^[a-zA-Z0-9\s\-_@.]+$`)
	return &UserDataProcessor{
		allowedChars: re,
		maxLength:    100,
	}
}

func (p *UserDataProcessor) ValidateUsername(input string) (bool, string) {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) == 0 {
		return false, "Username cannot be empty"
	}
	if len(trimmed) > p.maxLength {
		return false, "Username exceeds maximum length"
	}
	if !p.allowedChars.MatchString(trimmed) {
		return false, "Username contains invalid characters"
	}
	return true, ""
}

func (p *UserDataProcessor) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(trimmed, " ")
}

func (p *UserDataProcessor) ProcessUserInput(rawInput string) (string, error) {
	sanitized := p.SanitizeInput(rawInput)
	valid, errMsg := p.ValidateUsername(sanitized)
	if !valid {
		return "", &InputError{Message: errMsg}
	}
	return sanitized, nil
}

type InputError struct {
	Message string
}

func (e *InputError) Error() string {
	return e.Message
}package main

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
		return data, ErrInvalidUsername
	}

	data.Email = SanitizeEmail(data.Email)

	if !ValidateUserAge(data.Age) {
		return data, ErrInvalidAge
	}

	return data, nil
}

var (
	ErrInvalidUsername = errors.New("invalid username format")
	ErrInvalidAge      = errors.New("age must be between 13 and 120")
)