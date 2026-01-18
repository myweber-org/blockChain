
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

func encryptFile(inputPath, outputPath string, key []byte) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file error: %v", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("cipher creation error: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("GCM creation error: %v", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return fmt.Errorf("nonce generation error: %v", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

    encoded := base64.StdEncoding.EncodeToString(ciphertext)
    return os.WriteFile(outputPath, []byte(encoded), 0644)
}

func decryptFile(inputPath, outputPath string, key []byte) error {
    encoded, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file error: %v", err)
    }

    ciphertext, err := base64.StdEncoding.DecodeString(string(encoded))
    if err != nil {
        return fmt.Errorf("base64 decode error: %v", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("cipher creation error: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("GCM creation error: %v", err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("decryption error: %v", err)
    }

    return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
    key := []byte("32-byte-long-key-here-1234567890")
    
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt> <input> <output>")
        return
    }

    mode := os.Args[1]
    inputFile := os.Args[2]
    outputFile := os.Args[3]

    switch mode {
    case "encrypt":
        if err := encryptFile(inputFile, outputFile, key); err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
        } else {
            fmt.Println("File encrypted successfully")
        }
    case "decrypt":
        if err := decryptFile(inputFile, outputFile, key); err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
        } else {
            fmt.Println("File decrypted successfully")
        }
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
    }
}