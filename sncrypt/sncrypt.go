package sncrypt

import (
	"Snow/snowUser"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"log"
	"os"
)

func deriveKey(password []byte, salt []byte) []byte {
	return pbkdf2.Key(password, salt, 10000, 32, sha256.New)
}

func WriteEncryption(user *snowUser.User) error {
	bytes, err := encrypt(user)

	if err != nil {
		return err
	}

	err = os.WriteFile(user.GetPath(), bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finished!!!!!")
	return nil
}

func encrypt(user *snowUser.User) ([]byte, error) {
	var salt []byte

	saltyUUID, err := user.GetUUID().ToBytesInplace(salt)

	if err != nil {
		return nil, err
	}

	key := deriveKey(user.Passpharse, saltyUUID)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, user.GetData(), nil), nil
}

func WriteDecryption(user *snowUser.User) error {
	bytes, err := decrypt(user)

	if err != nil {
		return err
	}

	err = os.WriteFile(user.GetPath(), bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finished!!!!!")
	return nil
}

func decrypt(user *snowUser.User) ([]byte, error) {
	var salt []byte

	saltyUUID, err := user.GetUUID().ToBytesInplace(salt)

	if err != nil {
		return nil, err
	}

	key := deriveKey(user.Passpharse, saltyUUID)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	data := user.GetData()
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
