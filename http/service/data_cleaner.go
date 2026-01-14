package main

import "fmt"

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

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}package datautils

import "sort"

func RemoveDuplicates[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(input))

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func RemoveDuplicatesSorted[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	result := make([]T, 0, len(input))
	sort.Slice(input, func(i, j int) bool {
		switch v := any(input[i]).(type) {
		case int:
			return v < any(input[j]).(int)
		case string:
			return v < any(input[j]).(string)
		default:
			return false
		}
	})

	prev := input[0]
	result = append(result, prev)

	for i := 1; i < len(input); i++ {
		if input[i] != prev {
			result = append(result, input[i])
			prev = input[i]
		}
	}

	return result
}