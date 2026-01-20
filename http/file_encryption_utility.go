
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
    "os"
)

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encrypt(plaintext []byte, passphrase string) (string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    key := deriveKey(passphrase, salt)
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
    ciphertext = append(salt, ciphertext...)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encodedCiphertext string, passphrase string) ([]byte, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
    if err != nil {
        return nil, err
    }

    if len(ciphertext) < 16 {
        return nil, errors.New("ciphertext too short")
    }

    salt := ciphertext[:16]
    ciphertext = ciphertext[16:]

    key := deriveKey(passphrase, salt)
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
    return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
    secretMessage := "Sensitive data requiring protection"
    password := "securePass123!"

    encrypted, err := encrypt([]byte(secretMessage), password)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Encryption failed: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Encrypted: %s\n", encrypted)

    decrypted, err := decrypt(encrypted, password)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Decryption failed: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Decrypted: %s\n", string(decrypted))
}