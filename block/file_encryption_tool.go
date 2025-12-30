package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "fmt"
    "io"
    "os"
)

func encryptData(plaintext []byte, key []byte) ([]byte, error) {
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

func decryptData(ciphertext []byte, key []byte) ([]byte, error) {
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

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run file_encryption_tool.go <encrypt|decrypt|generate>")
        return
    }

    switch os.Args[1] {
    case "generate":
        key, err := generateKey()
        if err != nil {
            fmt.Printf("Key generation failed: %v\n", err)
            return
        }
        fmt.Printf("Generated key: %x\n", key)

    case "encrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryption_tool.go encrypt <input_file> <key_hex>")
            return
        }

        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Printf("File read error: %v\n", err)
            return
        }

        var key []byte
        fmt.Sscanf(os.Args[3], "%x", &key)

        encrypted, err := encryptData(data, key)
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            return
        }

        err = os.WriteFile(os.Args[2]+".enc", encrypted, 0644)
        if err != nil {
            fmt.Printf("File write error: %v\n", err)
            return
        }
        fmt.Println("File encrypted successfully")

    case "decrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryption_tool.go decrypt <encrypted_file> <key_hex>")
            return
        }

        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Printf("File read error: %v\n", err)
            return
        }

        var key []byte
        fmt.Sscanf(os.Args[3], "%x", &key)

        decrypted, err := decryptData(data, key)
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            return
        }

        outputFile := os.Args[2]
        if len(outputFile) > 4 && outputFile[len(outputFile)-4:] == ".enc" {
            outputFile = outputFile[:len(outputFile)-4]
        }

        err = os.WriteFile(outputFile+".dec", decrypted, 0644)
        if err != nil {
            fmt.Printf("File write error: %v\n", err)
            return
        }
        fmt.Println("File decrypted successfully")

    default:
        fmt.Println("Invalid command. Use: encrypt, decrypt, or generate")
    }
}