
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

    "golang.org/x/crypto/pbkdf2"
)

const (
    saltSize      = 16
    nonceSize     = 12
    keyIterations = 100000
    keyLength     = 32
)

type EncryptedData struct {
    Ciphertext string `json:"ciphertext"`
    Salt       string `json:"salt"`
    Nonce      string `json:"nonce"`
}

func deriveKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, keyIterations, keyLength, sha256.New)
}

func Encrypt(plaintext, password string) (*EncryptedData, error) {
    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return nil, fmt.Errorf("failed to generate salt: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %w", err)
    }

    ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

    return &EncryptedData{
        Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
        Salt:       base64.StdEncoding.EncodeToString(salt),
        Nonce:      base64.StdEncoding.EncodeToString(nonce),
    }, nil
}

func Decrypt(data *EncryptedData, password string) (string, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(data.Ciphertext)
    if err != nil {
        return "", fmt.Errorf("failed to decode ciphertext: %w", err)
    }

    salt, err := base64.StdEncoding.DecodeString(data.Salt)
    if err != nil {
        return "", fmt.Errorf("failed to decode salt: %w", err)
    }

    nonce, err := base64.StdEncoding.DecodeString(data.Nonce)
    if err != nil {
        return "", fmt.Errorf("failed to decode nonce: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", errors.New("decryption failed: invalid password or corrupted data")
    }

    return string(plaintext), nil
}

func main() {
    secretMessage := "This is a confidential message for authorized personnel only."
    password := "StrongPassw0rd!2024"

    fmt.Println("Original message:", secretMessage)

    encrypted, err := Encrypt(secretMessage, password)
    if err != nil {
        fmt.Printf("Encryption failed: %v\n", err)
        return
    }

    fmt.Printf("Encrypted data:\n")
    fmt.Printf("  Ciphertext: %s...\n", strings.Split(encrypted.Ciphertext, "")[0])
    fmt.Printf("  Salt: %s\n", encrypted.Salt)
    fmt.Printf("  Nonce: %s\n", encrypted.Nonce)

    decrypted, err := Decrypt(encrypted, password)
    if err != nil {
        fmt.Printf("Decryption failed: %v\n", err)
        return
    }

    fmt.Println("Decrypted message:", decrypted)

    wrongPassword := "WrongPassword123"
    _, err = Decrypt(encrypted, wrongPassword)
    if err != nil {
        fmt.Println("Expected error with wrong password:", err)
    }
}
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
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

func deriveKey(password string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	for i := 0; i < iterations-1; i++ {
		hash.Write(hash.Sum(nil))
	}
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, password string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(password, salt)

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	outputData := append(salt, nonce...)
	outputData = append(outputData, ciphertext...)

	return os.WriteFile(outputPath, outputData, 0644)
}

func decryptFile(inputPath, outputPath, password string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if len(data) < saltSize+nonceSize {
		return errors.New("invalid encrypted file format")
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("decryption failed - incorrect password or corrupted file")
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("Usage: %s <encrypt|decrypt> <input> <output> <password>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	password := os.Args[4]

	var err error
	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, password)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, password)
	default:
		err = errors.New("invalid operation - use 'encrypt' or 'decrypt'")
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully\n")
}