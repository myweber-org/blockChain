
package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
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
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, line := range input {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	fmt.Println("Enter data lines (press Ctrl+D to finish):")
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	cleaned := cleanData(lines)
	fmt.Println("\nCleaned data:")
	for _, item := range cleaned {
		fmt.Println(item)
	}
}