package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const maxLogSize = 10 * 1024 * 1024
const maxLogFiles = 5

func rotateLogFile(logPath string) error {
	info, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if info.Size() < maxLogSize {
		return nil
	}

	dir := filepath.Dir(logPath)
	base := filepath.Base(logPath)
	ext := filepath.Ext(logPath)
	name := strings.TrimSuffix(base, ext)

	for i := maxLogFiles - 1; i > 0; i-- {
		oldName := name + "." + strconv.Itoa(i) + ext
		newName := name + "." + strconv.Itoa(i+1) + ext
		oldPath := filepath.Join(dir, oldName)
		newPath := filepath.Join(dir, newName)

		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	firstRotated := filepath.Join(dir, name+".1"+ext)
	os.Rename(logPath, firstRotated)

	return nil
}

func setupLogger() (*log.Logger, error) {
	logPath := "application.log"
	if err := rotateLogFile(logPath); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	logger := log.New(file, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger, nil
}

func main() {
	logger, err := setupLogger()
	if err != nil {
		log.Fatal(err)
	}

	logger.Println("Application started")
	logger.Println("Performing system check...")
	logger.Println("System check completed successfully")
}