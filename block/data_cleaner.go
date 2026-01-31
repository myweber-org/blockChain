package main

import (
	"strings"
)

func CleanData(input string) string {
	lines := strings.Split(input, "\n")
	seen := make(map[string]bool)
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, "\n")
}