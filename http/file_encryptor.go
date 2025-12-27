package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

func encryptFile(inputPath, outputPath, keyHex string) error {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return fmt.Errorf("invalid key: %v", err)
	}
	if len(key) != 32 {
		return errors.New("key must be 32 bytes for AES-256")
	}

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("gcm creation error: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation error: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file error: %v", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, keyHex string) error {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return fmt.Errorf("invalid key: %v", err)
	}
	if len(key) != 32 {
		return errors.New("key must be 32 bytes for AES-256")
	}

	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("gcm creation error: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption error: %v", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file error: %v", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output> <hex_key>")
		os.Exit(1)
	}

	action := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	keyHex := os.Args[4]

	switch action {
	case "encrypt":
		if err := encryptFile(inputPath, outputPath, keyHex); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Encryption successful")
	case "decrypt":
		if err := decryptFile(inputPath, outputPath, keyHex); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Decryption successful")
	default:
		fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
}