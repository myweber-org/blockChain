package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput removes leading/trailing whitespace, reduces multiple spaces to single,
// and strips potentially dangerous characters from user input strings.
func SanitizeInput(input string) string {
	// Trim whitespace from both ends
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned := spaceRegex.ReplaceAllString(trimmed, " ")
	
	// Remove special characters that could be used for injection attacks
	// Allow alphanumeric, spaces, and common punctuation
	specialRegex := regexp.MustCompile(`[^a-zA-Z0-9\s.,!?-]`)
	sanitized := specialRegex.ReplaceAllString(cleaned, "")
	
	return sanitized
}

// ValidateEmail checks if a string is a valid email address format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}