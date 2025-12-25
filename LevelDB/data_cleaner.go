
package main

import (
	"strings"
)

// CleanString removes duplicate spaces and trims leading/trailing whitespace
func CleanString(input string) string {
	// Remove duplicate spaces
	var result strings.Builder
	prevSpace := false
	for _, r := range input {
		if r == ' ' {
			if !prevSpace {
				result.WriteRune(r)
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}
	
	// Trim and return
	return strings.TrimSpace(result.String())
}

// RemoveDuplicates removes duplicate strings from a slice
func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// CleanSlice applies CleanString to each element and removes duplicates
func CleanSlice(slice []string) []string {
	cleaned := make([]string, len(slice))
	for i, s := range slice {
		cleaned[i] = CleanString(s)
	}
	return RemoveDuplicates(cleaned)
}
package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeInput(input string) string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)

	// Remove any null characters
	trimmed = strings.ReplaceAll(trimmed, "\x00", "")

	// Normalize multiple spaces to single space
	spaceRegex := regexp.MustCompile(`\s+`)
	trimmed = spaceRegex.ReplaceAllString(trimmed, " ")

	// Remove any non-printable characters except newline and tab
	var result strings.Builder
	for _, r := range trimmed {
		if unicode.IsPrint(r) || r == '\n' || r == '\t' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func NormalizeWhitespace(input string) string {
	// Replace various whitespace characters with standard space
	whitespaceRegex := regexp.MustCompile(`[\t\n\r\f\v]+`)
	normalized := whitespaceRegex.ReplaceAllString(input, " ")
	
	// Collapse multiple spaces
	normalized = strings.Join(strings.Fields(normalized), " ")
	
	return normalized
}