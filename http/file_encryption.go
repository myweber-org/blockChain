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
	for i := 0; i < iterations; i++ {
		hash.Write(hash.Sum(nil))
	}
	return hash.Sum(nil)[:keySize]
}

func encryptFile(inputPath, outputPath, password string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("salt generation failed: %w", err)
	}

	key := deriveKey(password, salt)

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("nonce generation failed: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation failed: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	outputData := append(salt, nonce...)
	outputData = append(outputData, ciphertext...)

	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("write file failed: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, password string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	if len(ciphertext) < saltSize+nonceSize {
		return errors.New("file too short")
	}

	salt := ciphertext[:saltSize]
	nonce := ciphertext[saltSize : saltSize+nonceSize]
	encryptedData := ciphertext[saltSize+nonceSize:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation failed: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return errors.New("decryption failed - wrong password or corrupted file")
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file failed: %w", err)
	}

	return nil
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
		fmt.Println("Operation must be 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Operation completed successfully")
}package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "errors"
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

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return os.WriteFile(outputPath, ciphertext, 0644)
}

func decryptFile(inputPath, outputPath, passphrase string) error {
    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    if len(ciphertext) < aes.BlockSize {
        return errors.New("ciphertext too short")
    }

    key := deriveKey(passphrase)
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return os.WriteFile(outputPath, ciphertext, 0644)
}

func main() {
    if len(os.Args) != 5 {
        fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt> <input> <output> <passphrase>")
        os.Exit(1)
    }

    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    passphrase := os.Args[4]

    var err error
    switch mode {
    case "encrypt":
        err = encryptFile(inputPath, outputPath, passphrase)
    case "decrypt":
        err = decryptFile(inputPath, outputPath, passphrase)
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully: %s -> %s\n", inputPath, outputPath)
}