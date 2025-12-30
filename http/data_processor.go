
package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(trimmed, " ")
	return cleaned
}

func ValidateInput(input string) bool {
	if len(input) == 0 {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9\s.,!?]+$`)
	return re.MatchString(input)
}

func ProcessData(input string) (string, bool) {
	cleaned := CleanInput(input)
	if !ValidateInput(cleaned) {
		return "", false
	}
	return cleaned, true
}