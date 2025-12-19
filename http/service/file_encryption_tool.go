package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

func encryptData(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptData(encrypted string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_encryption_tool.go <command>")
		fmt.Println("Commands: generate-key, encrypt <text>, decrypt <encrypted>")
		return
	}

	switch os.Args[1] {
	case "generate-key":
		key, err := generateKey()
		if err != nil {
			fmt.Printf("Error generating key: %v\n", err)
			return
		}
		fmt.Printf("Generated key: %s\n", base64.StdEncoding.EncodeToString(key))

	case "encrypt":
		if len(os.Args) < 4 {
			fmt.Println("Usage: encrypt <key> <text>")
			return
		}
		key, err := base64.StdEncoding.DecodeString(os.Args[2])
		if err != nil {
			fmt.Printf("Invalid key: %v\n", err)
			return
		}
		encrypted, err := encryptData([]byte(os.Args[3]), key)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			return
		}
		fmt.Printf("Encrypted: %s\n", encrypted)

	case "decrypt":
		if len(os.Args) < 4 {
			fmt.Println("Usage: decrypt <key> <encrypted>")
			return
		}
		key, err := base64.StdEncoding.DecodeString(os.Args[2])
		if err != nil {
			fmt.Printf("Invalid key: %v\n", err)
			return
		}
		decrypted, err := decryptData(os.Args[3], key)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			return
		}
		fmt.Printf("Decrypted: %s\n", decrypted)

	default:
		fmt.Println("Unknown command")
	}
}