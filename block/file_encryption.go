
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "strings"

    "golang.org/x/crypto/pbkdf2"
)

const (
    saltSize      = 16
    nonceSize     = 12
    keyIterations = 100000
    keyLength     = 32
)

type EncryptedData struct {
    Ciphertext string `json:"ciphertext"`
    Salt       string `json:"salt"`
    Nonce      string `json:"nonce"`
}

func deriveKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, keyIterations, keyLength, sha256.New)
}

func Encrypt(plaintext, password string) (*EncryptedData, error) {
    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return nil, fmt.Errorf("failed to generate salt: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %w", err)
    }

    ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

    return &EncryptedData{
        Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
        Salt:       base64.StdEncoding.EncodeToString(salt),
        Nonce:      base64.StdEncoding.EncodeToString(nonce),
    }, nil
}

func Decrypt(data *EncryptedData, password string) (string, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(data.Ciphertext)
    if err != nil {
        return "", fmt.Errorf("failed to decode ciphertext: %w", err)
    }

    salt, err := base64.StdEncoding.DecodeString(data.Salt)
    if err != nil {
        return "", fmt.Errorf("failed to decode salt: %w", err)
    }

    nonce, err := base64.StdEncoding.DecodeString(data.Nonce)
    if err != nil {
        return "", fmt.Errorf("failed to decode nonce: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", errors.New("decryption failed: invalid password or corrupted data")
    }

    return string(plaintext), nil
}

func main() {
    secretMessage := "This is a confidential message for authorized personnel only."
    password := "StrongPassw0rd!2024"

    fmt.Println("Original message:", secretMessage)

    encrypted, err := Encrypt(secretMessage, password)
    if err != nil {
        fmt.Printf("Encryption failed: %v\n", err)
        return
    }

    fmt.Printf("Encrypted data:\n")
    fmt.Printf("  Ciphertext: %s...\n", strings.Split(encrypted.Ciphertext, "")[0])
    fmt.Printf("  Salt: %s\n", encrypted.Salt)
    fmt.Printf("  Nonce: %s\n", encrypted.Nonce)

    decrypted, err := Decrypt(encrypted, password)
    if err != nil {
        fmt.Printf("Decryption failed: %v\n", err)
        return
    }

    fmt.Println("Decrypted message:", decrypted)

    wrongPassword := "WrongPassword123"
    _, err = Decrypt(encrypted, wrongPassword)
    if err != nil {
        fmt.Println("Expected error with wrong password:", err)
    }
}