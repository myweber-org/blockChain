package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "strings"
)

func encrypt(plaintext, key string) (string, error) {
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return "", err
    }

    plaintextBytes := []byte(plaintext)
    blockSize := block.BlockSize()
    padding := blockSize - len(plaintextBytes)%blockSize
    padText := append(plaintextBytes, bytes.Repeat([]byte{byte(padding)}, padding)...)

    ciphertext := make([]byte, blockSize+len(padText))
    iv := ciphertext[:blockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(ciphertext[blockSize:], padText)

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(ciphertext, key string) (string, error) {
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return "", err
    }

    decoded, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    if len(decoded) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    iv := decoded[:aes.BlockSize]
    decoded = decoded[aes.BlockSize:]

    mode := cipher.NewCBCDecrypter(block, iv)
    mode.CryptBlocks(decoded, decoded)

    padding := decoded[len(decoded)-1]
    if padding > aes.BlockSize || padding == 0 {
        return "", errors.New("invalid padding")
    }

    for i := len(decoded) - int(padding); i < len(decoded); i++ {
        if decoded[i] != padding {
            return "", errors.New("invalid padding")
        }
    }

    return string(decoded[:len(decoded)-int(padding)]), nil
}

func validateKey(key string) error {
    validLengths := []int{16, 24, 32}
    for _, length := range validLengths {
        if len(key) == length {
            return nil
        }
    }
    return fmt.Errorf("key must be %s bytes long", strings.Join(strings.Split(fmt.Sprint(validLengths), " "), ", "))
}

func main() {
    key := "thisis32bitlongpassphraseimusing"
    plaintext := "Sensitive data to protect"

    if err := validateKey(key); err != nil {
        fmt.Printf("Key validation failed: %v\n", err)
        return
    }

    encrypted, err := encrypt(plaintext, key)
    if err != nil {
        fmt.Printf("Encryption failed: %v\n", err)
        return
    }
    fmt.Printf("Encrypted: %s\n", encrypted)

    decrypted, err := decrypt(encrypted, key)
    if err != nil {
        fmt.Printf("Decryption failed: %v\n", err)
        return
    }
    fmt.Printf("Decrypted: %s\n", decrypted)
}