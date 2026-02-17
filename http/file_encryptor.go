
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
    saltSize      = 16
    keyIterations = 100000
    keyLength     = 32
)

func deriveKey(password string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(password))
    hash.Write(salt)
    for i := 0; i < keyIterations-1; i++ {
        hash.Write(hash.Sum(nil))
    }
    return hash.Sum(nil)[:keyLength]
}

func encryptFile(inputPath, outputPath, password string) error {
    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("salt generation failed: %w", err)
    }

    key := deriveKey(password, salt)

    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file failed: %w", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("cipher creation failed: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("GCM mode failed: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return fmt.Errorf("nonce generation failed: %w", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    encryptedData := append(salt, ciphertext...)

    if err := os.WriteFile(outputPath, encryptedData, 0644); err != nil {
        return fmt.Errorf("write output file failed: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, password string) error {
    encryptedData, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file failed: %w", err)
    }

    if len(encryptedData) < saltSize {
        return errors.New("invalid encrypted file format")
    }

    salt := encryptedData[:saltSize]
    ciphertext := encryptedData[saltSize:]

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("cipher creation failed: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("GCM mode failed: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return errors.New("decryption failed - incorrect password or corrupted file")
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("write output file failed: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) < 5 {
        fmt.Printf("Usage: %s <encrypt|decrypt> <input> <output> <password>\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }

    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    password := os.Args[4]

    var err error
    switch mode {
    case "encrypt":
        err = encryptFile(inputPath, outputPath, password)
    case "decrypt":
        err = decryptFile(inputPath, outputPath, password)
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Operation failed: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully: %s -> %s\n", inputPath, outputPath)
}