package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	retentionDays = 30
	logExtension  = ".log"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_cleaner <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	err := cleanOldLogs(dir)
	if err != nil {
		fmt.Printf("Error cleaning logs: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Log cleanup completed successfully")
}

func cleanOldLogs(dir string) error {
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != logExtension {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old log file: %s (modified: %v)\n", path, info.ModTime())
			return os.Remove(path)
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
	cutoff := time.Now().AddDate(0, 0, -7)
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.ModTime().Before(cutoff) && !info.IsDir() {
			fmt.Printf("Removing old file: %s\n", path)
			os.Remove(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error cleaning temp files: %v\n", err)
	}
}
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
}