package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tempDir := os.TempDir()
	cutoffTime := time.Now().AddDate(0, 0, -7)
	var removedCount int

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.ModTime().Before(cutoffTime) && !info.IsDir() {
			if err := os.Remove(path); err == nil {
				removedCount++
				fmt.Printf("Removed: %s\n", path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	fmt.Printf("Cleaning completed. Removed %d files.\n", removedCount)
}package main

import (
	"bufio"
	"fmt"
	"os"
)

func removeDuplicates(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	seen := make(map[string]bool)
	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	for _, line := range lines {
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

	fmt.Printf("Duplicate removal completed. Output saved to %s\n", outputFile)
}
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	tempFilePattern = "temp_*"
	maxAgeHours     = 24
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	err := cleanTempFiles(dir)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Temporary file cleanup completed")
}

func cleanTempFiles(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, tempFilePattern))
	if err != nil {
		return err
	}

	cutoffTime := time.Now().Add(-maxAgeHours * time.Hour)
	removedCount := 0

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoffTime) {
			err := os.Remove(file)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", file, err)
			} else {
				removedCount++
				fmt.Printf("Removed: %s\n", file)
			}
		}
	}

	fmt.Printf("Total files removed: %d\n", removedCount)
	return nil
}