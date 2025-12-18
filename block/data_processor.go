
package main

import (
	"fmt"
)

// FilterAndDouble filters out negative numbers from a slice and doubles the remaining values.
func FilterAndDouble(numbers []int) []int {
	var result []int
	for _, num := range numbers {
		if num >= 0 {
			result = append(result, num*2)
		}
	}
	return result
}

func main() {
	input := []int{-5, 2, 0, 7, -1, 4}
	output := FilterAndDouble(input)
	fmt.Printf("Input: %v\n", input)
	fmt.Printf("Output: %v\n", output)
}