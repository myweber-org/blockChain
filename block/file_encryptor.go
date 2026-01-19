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

type Encryptor struct {
	key []byte
}

func NewEncryptor(key string) (*Encryptor, error) {
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("invalid hex key: %w", err)
	}
	if len(decodedKey) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}
	return &Encryptor{key: decodedKey}, nil
}

func (e *Encryptor) EncryptFile(inputPath, outputPath string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func (e *Encryptor) DecryptFile(inputPath, outputPath string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func generateRandomKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("generate key: %w", err)
	}
	return hex.EncodeToString(key), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_encryptor <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  genkey                     Generate random encryption key")
		fmt.Println("  encrypt <key> <input> <output>  Encrypt file")
		fmt.Println("  decrypt <key> <input> <output>  Decrypt file")
		return
	}

	switch os.Args[1] {
	case "genkey":
		key, err := generateRandomKey()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated key: %s\n", key)

	case "encrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryptor encrypt <key> <input> <output>")
			os.Exit(1)
		}
		encryptor, err := NewEncryptor(os.Args[2])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := encryptor.EncryptFile(os.Args[3], os.Args[4]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File encrypted successfully")

	case "decrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryptor decrypt <key> <input> <output>")
			os.Exit(1)
		}
		encryptor, err := NewEncryptor(os.Args[2])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := encryptor.DecryptFile(os.Args[3], os.Args[4]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}