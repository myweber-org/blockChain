package datautil

import "sort"

func RemoveDuplicates[T comparable](input []T) []T {
    if len(input) == 0 {
        return input
    }

    seen := make(map[T]bool)
    result := make([]T, 0, len(input))

    for _, item := range input {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }

    return result
}

func RemoveDuplicatesSorted[T comparable](input []T) []T {
    if len(input) == 0 {
        return input
    }

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

    result := make([]T, 0, len(input))
    result = append(result, input[0])

    for i := 1; i < len(input); i++ {
        if input[i] != input[i-1] {
            result = append(result, input[i])
        }
    }

    return result
}
package main

import "fmt"

func removeDuplicates(nums []int) []int {
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
	cleaned := removeDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}