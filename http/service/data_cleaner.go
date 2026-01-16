
package main

import (
	"fmt"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	return lower
}

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func ProcessData(data []string) []string {
	cleaned := []string{}
	for _, item := range data {
		cleaned = append(cleaned, CleanString(item))
	}
	return RemoveDuplicates(cleaned)
}

func main() {
	sampleData := []string{"  Apple ", "banana", "  APPLE", "Banana ", "Cherry"}
	result := ProcessData(sampleData)
	fmt.Println("Cleaned data:", result)
}