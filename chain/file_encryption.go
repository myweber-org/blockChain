
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
    "strings"

    "golang.org/x/crypto/pbkdf2"
)

const (
    saltSize      = 16
    nonceSize     = 12
    keyIterations = 100000
    keyLength     = 32
)

func deriveKey(password, salt []byte) []byte {
    return pbkdf2.Key(password, salt, keyIterations, keyLength, sha256.New)
}

func encryptFile(inputPath, outputPath, password string) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return fmt.Errorf("generate salt: %w", err)
    }

    key := deriveKey([]byte(password), salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return fmt.Errorf("generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

    outputData := append(salt, nonce...)
    outputData = append(outputData, ciphertext...)

    return os.WriteFile(outputPath, outputData, 0644)
}

func decryptFile(inputPath, outputPath, password string) error {
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    if len(data) < saltSize+nonceSize {
        return errors.New("invalid encrypted file format")
    }

    salt := data[:saltSize]
    nonce := data[saltSize : saltSize+nonceSize]
    ciphertext := data[saltSize+nonceSize:]

    key := deriveKey([]byte(password), salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM: %w", err)
    }

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("decrypt data: %w", err)
    }

    return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt> <input> <output>")
        fmt.Println("Example: go run file_encryption.go encrypt secret.txt secret.enc")
        os.Exit(1)
    }

    mode := strings.ToLower(os.Args[1])
    inputPath := os.Args[2]
    outputPath := os.Args[3]

    fmt.Print("Enter password: ")
    var password string
    fmt.Scanln(&password)

    switch mode {
    case "encrypt":
        if err := encryptFile(inputPath, outputPath, password); err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File encrypted successfully")
    case "decrypt":
        if err := decryptFile(inputPath, outputPath, password); err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File decrypted successfully")
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}