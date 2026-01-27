
package main

import (
	"bufio"
	"fmt"
	"os"
)

func removeDuplicates(inputFile, outputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	seen := make(map[string]bool)
	var uniqueLines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			uniqueLines = append(uniqueLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	for _, line := range uniqueLines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run file_cleaner.go <input_file> <output_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	err := removeDuplicates(inputFile, outputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully removed duplicates. Output written to %s\n", outputFile)
}