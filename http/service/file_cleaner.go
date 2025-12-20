package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	maxAgeDays   = 7
	fileExt      = ".tmp"
)

func main() {
	if err := cleanOldFiles(); err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles() error {
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -maxAgeDays)
	removedCount := 0

	for _, file := range files {
		if filepath.Ext(file.Name()) != fileExt {
			continue
		}

		if file.ModTime().Before(cutoffTime) {
			filePath := filepath.Join(tempDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("Warning: failed to remove %s: %v\n", filePath, err)
				continue
			}
			removedCount++
		}
	}

	fmt.Printf("Removed %d temporary files older than %d days\n", removedCount, maxAgeDays)
	return nil
}