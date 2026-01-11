
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

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(item) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && len(item) < 100
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"apple", "Apple", "banana", "  BANANA  ", "", "cherry", "cherry"}
	
	cleaned := cleaner.RemoveDuplicates(data)
	
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
	
	cleaner.Reset()
	
	moreData := []string{"grape", "GRAPE", "kiwi"}
	secondBatch := cleaner.RemoveDuplicates(moreData)
	
	fmt.Println("Second batch:", secondBatch)
}