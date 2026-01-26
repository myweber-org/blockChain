package csvutils

import (
	"strings"
	"unicode"
)

// SanitizeString removes potentially problematic characters from CSV field values.
// It trims spaces, removes newlines and carriage returns, and replaces
// consecutive spaces with a single space.
func SanitizeString(input string) string {
	// Remove leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove newlines and carriage returns
	noNewlines := strings.ReplaceAll(trimmed, "\n", "")
	noNewlines = strings.ReplaceAll(noNewlines, "\r", "")
	
	// Replace multiple spaces with single space
	var result strings.Builder
	result.Grow(len(noNewlines))
	
	prevSpace := false
	for _, r := range noNewlines {
		if unicode.IsSpace(r) {
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}
	
	return result.String()
}