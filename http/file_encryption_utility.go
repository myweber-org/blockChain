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

func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func encrypt(plaintext []byte, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase)
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

func decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase)
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

func encryptFile(inputPath, outputPath, passphrase string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(data, passphrase)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, encrypted, 0644)
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	decrypted, err := decrypt(data, passphrase)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, decrypted, 0644)
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <passphrase>")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	passphrase := os.Args[4]

	switch operation {
	case "encrypt":
		err := encryptFile(inputPath, outputPath, passphrase)
		if err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File encrypted successfully")
	case "decrypt":
		err := decryptFile(inputPath, outputPath, passphrase)
		if err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
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
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	saltSize      = 16
	nonceSize     = 12
	keyIterations = 100000
)

func deriveKey(password string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	for i := 0; i < keyIterations; i++ {
		hash.Write(hash.Sum(nil))
	}
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, password string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("salt generation failed: %w", err)
	}

	key := deriveKey(password, salt)

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("nonce generation failed: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation failed: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	outputData := append(salt, nonce...)
	outputData = append(outputData, ciphertext...)

	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("write file failed: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, password string) error {
	encryptedData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	if len(encryptedData) < saltSize+nonceSize {
		return errors.New("invalid encrypted file format")
	}

	salt := encryptedData[:saltSize]
	nonce := encryptedData[saltSize : saltSize+nonceSize]
	ciphertext := encryptedData[saltSize+nonceSize:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation failed: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file failed: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <password>")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	password := os.Args[4]

	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		fmt.Printf("Input path resolution failed: %v\n", err)
		os.Exit(1)
	}

	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		fmt.Printf("Output path resolution failed: %v\n", err)
		os.Exit(1)
	}

	switch operation {
	case "encrypt":
		if err := encryptFile(absInput, absOutput, password); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File encrypted successfully: %s -> %s\n", absInput, absOutput)
	case "decrypt":
		if err := decryptFile(absInput, absOutput, password); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File decrypted successfully: %s -> %s\n", absInput, absOutput)
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return os.WriteFile(outputPath, ciphertext, 0644)
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	key := deriveKey(passphrase)
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
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <encrypt|decrypt> <input> <output>\n", filepath.Base(os.Args[0]))
		fmt.Println("Passphrase will be read from ENCRYPTION_KEY environment variable")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	passphrase := os.Getenv("ENCRYPTION_KEY")
	if passphrase == "" {
		fmt.Println("Error: ENCRYPTION_KEY environment variable not set")
		os.Exit(1)
	}

	var err error
	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, passphrase)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, passphrase)
	default:
		fmt.Println("Error: operation must be 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully: %s -> %s\n", inputPath, outputPath)
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 16
	nonceSize  = 12
	keyIter    = 100000
	keyLen     = 32
)

func deriveKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, keyIter, keyLen, sha256.New)
}

func encryptFile(inputPath, outputPath, password string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("generate salt: %w", err)
	}

	key := deriveKey([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, input, nil)

	output := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	output = append(output, salt...)
	output = append(output, nonce...)
	output = append(output, ciphertext...)

	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, password string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	if len(input) < saltSize+nonceSize {
		return fmt.Errorf("invalid encrypted file format")
	}

	salt := input[:saltSize]
	nonce := input[saltSize : saltSize+nonceSize]
	ciphertext := input[saltSize+nonceSize:]

	key := deriveKey([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <password>")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	password := os.Args[4]

	var err error
	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, password)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, password)
	default:
		fmt.Printf("Invalid operation: %s\n", operation)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation %s completed successfully\n", operation)
}