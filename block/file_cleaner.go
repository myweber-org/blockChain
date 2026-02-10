package main

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

	fmt.Printf("Duplicate removal complete. Output written to %s\n", outputFile)
}package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const tempDir = "/tmp/myapp"
const retentionDays = 7

func main() {
	err := cleanOldFiles(tempDir, retentionDays)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, days int) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	cutoffTime := time.Now().AddDate(0, 0, -days)
	for _, file := range files {
		if file.ModTime().Before(cutoffTime) {
			fullPath := filepath.Join(dirPath, file.Name())
			err := os.Remove(fullPath)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", fullPath, err)
			} else {
				fmt.Printf("Removed: %s\n", fullPath)
			}
		}
	}
	return nil
}package main

import (
    "os"
    "path/filepath"
    "time"
)

func main() {
    tempDir := os.TempDir()
    cutoff := time.Now().AddDate(0, 0, -7)

    filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if info.IsDir() {
            return nil
        }
        if info.ModTime().Before(cutoff) {
            os.Remove(path)
        }
        return nil
    })
}package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	daysToKeep = 7
	tempDir    = "/tmp"
)

func main() {
	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)
	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.ModTime().Before(cutoffTime) {
			if err := os.Remove(path); err != nil {
				fmt.Printf("Failed to remove %s: %v\n", path, err)
			} else {
				fmt.Printf("Removed old file: %s\n", path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	retentionDays = 7
)

func main() {
	err := cleanOldFiles(tempDir, retentionDays)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old file: %s (modified: %v)\n", path, info.ModTime())
			return os.Remove(path)
		}

		return nil
	})
}