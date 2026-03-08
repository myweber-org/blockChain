package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput cleans and normalizes user-provided string input
func SanitizeInput(input string) string {
	// Remove leading/trailing whitespace
	cleaned := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")
	
	// Remove potentially dangerous characters (basic example)
	dangerousChars := regexp.MustCompile(`[<>{}]`)
	cleaned = dangerousChars.ReplaceAllString(cleaned, "")
	
	return cleaned
}