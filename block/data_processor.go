
package utils

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToValidUTF8(trimmed, "")
}

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}