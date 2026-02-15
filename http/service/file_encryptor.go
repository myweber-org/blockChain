
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
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

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func main() {
	key, err := generateRandomKey()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		os.Exit(1)
	}

	originalText := []byte("Sensitive data requiring encryption")
	fmt.Printf("Original text: %s\n", originalText)

	encrypted, err := encryptData(originalText, key)
	if err != nil {
		fmt.Printf("Encryption error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Encrypted (hex): %s\n", hex.EncodeToString(encrypted))

	decrypted, err := decryptData(encrypted, key)
	if err != nil {
		fmt.Printf("Decryption error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Decrypted text: %s\n", decrypted)

	if string(originalText) == string(decrypted) {
		fmt.Println("Encryption/decryption successful")
	} else {
		fmt.Println("Encryption/decryption failed")
		os.Exit(1)
	}
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

func encryptFile(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("generate iv: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return errors.New("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func main() {
	key := []byte("examplekey123456789012345678901234")
	
	err := encryptFile("plain.txt", "encrypted.bin", key)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}
	fmt.Println("File encrypted successfully")

	err = decryptFile("encrypted.bin", "decrypted.txt", key)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Println("File decrypted successfully")
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

func encryptFile(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypt data: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output>")
		fmt.Println("Key must be 32 bytes for AES-256")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		fmt.Printf("Generate key error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated key (hex): %x\n", key)

	switch operation {
	case "encrypt":
		if err := encryptFile(inputPath, outputPath, key); err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File encrypted successfully")
	case "decrypt":
		if err := decryptFile(inputPath, outputPath, key); err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully")
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

func encryptFile(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode error: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation error: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode error: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption error: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func main() {
	key := []byte("32-byte-long-key-here-123456789012")
	
	err := encryptFile("plain.txt", "encrypted.bin", key)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}
	fmt.Println("File encrypted successfully")

	err = decryptFile("encrypted.bin", "decrypted.txt", key)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Println("File decrypted successfully")
}
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

const saltSize = 16

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("generate salt: %w", err)
    }

    key := deriveKey(passphrase, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return fmt.Errorf("generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
    encryptedData := append(salt, nonce...)
    encryptedData = append(encryptedData, ciphertext...)

    if err := os.WriteFile(outputPath, encryptedData, 0644); err != nil {
        return fmt.Errorf("write output file: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
    encryptedData, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    if len(encryptedData) < saltSize {
        return errors.New("invalid encrypted file format")
    }

    salt := encryptedData[:saltSize]
    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(encryptedData) < saltSize+nonceSize {
        return errors.New("invalid encrypted file format")
    }

    nonce := encryptedData[saltSize : saltSize+nonceSize]
    ciphertext := encryptedData[saltSize+nonceSize:]

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return errors.New("decryption failed: invalid passphrase or corrupted data")
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("write output file: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) < 5 {
        fmt.Printf("Usage: %s <encrypt|decrypt> <input> <output> <passphrase>\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }

    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    passphrase := os.Args[4]

    var err error
    switch mode {
    case "encrypt":
        err = encryptFile(inputPath, outputPath, passphrase)
    case "decrypt":
        err = decryptFile(inputPath, outputPath, passphrase)
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully: %s -> %s\n", inputPath, outputPath)
}