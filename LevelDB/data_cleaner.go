package datautils

func RemoveDuplicates[T comparable](slice []T) []T {
    seen := make(map[T]bool)
    result := []T{}
    for _, item := range slice {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    return result
}
package main

import (
	"fmt"
	"strings"
)

func CleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range input {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}

func main() {
	data := []string{" Apple", "banana", "apple", " Banana ", "Cherry"}
	cleaned := CleanData(data)
	fmt.Println("Cleaned data:", cleaned)
}