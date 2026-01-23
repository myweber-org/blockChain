package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_cache"
	maxAge       = 7 * 24 * time.Hour
	fileModePerm = 0750
)

func main() {
	if err := cleanOldFiles(); err != nil {
		fmt.Printf("Cleanup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles() error {
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return fmt.Errorf("temp directory does not exist: %s", tempDir)
	}

	cutoffTime := time.Now().Add(-maxAge)
	var cleanupErr error

	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.ModTime().Before(cutoffTime) {
			if err := os.Remove(path); err != nil {
				cleanupErr = fmt.Errorf("failed to remove %s: %w", path, err)
				return nil
			}
			fmt.Printf("Removed: %s (modified: %v)\n", path, info.ModTime())
		}
		return nil
	})

	if cleanupErr != nil {
		return cleanupErr
	}
	return err
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

const retentionDays = 7

func main() {
	tempDir := os.TempDir()
	fmt.Printf("Cleaning temporary directory: %s\n", tempDir)

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	var removedCount int

	filepath.Walk(tempDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			err := os.Remove(path)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", path, err)
			} else {
				removedCount++
				fmt.Printf("Removed: %s\n", path)
			}
		}
		return nil
	})

	fmt.Printf("Cleanup completed. Removed %d files.\n", removedCount)
}