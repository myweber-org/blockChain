
package data

import (
	"regexp"
	"strings"
)

// ValidateEmail checks if the provided string is a valid email address.
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// SanitizeInput removes leading and trailing whitespace from a string.
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// TruncateString shortens a string to the specified maximum length.
// If the string is shorter than maxLen, it is returned unchanged.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}