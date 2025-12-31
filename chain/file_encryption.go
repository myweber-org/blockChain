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

const keySize = 32

func generateKey() ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

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

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func saveToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func readFromFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func main() {
	key, err := generateKey()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}

	originalText := []byte("This is a secret message that needs encryption")
	fmt.Printf("Original text: %s\n", originalText)

	encrypted, err := encryptData(originalText, key)
	if err != nil {
		fmt.Printf("Encryption error: %v\n", err)
		return
	}

	err = saveToFile("encrypted.dat", encrypted)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		return
	}

	savedData, err := readFromFile("encrypted.dat")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	decrypted, err := decryptData(savedData, key)
	if err != nil {
		fmt.Printf("Decryption error: %v\n", err)
		return
	}

	fmt.Printf("Decrypted text: %s\n", decrypted)
	fmt.Printf("Encryption/decryption successful\n")
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func encrypt(plaintext []byte, key []byte) (string, error) {
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
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encodedCiphertext string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, err
	}

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
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	secret := "Sensitive data requiring protection"
	encrypted, err := encrypt([]byte(secret), key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Encrypted: %s\n", encrypted)

	decrypted, err := decrypt(encrypted, key)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decrypted: %s\n", string(decrypted))
}
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

func encryptData(key []byte, plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return ciphertext, nil
}

func decryptData(key []byte, ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    if len(ciphertext) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return ciphertext, nil
}

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, key); err != nil {
        return nil, err
    }
    return key, nil
}

func processFile(inputPath string, outputPath string, encrypt bool, key []byte) error {
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    var processedData []byte
    if encrypt {
        processedData, err = encryptData(key, data)
        if err != nil {
            return err
        }
    } else {
        processedData, err = decryptData(key, data)
        if err != nil {
            return err
        }
    }

    return os.WriteFile(outputPath, processedData, 0644)
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt> <input_file> <output_file>")
        return
    }

    mode := os.Args[1]
    inputFile := os.Args[2]
    outputFile := os.Args[3]

    key, err := generateKey()
    if err != nil {
        fmt.Printf("Key generation failed: %v\n", err)
        return
    }

    keyB64 := base64.StdEncoding.EncodeToString(key)
    fmt.Printf("Generated key (base64): %s\n", keyB64)

    encrypt := mode == "encrypt"
    err = processFile(inputFile, outputFile, encrypt, key)
    if err != nil {
        fmt.Printf("Operation failed: %v\n", err)
        return
    }

    action := "Encrypted"
    if !encrypt {
        action = "Decrypted"
    }
    fmt.Printf("%s %s to %s successfully\n", action, inputFile, outputFile)
}