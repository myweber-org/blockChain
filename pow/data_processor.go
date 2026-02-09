package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	return err == nil && matched
}

func TrimAndLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func ContainsSQLInjection(input string) bool {
	patterns := []string{
		`(?i)select.*from`,
		`(?i)insert.*into`,
		`(?i)update.*set`,
		`(?i)delete.*from`,
		`(?i)drop.*table`,
		`(?i)union.*select`,
		`--`,
		`;`,
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, input)
		if err == nil && matched {
			return true
		}
	}
	return false
}