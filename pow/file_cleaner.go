package main

import (
	"bufio"
	"fmt"
	"os"
)

func removeDuplicates(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	seen := make(map[string]bool)
	var uniqueLines []string

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			uniqueLines = append(uniqueLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for _, line := range uniqueLines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to output file: %w", err)
		}
	}
	writer.Flush()

	return nil
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
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tempDir := os.TempDir()
	cutoffTime := time.Now().AddDate(0, 0, -7)

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old file: %s\n", path)
			os.RemoveAll(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error cleaning temp directory: %v\n", err)
	}
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const retentionDays = 7

func main() {
	tempDir := os.TempDir()
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old file: %s\n", path)
			os.Remove(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error cleaning temp directory: %v\n", err)
	}
}