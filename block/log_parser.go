package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	errorPattern := regexp.MustCompile(`(?i)error|exception|fail|fatal`)
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		if errorPattern.MatchString(line) {
			cleanedLine := strings.TrimSpace(line)
			fmt.Printf("Line %d: %s\n", lineNumber, cleanedLine)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}