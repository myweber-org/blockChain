
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func calculateFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func verifyFileIntegrity(filePath, expectedChecksum string) (bool, error) {
	actualChecksum, err := calculateFileChecksum(filePath)
	if err != nil {
		return false, err
	}

	return actualChecksum == expectedChecksum, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run file_checksum_verifier.go <file_path> <expected_checksum>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	expectedChecksum := os.Args[2]

	isValid, err := verifyFileIntegrity(filePath, expectedChecksum)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if isValid {
		fmt.Println("File integrity verified: checksum matches")
	} else {
		fmt.Println("WARNING: File integrity check failed - checksum mismatch")
		os.Exit(1)
	}
}
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func calculateFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

func verifyFileIntegrity(filePath, expectedChecksum string) (bool, error) {
	actualChecksum, err := calculateFileChecksum(filePath)
	if err != nil {
		return false, err
	}

	return actualChecksum == expectedChecksum, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run file_checksum_verifier.go <file_path> <expected_checksum>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	expectedChecksum := os.Args[2]

	isValid, err := verifyFileIntegrity(filePath, expectedChecksum)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if isValid {
		fmt.Println("File integrity verified: checksum matches")
	} else {
		fmt.Println("WARNING: File integrity check failed - checksum mismatch")
		os.Exit(1)
	}
}
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
)

func calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func verifyChecksum(filePath, expectedChecksum string) (bool, error) {
	actualChecksum, err := calculateChecksum(filePath)
	if err != nil {
		return false, err
	}
	return actualChecksum == expectedChecksum, nil
}

func main() {
	action := flag.String("action", "calculate", "Action to perform: calculate or verify")
	file := flag.String("file", "", "Path to the file")
	checksum := flag.String("checksum", "", "Expected checksum for verification")
	flag.Parse()

	if *file == "" {
		fmt.Println("Error: file path is required")
		flag.Usage()
		os.Exit(1)
	}

	switch *action {
	case "calculate":
		result, err := calculateChecksum(*file)
		if err != nil {
			fmt.Printf("Error calculating checksum: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("SHA256 checksum: %s\n", result)

	case "verify":
		if *checksum == "" {
			fmt.Println("Error: checksum is required for verification")
			flag.Usage()
			os.Exit(1)
		}
		match, err := verifyChecksum(*file, *checksum)
		if err != nil {
			fmt.Printf("Error verifying checksum: %v\n", err)
			os.Exit(1)
		}
		if match {
			fmt.Println("Checksum verification passed")
		} else {
			fmt.Println("Checksum verification failed")
			os.Exit(1)
		}

	default:
		fmt.Printf("Error: unknown action '%s'\n", *action)
		flag.Usage()
		os.Exit(1)
	}
}