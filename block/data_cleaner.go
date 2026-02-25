package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput removes potentially harmful characters and trims whitespace
func SanitizeInput(input string) string {
	// Trim leading and trailing whitespace
	cleaned := strings.TrimSpace(input)
	
	// Remove any HTML/XML tags
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned = re.ReplaceAllString(cleaned, "")
	
	// Remove control characters except newline and tab
	re = regexp.MustCompile(`[\x00-\x08\x0B-\x0C\x0E-\x1F\x7F]`)
	cleaned = re.ReplaceAllString(cleaned, "")
	
	// Limit length to prevent excessive input
	maxLength := 1000
	if len(cleaned) > maxLength {
		cleaned = cleaned[:maxLength]
	}
	
	return cleaned
}

// NormalizeWhitespace converts multiple whitespace characters to single spaces
func NormalizeWhitespace(input string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(input, " ")
}