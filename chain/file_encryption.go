package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "os"
)

func encryptString(key []byte, plaintext string) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptString(key []byte, cryptoText string) (string, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    if len(ciphertext) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return string(ciphertext), nil
}

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt> <text>")
        os.Exit(1)
    }

    action := os.Args[1]
    text := os.Args[2]

    key, err := generateKey()
    if err != nil {
        fmt.Printf("Error generating key: %v\n", err)
        os.Exit(1)
    }

    switch action {
    case "encrypt":
        encrypted, err := encryptString(key, text)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Key: %s\n", base64.StdEncoding.EncodeToString(key))
        fmt.Printf("Encrypted: %s\n", encrypted)

    case "decrypt":
        if len(os.Args) < 4 {
            fmt.Println("For decryption, provide key as third argument")
            os.Exit(1)
        }
        key, err := base64.StdEncoding.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Invalid key: %v\n", err)
            os.Exit(1)
        }
        decrypted, err := decryptString(key, text)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted: %s\n", decrypted)

    default:
        fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}