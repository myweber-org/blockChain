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

	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		fmt.Printf("Generate key error: %v\n", err)
		os.Exit(1)
	}

	mode := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	switch mode {
	case "encrypt":
		if err := encryptFile(inputPath, outputPath, key); err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Encrypted successfully. Key: %x\n", key)
	case "decrypt":
		fmt.Print("Enter 64-character hex key: ")
		var hexKey string
		fmt.Scanln(&hexKey)
		if _, err := fmt.Sscanf(hexKey, "%x", &key); err != nil {
			fmt.Printf("Invalid key format: %v\n", err)
			os.Exit(1)
		}
		if err := decryptFile(inputPath, outputPath, key); err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Decrypted successfully")
	default:
		fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
}
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
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt|keygen>")
        return
    }

    switch os.Args[1] {
    case "encrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go encrypt <input_file> <key_hex>")
            return
        }

        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }

        key, err := hex.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Error decoding key: %v\n", err)
            return
        }

        encrypted, err := encryptData(data, key)
        if err != nil {
            fmt.Printf("Error encrypting: %v\n", err)
            return
        }

        outputFile := os.Args[2] + ".enc"
        if err := os.WriteFile(outputFile, encrypted, 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            return
        }

        fmt.Printf("File encrypted successfully: %s\n", outputFile)

    case "decrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go decrypt <input_file> <key_hex>")
            return
        }

        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }

        key, err := hex.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Error decoding key: %v\n", err)
            return
        }

        decrypted, err := decryptData(data, key)
        if err != nil {
            fmt.Printf("Error decrypting: %v\n", err)
            return
        }

        outputFile := os.Args[2] + ".dec"
        if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            return
        }

        fmt.Printf("File decrypted successfully: %s\n", outputFile)

    case "keygen":
        key, err := generateRandomKey()
        if err != nil {
            fmt.Printf("Error generating key: %v\n", err)
            return
        }

        fmt.Printf("Generated key: %s\n", hex.EncodeToString(key))

    default:
        fmt.Println("Invalid command. Use: encrypt, decrypt, or keygen")
    }
}