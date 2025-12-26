
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
}