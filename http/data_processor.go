package data_processor

import (
	"regexp"
	"strings"
)

type DataCleaner struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dc *DataCleaner) NormalizeString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := dc.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return strings.ToLower(normalized)
}

func (dc *DataCleaner) RemoveSpecialChars(input string, keepChars string) string {
	pattern := "[^a-zA-Z0-9"
	if keepChars != "" {
		pattern += regexp.QuoteMeta(keepChars)
	}
	pattern += "]+"
	
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, "")
}

func (dc *DataCleaner) Tokenize(input string, delimiter string) []string {
	normalized := dc.NormalizeString(input)
	if normalized == "" {
		return []string{}
	}
	
	return strings.Split(normalized, delimiter)
}