
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ExtractDomain(email string) (string, bool) {
	if !dp.ValidateEmail(email) {
		return "", false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", false
	}
	return parts[1], true
}

func (dp *DataProcessor) NormalizeWhitespace(input string) string {
	return dp.whitespaceRegex.ReplaceAllString(input, " ")
}
package main

import (
	"regexp"
	"strings"
)

func ProcessInput(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", nil
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	cleaned := re.ReplaceAllString(trimmed, "")

	reMultipleSpaces := regexp.MustCompile(`\s+`)
	final := reMultipleSpaces.ReplaceAllString(cleaned, " ")

	return strings.ToLower(final), nil
}