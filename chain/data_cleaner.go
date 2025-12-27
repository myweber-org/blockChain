
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

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) ProcessBatch(items []string) []string {
	var uniqueItems []string
	for _, item := range items {
		if !dc.IsDuplicate(item) {
			uniqueItems = append(uniqueItems, item)
		}
	}
	return uniqueItems
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"  Apple  ",
		"apple",
		"BANANA",
		"  banana  ",
		"Cherry",
		"Date",
		"date",
	}
	
	result := cleaner.ProcessBatch(data)
	fmt.Println("Unique items:", result)
	fmt.Println("Total unique count:", len(result))
}