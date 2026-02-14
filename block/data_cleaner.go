
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Clean(input []string) []string {
	var result []string
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !dc.seen[trimmed] {
			dc.seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func (dc *DataCleaner) Validate(item string) bool {
	return len(item) > 0 && len(item) <= 100
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"apple", " banana ", "apple", "", "cherry", "  "}
	cleaned := cleaner.Clean(data)
	fmt.Println("Cleaned data:", cleaned)
	for _, item := range cleaned {
		if cleaner.Validate(item) {
			fmt.Printf("'%s' is valid\n", item)
		}
	}
}