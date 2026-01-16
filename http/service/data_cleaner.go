
package main

import (
	"fmt"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	return lower
}

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func ProcessData(data []string) []string {
	cleaned := []string{}
	for _, item := range data {
		cleaned = append(cleaned, CleanString(item))
	}
	return RemoveDuplicates(cleaned)
}

func main() {
	sampleData := []string{"  Apple ", "banana", "  APPLE", "Banana ", "Cherry"}
	result := ProcessData(sampleData)
	fmt.Println("Cleaned data:", result)
}package datautils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Remove any null characters
	cleaned := strings.ReplaceAll(input, "\x00", "")
	
	// Trim whitespace from both ends
	cleaned = strings.TrimSpace(cleaned)
	
	// Normalize multiple spaces to single space
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")
	
	// Remove any non-printable characters except standard whitespace
	var result strings.Builder
	for _, r := range cleaned {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

func NormalizeWhitespace(input string) string {
	// Replace various whitespace characters with standard space
	whitespaceRegex := regexp.MustCompile(`[\t\n\r\f\v]+`)
	normalized := whitespaceRegex.ReplaceAllString(input, " ")
	
	// Trim and collapse multiple spaces
	return strings.Join(strings.Fields(normalized), " ")
}
package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}