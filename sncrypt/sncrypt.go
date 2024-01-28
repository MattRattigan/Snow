package sncrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	relativePath := "/cmd/temp.txt"
	fullPath := filepath.Join(cwd, relativePath)
	fmt.Println(fullPath)

	passphrase := "your-password-here" // Use a secure method to generate a strong passphrase

	// Read the file you want to encrypt
	data, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatal(err)
	}

	// Encrypt the file
	encryptedData := Encrypt(data, passphrase)
	err = os.WriteFile(fullPath, encryptedData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finished!!!!!")

}

func Encrypt(data []byte, passphrase string) []byte {
	// TODO: Check if the file in question

	// Ensure the passphrase is of a valid length (16, 24, 32 bytes)
	key := []byte(passphrase)
	switch len(passphrase) {
	case 16, 24, 32: // Valid AES key lengths
		// Key is already the correct length
	default:
		// Pad or trim the key to 32 bytes
		key = padOrTrim([]byte(passphrase), 32)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}
	return gcm.Seal(nonce, nonce, data, nil)
}

func padOrTrim(b []byte, size int) []byte {
	if len(b) > size {
		return b[:size]
	}
	padded := make([]byte, size)
	copy(padded, b)
	return padded
}

func Decrypt(data []byte, passphrase string) []byte {
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	return plaintext
}

func getFileLocation() {

}
