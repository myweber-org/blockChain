
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedCount int
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{processedCount: 0}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	dc.processedCount += len(items)
	return result
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

func (dc *DataCleaner) GetStats() int {
	return dc.processedCount
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"  apple ", "banana", "apple", "", "  banana  ", "cherry"}
	unique := cleaner.RemoveDuplicates(data)
	fmt.Println("Unique items:", unique)
	
	emails := []string{"test@example.com", "invalid-email", "user@domain"}
	for _, email := range emails {
		valid := cleaner.ValidateEmail(email)
		fmt.Printf("Email %s valid: %v\n", email, valid)
	}
	
	fmt.Printf("Total processed items: %d\n", cleaner.GetStats())
}
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
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