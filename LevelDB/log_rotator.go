package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

func rotateLogs(logDir, archiveDir, prefix string, maxFiles int) error {
    files, err := os.ReadDir(logDir)
    if err != nil {
        return fmt.Errorf("failed to read log directory: %w", err)
    }

    var logFiles []string
    for _, file := range files {
        if !file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
            logFiles = append(logFiles, filepath.Join(logDir, file.Name()))
        }
    }

    if len(logFiles) <= maxFiles {
        return nil
    }

    if err := os.MkdirAll(archiveDir, 0755); err != nil {
        return fmt.Errorf("failed to create archive directory: %w", err)
    }

    for i := 0; i < len(logFiles)-maxFiles; i++ {
        oldPath := logFiles[i]
        timestamp := time.Now().Format("20060102_150405")
        newName := fmt.Sprintf("%s_%s", filepath.Base(oldPath), timestamp)
        newPath := filepath.Join(archiveDir, newName)

        srcFile, err := os.Open(oldPath)
        if err != nil {
            return fmt.Errorf("failed to open source file: %w", err)
        }

        dstFile, err := os.Create(newPath)
        if err != nil {
            srcFile.Close()
            return fmt.Errorf("failed to create destination file: %w", err)
        }

        if _, err := io.Copy(dstFile, srcFile); err != nil {
            srcFile.Close()
            dstFile.Close()
            return fmt.Errorf("failed to copy file content: %w", err)
        }

        srcFile.Close()
        dstFile.Close()

        if err := os.Remove(oldPath); err != nil {
            return fmt.Errorf("failed to remove old log file: %w", err)
        }

        fmt.Printf("Rotated %s to %s\n", oldPath, newPath)
    }

    return nil
}

func main() {
    logDir := "./logs"
    archiveDir := "./archive"
    prefix := "app_"
    maxFiles := 5

    if err := rotateLogs(logDir, archiveDir, prefix, maxFiles); err != nil {
        fmt.Printf("Error rotating logs: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Log rotation completed successfully")
}