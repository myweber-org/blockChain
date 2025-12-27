
package main

import (
	"fmt"
	"strings"
)

func FilterAndUppercase(input []string, prefix string) []string {
	var result []string
	for _, s := range input {
		if strings.HasPrefix(s, prefix) {
			result = append(result, strings.ToUpper(s))
		}
	}
	return result
}

func main() {
	data := []string{"apple", "application", "banana", "appetizer"}
	filtered := FilterAndUppercase(data, "app")
	fmt.Println("Processed slice:", filtered)
}