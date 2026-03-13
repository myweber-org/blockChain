
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

func encryptData(plaintext, password []byte) (string, error) {
    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

    result := hex.EncodeToString(salt) + ":" + hex.EncodeToString(nonce) + ":" + hex.EncodeToString(ciphertext)
    return result, nil
}

func decryptData(encrypted string, password []byte) ([]byte, error) {
    parts := strings.Split(encrypted, ":")
    if len(parts) != 3 {
        return nil, errors.New("invalid encrypted data format")
    }

    salt, err := hex.DecodeString(parts[0])
    if err != nil {
        return nil, err
    }

    nonce, err := hex.DecodeString(parts[1])
    if err != nil {
        return nil, err
    }

    ciphertext, err := hex.DecodeString(parts[2])
    if err != nil {
        return nil, err
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <output_file>")
        fmt.Println("Password will be read from ENCRYPTION_PASSWORD environment variable")
        return
    }

    action := os.Args[1]
    inputFile := os.Args[2]
    outputFile := os.Args[3]

    password := os.Getenv("ENCRYPTION_PASSWORD")
    if password == "" {
        fmt.Println("Error: ENCRYPTION_PASSWORD environment variable not set")
        os.Exit(1)
    }

    inputData, err := os.ReadFile(inputFile)
    if err != nil {
        fmt.Printf("Error reading input file: %v\n", err)
        os.Exit(1)
    }

    switch action {
    case "encrypt":
        encrypted, err := encryptData(inputData, []byte(password))
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        err = os.WriteFile(outputFile, []byte(encrypted), 0644)
        if err != nil {
            fmt.Printf("Error writing output file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File encrypted successfully: %s\n", outputFile)

    case "decrypt":
        decrypted, err := decryptData(string(inputData), []byte(password))
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        err = os.WriteFile(outputFile, decrypted, 0644)
        if err != nil {
            fmt.Printf("Error writing output file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File decrypted successfully: %s\n", outputFile)

    default:
        fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}