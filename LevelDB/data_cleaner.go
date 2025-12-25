
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