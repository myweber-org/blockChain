package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const retentionDays = 7

func main() {
	tempDir := os.TempDir()
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

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
			fmt.Printf("Removing old file: %s\n", path)
			os.Remove(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error cleaning temp directory: %v\n", err)
	}
}package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/myapp"
	maxAgeHours  = 168 // 7 days
)

func main() {
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	now := time.Now()
	removedCount := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileAge := now.Sub(file.ModTime())
		if fileAge.Hours() > maxAgeHours {
			filePath := filepath.Join(tempDir, file.Name())
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", file.Name(), err)
			} else {
				removedCount++
				fmt.Printf("Removed old file: %s\n", file.Name())
			}
		}
	}

	fmt.Printf("Cleanup completed. Removed %d files.\n", removedCount)
}