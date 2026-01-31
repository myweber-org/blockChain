
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

const (
    saltSize = 16
    keySize  = 32
)

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

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

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    ciphertext = append(salt, ciphertext...)

    return os.WriteFile(outputPath, ciphertext, 0644)
}

func decryptFile(inputPath, outputPath, passphrase string) error {
    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    if len(ciphertext) < saltSize {
        return errors.New("file too short")
    }

    salt := ciphertext[:saltSize]
    ciphertext = ciphertext[saltSize:]

    key := deriveKey(passphrase, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return err
    }

    return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
    if len(os.Args) != 5 {
        fmt.Println("Usage: file_encryptor <encrypt|decrypt> <input> <output> <passphrase>")
        os.Exit(1)
    }

    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    passphrase := os.Args[4]

    switch mode {
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
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}