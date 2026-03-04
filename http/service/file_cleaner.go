package main

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
    "os"
    "path/filepath"
    "time"
)

func main() {
    dir := "/tmp"
    cutoff := time.Now().AddDate(0, 0, -7)

    filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

	fmt.Printf("Cleanup completed. Removed %d files.\n", removedCount)
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	dir := "./tmp"
	cleanupAge := 24 * time.Hour

	err := cleanOldFiles(dir, cleanupAge)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dir string, maxAge time.Duration) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if time.Since(info.ModTime()) > maxAge {
			fmt.Printf("Removing old file: %s\n", path)
			return os.Remove(path)
		}
		return nil
	})
}