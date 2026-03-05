
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	retentionDays = 7
)

func main() {
	if err := cleanOldFiles(); err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles() error {
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	return filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
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
				return fmt.Errorf("failed to remove %s: %w", path, err)
			}
			fmt.Printf("Removed: %s\n", path)
		}

		return nil
	})
}
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	maxAgeHours  = 168
)

func main() {
	err := cleanOldFiles(tempDir, maxAgeHours)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Temporary files cleaned successfully")
}

func cleanOldFiles(dirPath string, maxAgeHours int) error {
	cutoffTime := time.Now().Add(-time.Duration(maxAgeHours) * time.Hour)

	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
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
			fmt.Printf("Removing old file: %s (modified: %v)\n", path, info.ModTime())
			return os.Remove(path)
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
    tempDir := os.TempDir()
    cutoffTime := time.Now().AddDate(0, 0, -7)

    filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }

        if info.IsDir() {
            return nil
        }

        if info.ModTime().Before(cutoffTime) {
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
			err := os.Remove(path)
			if err != nil {
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
}