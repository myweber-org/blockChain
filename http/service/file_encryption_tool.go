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

func encrypt(plaintext []byte, key []byte) (string, error) {
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

func decrypt(encodedCiphertext string, key []byte) ([]byte, error) {
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

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}

	message := "Sensitive data requiring encryption"
	encrypted, err := encrypt([]byte(message), key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Encrypted: %s\n", encrypted)

	decrypted, err := decrypt(encrypted, key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decrypted: %s\n", string(decrypted))
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
    "path/filepath"
)

const (
    saltSize   = 16
    nonceSize  = 12
    keySize    = 32
    bufferSize = 4096
)

func deriveKey(password, salt []byte) []byte {
    hash := sha256.New()
    hash.Write(password)
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, password string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("failed to generate salt: %w", err)
    }

    key := deriveKey([]byte(password), salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return fmt.Errorf("failed to generate nonce: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %w", err)
    }

    if _, err := outputFile.Write(salt); err != nil {
        return fmt.Errorf("failed to write salt: %w", err)
    }
    if _, err := outputFile.Write(nonce); err != nil {
        return fmt.Errorf("failed to write nonce: %w", err)
    }

    buffer := make([]byte, bufferSize)
    for {
        n, err := inputFile.Read(buffer)
        if err != nil && err != io.EOF {
            return fmt.Errorf("failed to read input file: %w", err)
        }

        if n > 0 {
            ciphertext := aesgcm.Seal(nil, nonce, buffer[:n], nil)
            if _, err := outputFile.Write(ciphertext); err != nil {
                return fmt.Errorf("failed to write encrypted data: %w", err)
            }
        }

        if err == io.EOF {
            break
        }
    }

    return nil
}

func decryptFile(inputPath, outputPath, password string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(inputFile, salt); err != nil {
        return fmt.Errorf("failed to read salt: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := io.ReadFull(inputFile, nonce); err != nil {
        return fmt.Errorf("failed to read nonce: %w", err)
    }

    key := deriveKey([]byte(password), salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %w", err)
    }

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    buffer := make([]byte, bufferSize+aesgcm.Overhead())
    for {
        n, err := inputFile.Read(buffer)
        if err != nil && err != io.EOF {
            return fmt.Errorf("failed to read encrypted data: %w", err)
        }

        if n > 0 {
            plaintext, err := aesgcm.Open(nil, nonce, buffer[:n], nil)
            if err != nil {
                return fmt.Errorf("failed to decrypt data: %w", err)
            }
            if _, err := outputFile.Write(plaintext); err != nil {
                return fmt.Errorf("failed to write decrypted data: %w", err)
            }
        }

        if err == io.EOF {
            break
        }
    }

    return nil
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: go run file_encryption_tool.go <encrypt|decrypt> <input> <output> <password>")
        os.Exit(1)
    }

    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    password := os.Args[4]

    switch mode {
    case "encrypt":
        if err := encryptFile(inputPath, outputPath, password); err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File encrypted successfully: %s\n", outputPath)

    case "decrypt":
        if err := decryptFile(inputPath, outputPath, password); err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File decrypted successfully: %s\n", outputPath)

    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}