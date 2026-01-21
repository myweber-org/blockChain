package utils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)

	// Remove any HTML/XML tags
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(trimmed, "")

	// Escape potentially dangerous characters
	re = regexp.MustCompile(`['"\\;]`)
	escaped := re.ReplaceAllStringFunc(cleaned, func(match string) string {
		return "\\" + match
	})

	return escaped
}