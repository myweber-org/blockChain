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

func (dc *DataCleaner) IsDuplicate(item string) bool {
	normalized := dc.Normalize(item)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) ProcessBatch(items []string) []string {
	var unique []string
	for _, item := range items {
		if !dc.IsDuplicate(item) {
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", "apple ", " BANANA", "banana", "Cherry"}
	
	fmt.Println("Original data:", data)
	
	cleaned := cleaner.ProcessBatch(data)
	fmt.Println("Cleaned data:", cleaned)
	
	cleaner.Reset()
	
	moreData := []string{"grape", "Grape", "PEACH"}
	secondBatch := cleaner.ProcessBatch(moreData)
	fmt.Println("Second batch:", secondBatch)
}