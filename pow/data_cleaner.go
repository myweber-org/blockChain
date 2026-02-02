package utils

import (
	"regexp"
	"strings"
)

// SanitizeString removes leading/trailing whitespace and collapses multiple spaces
func SanitizeString(input string) string {
	// Trim spaces from start and end
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(trimmed, " ")
}