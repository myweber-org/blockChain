
package main

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

const saltSize = 16

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("generate salt: %w", err)
    }

    key := deriveKey(passphrase, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return fmt.Errorf("generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
    outputData := append(salt, nonce...)
    outputData = append(outputData, ciphertext...)

    return os.WriteFile(outputPath, outputData, 0644)
}

func decryptFile(inputPath, outputPath, passphrase string) error {
    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    if len(ciphertext) < saltSize {
        return errors.New("invalid ciphertext length")
    }

    salt := ciphertext[:saltSize]
    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < saltSize+nonceSize {
        return errors.New("invalid ciphertext length")
    }

    nonce := ciphertext[saltSize : saltSize+nonceSize]
    actualCiphertext := ciphertext[saltSize+nonceSize:]

    plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
    if err != nil {
        return fmt.Errorf("decrypt: %w", err)
    }

    return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: go run file_encryption_tool.go <encrypt|decrypt> <input> <output> <passphrase>")
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
package main

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

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptData(plaintext []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := deriveKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	result := append(salt, nonce...)
	result = append(result, ciphertext...)
	return result, nil
}

func decryptData(ciphertext []byte, passphrase string) ([]byte, error) {
	if len(ciphertext) < 48 {
		return nil, errors.New("ciphertext too short")
	}

	salt := ciphertext[:16]
	nonce := ciphertext[16:32]
	actualCiphertext := ciphertext[32:]

	key := deriveKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run file_encryption_tool.go <encrypt|decrypt> <input_file> <output_file>")
		fmt.Println("Passphrase will be read from environment variable ENCRYPTION_KEY")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	passphrase := os.Getenv("ENCRYPTION_KEY")
	if passphrase == "" {
		fmt.Println("Error: ENCRYPTION_KEY environment variable not set")
		os.Exit(1)
	}

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	var outputData []byte
	switch operation {
	case "encrypt":
		outputData, err = encryptData(inputData, passphrase)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			os.Exit(1)
		}
	case "decrypt":
		outputData, err = decryptData(inputData, passphrase)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully. Output written to %s\n", outputFile)
}