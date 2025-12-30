package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const saltSize = 16

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("salt generation failed: %v", err)
	}

	key := deriveKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %v", err)
	}

	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("input file open failed: %v", err)
	}
	defer input.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("output file creation failed: %v", err)
	}
	defer output.Close()

	if _, err := output.Write(salt); err != nil {
		return fmt.Errorf("salt write failed: %v", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return fmt.Errorf("iv generation failed: %v", err)
	}

	if _, err := output.Write(iv); err != nil {
		return fmt.Errorf("iv write failed: %v", err)
	}

	stream := cipher.NewCTR(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: output}

	if _, err := io.Copy(writer, input); err != nil {
		return fmt.Errorf("encryption copy failed: %v", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("input file open failed: %v", err)
	}
	defer input.Close()

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(input, salt); err != nil {
		return fmt.Errorf("salt read failed: %v", err)
	}

	key := deriveKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %v", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(input, iv); err != nil {
		return fmt.Errorf("iv read failed: %v", err)
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("output file creation failed: %v", err)
	}
	defer output.Close()

	stream := cipher.NewCTR(block, iv)
	reader := &cipher.StreamReader{S: stream, R: input}

	if _, err := io.Copy(output, reader); err != nil {
		return fmt.Errorf("decryption copy failed: %v", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("Usage: %s <encrypt|decrypt> <input> <output> <passphrase>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	passphrase := os.Args[4]

	switch operation {
	case "encrypt":
		if err := encryptFile(inputPath, outputPath, passphrase); err != nil {
			fmt.Fprintf(os.Stderr, "Encryption error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File encrypted successfully: %s -> %s\n", inputPath, outputPath)
	case "decrypt":
		if err := decryptFile(inputPath, outputPath, passphrase); err != nil {
			fmt.Fprintf(os.Stderr, "Decryption error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File decrypted successfully: %s -> %s\n", inputPath, outputPath)
	default:
		fmt.Fprintf(os.Stderr, "Invalid operation: %s\n", operation)
		os.Exit(1)
	}
}