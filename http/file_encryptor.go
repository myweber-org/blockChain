
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

func encryptFile(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation error: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation error: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation error: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption error: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("key generation error: %w", err)
	}
	return key, nil
}

func main() {
	key, err := generateRandomKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		return
	}

	fmt.Printf("Generated key: %x\n", key)

	testData := []byte("This is a secret message for encryption testing.")
	if err := os.WriteFile("test_input.txt", testData, 0644); err != nil {
		fmt.Printf("Test file creation failed: %v\n", err)
		return
	}
	defer os.Remove("test_input.txt")
	defer os.Remove("test_encrypted.bin")
	defer os.Remove("test_decrypted.txt")

	if err := encryptFile("test_input.txt", "test_encrypted.bin", key); err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}

	if err := decryptFile("test_encrypted.bin", "test_decrypted.txt", key); err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}

	decryptedData, err := os.ReadFile("test_decrypted.txt")
	if err != nil {
		fmt.Printf("Read decrypted file failed: %v\n", err)
		return
	}

	if string(decryptedData) == string(testData) {
		fmt.Println("Encryption and decryption successful")
	} else {
		fmt.Println("Encryption/decryption verification failed")
	}
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

const keySize = 32

func generateKey() ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func encryptFile(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return os.WriteFile(outputPath, ciphertext, 0644)
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output>")
		os.Exit(1)
	}

	key, err := generateKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		os.Exit(1)
	}

	mode := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	switch mode {
	case "encrypt":
		err = encryptFile(inputFile, outputFile, key)
	case "decrypt":
		err = decryptFile(inputFile, outputFile, key)
	default:
		fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Key (hex): %x\n", key)
	fmt.Printf("Operation completed successfully\n")
}package main

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
		return fmt.Errorf("read file failed: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation failed: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation failed: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file failed: %v", err)
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
		return fmt.Errorf("read file failed: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation failed: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed: %v", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output> <hex_key>")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	keyHex := os.Args[4]

	var err error
	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, keyHex)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, keyHex)
	default:
		fmt.Println("Operation must be 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation %s completed successfully\n", operation)
}