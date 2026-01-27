
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	emailRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return &DataProcessor{emailRegex: regex}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToLower(trimmed)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) RemoveSpecialChars(input string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return reg.ReplaceAllString(input, "")
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, string, bool) {
	cleanName := dp.CleanString(name)
	cleanEmail := dp.CleanString(email)
	validEmail := dp.ValidateEmail(cleanEmail)
	
	if !validEmail {
		return cleanName, cleanEmail, false
	}
	
	safeName := dp.RemoveSpecialChars(cleanName)
	return safeName, cleanEmail, true
}