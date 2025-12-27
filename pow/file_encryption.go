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

const (
	saltSize   = 16
	nonceSize  = 12
	keySize    = 32
	iterations = 100000
)

func deriveKey(passphrase, salt []byte) []byte {
	hash := sha256.New()
	hash.Write(passphrase)
	hash.Write(salt)
	
	for i := 0; i < iterations-1; i++ {
		hash.Write(hash.Sum(nil))
	}
	
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %v", err)
	}

	key := deriveKey([]byte(passphrase), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %v", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %v", err)
	}

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, inputData, nil)

	outputData := append(salt, nonce...)
	outputData = append(outputData, ciphertext...)

	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	encryptedData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %v", err)
	}

	if len(encryptedData) < saltSize+nonceSize {
		return fmt.Errorf("invalid encrypted file format")
	}

	salt := encryptedData[:saltSize]
	nonce := encryptedData[saltSize : saltSize+nonceSize]
	ciphertext := encryptedData[saltSize+nonceSize:]

	key := deriveKey([]byte(passphrase), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %v", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed: %v", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("failed to write decrypted file: %v", err)
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

	var err error
	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, passphrase)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, passphrase)
	default:
		fmt.Printf("Invalid operation: %s\n", operation)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully\n")
}