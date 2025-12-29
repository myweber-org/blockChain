package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func encryptData(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptData(encodedCiphertext string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		fmt.Println("Key generation failed:", err)
		return
	}

	secretMessage := "This is a confidential message"
	encrypted, err := encryptData([]byte(secretMessage), key)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}

	fmt.Println("Encrypted:", encrypted)

	decrypted, err := decryptData(encrypted, key)
	if err != nil {
		fmt.Println("Decryption failed:", err)
		return
	}

	fmt.Println("Decrypted:", string(decrypted))
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
		return
	}

	key, err := generateKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		return
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
		return
	}

	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
		return
	}

	fmt.Printf("Key (hex): %x\n", key)
	fmt.Printf("Operation completed successfully\n")
}