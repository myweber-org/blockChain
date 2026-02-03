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
)

const saltSize = 16

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	plaintext, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if _, err := outputFile.Write(salt); err != nil {
		return err
	}
	if _, err := outputFile.Write(nonce); err != nil {
		return err
	}
	if _, err := outputFile.Write(ciphertext); err != nil {
		return err
	}

	return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(inputFile, salt); err != nil {
		return err
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(inputFile, nonce); err != nil {
		return err
	}

	ciphertext, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("decryption failed: invalid passphrase or corrupted file")
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if _, err := outputFile.Write(plaintext); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <passphrase>")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	passphrase := os.Args[4]

	switch operation {
	case "encrypt":
		if err := encryptFile(inputPath, outputPath, passphrase); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File encrypted successfully")
	case "decrypt":
		if err := decryptFile(inputPath, outputPath, passphrase); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully")
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
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

func deriveKey(passphrase string) []byte {
    hash := sha256.Sum256([]byte(passphrase))
    return hash[:]
}

func encryptData(plaintext []byte, passphrase string) ([]byte, error) {
    key := deriveKey(passphrase)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

func decryptData(ciphertext []byte, passphrase string) ([]byte, error) {
    key := deriveKey(passphrase)
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
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <output_file>")
        fmt.Println("Passphrase will be read from environment variable ENCRYPTION_KEY")
        os.Exit(1)
    }

    mode := os.Args[1]
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
    switch mode {
    case "encrypt":
        outputData, err = encryptData(inputData, passphrase)
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Encrypted %d bytes to %s\n", len(outputData), outputFile)
    case "decrypt":
        outputData, err = decryptData(inputData, passphrase)
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted %d bytes to %s\n", len(outputData), outputFile)
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }

    err = os.WriteFile(outputFile, outputData, 0644)
    if err != nil {
        fmt.Printf("Error writing output file: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Operation completed successfully")
}