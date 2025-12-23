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
)

func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	key := deriveKey(passphrase)
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

func decryptFile(inputPath, outputPath, passphrase string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	key := deriveKey(passphrase)
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
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <passphrase>")
		os.Exit(1)
	}

	action := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	passphrase := os.Args[4]

	var err error
	switch action {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, passphrase)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, passphrase)
	default:
		fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully. Output saved to %s\n", outputPath)
}
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "strings"
)

const (
    saltSize        = 16
    nonceSize       = 12
    keyIterations   = 100000
    keyLength       = 32
)

func deriveKey(password, salt []byte) []byte {
    hash := sha256.New()
    hash.Write(password)
    hash.Write(salt)
    for i := 0; i < keyIterations-1; i++ {
        hash.Write(hash.Sum(nil))
    }
    return hash.Sum(nil)[:keyLength]
}

func encrypt(plaintext, password string) (string, error) {
    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return "", fmt.Errorf("salt generation failed: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", fmt.Errorf("nonce generation failed: %w", err)
    }

    key := deriveKey([]byte(password), salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("cipher creation failed: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("GCM mode initialization failed: %w", err)
    }

    ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
    combined := append(salt, nonce...)
    combined = append(combined, ciphertext...)

    return base64.RawURLEncoding.EncodeToString(combined), nil
}

func decrypt(encrypted, password string) (string, error) {
    data, err := base64.RawURLEncoding.DecodeString(encrypted)
    if err != nil {
        return "", fmt.Errorf("base64 decoding failed: %w", err)
    }

    if len(data) < saltSize+nonceSize {
        return "", errors.New("encrypted data too short")
    }

    salt := data[:saltSize]
    nonce := data[saltSize : saltSize+nonceSize]
    ciphertext := data[saltSize+nonceSize:]

    key := deriveKey([]byte(password), salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("cipher creation failed: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("GCM mode initialization failed: %w", err)
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", fmt.Errorf("decryption failed: %w", err)
    }

    return string(plaintext), nil
}

func main() {
    secretMessage := "confidential data for secure transmission"
    password := "strongPassw0rd!2024"

    fmt.Println("Original message:", secretMessage)

    encrypted, err := encrypt(secretMessage, password)
    if err != nil {
        fmt.Println("Encryption error:", err)
        return
    }
    fmt.Println("Encrypted:", encrypted)
    fmt.Println("Encrypted length:", len(encrypted))

    if strings.Contains(encrypted, ".") || strings.Contains(encrypted, "+") || strings.Contains(encrypted, "/") {
        fmt.Println("Warning: Base64 contains URL-unsafe characters")
    }

    decrypted, err := decrypt(encrypted, password)
    if err != nil {
        fmt.Println("Decryption error:", err)
        return
    }
    fmt.Println("Decrypted:", decrypted)

    if secretMessage == decrypted {
        fmt.Println("SUCCESS: Message integrity verified")
    } else {
        fmt.Println("ERROR: Message corrupted during encryption/decryption")
    }
}