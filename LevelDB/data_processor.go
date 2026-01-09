
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor(pattern string) (*DataProcessor, error) {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &DataProcessor{allowedPattern: compiled}, nil
}

func (dp *DataProcessor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.allowedPattern.FindString(trimmed)
}

func (dp *DataProcessor) Validate(input string) bool {
	return dp.allowedPattern.MatchString(input)
}

func (dp *DataProcessor) ProcessBatch(inputs []string) []string {
	var results []string
	for _, item := range inputs {
		cleaned := dp.CleanInput(item)
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}