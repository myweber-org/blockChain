
package main

import (
	"strings"
	"unicode"
)

// CleanInput removes extra whitespace and normalizes line endings
func CleanInput(input string) string {
	// Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	var result strings.Builder
	result.Grow(len(trimmed))
	spaceFlag := false
	
	for _, r := range trimmed {
		if unicode.IsSpace(r) {
			if !spaceFlag {
				result.WriteRune(' ')
				spaceFlag = true
			}
		} else {
			result.WriteRune(r)
			spaceFlag = false
		}
	}
	
	return result.String()
}

// NormalizeWhitespace converts all whitespace characters to single spaces
func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

// IsValidInput checks if input contains only printable characters
func IsValidInput(input string) bool {
	for _, r := range input {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}