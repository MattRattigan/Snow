package sncrypt

import (
	"Snow/mimetype"
	"Snow/snFlags"
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
	"strings"
)

func deriveKey(password []byte, salt []byte) []byte {
	return pbkdf2.Key(password, salt, 10000, 32, sha256.New)
}

func snToFile(path string) string {
	// check for original file type extension using MIMEType,
	// creates file extension to later append to file path

	var builder strings.Builder

	// Splits the path based on . separator
	parts := strings.Split(path, ".")

	mime := mimetype.MIMEMap()
	mimeStr, _ := mime.CheckFileType(path)
	fmt.Println(mimeStr)

	if mimeStr == "" {
		builder.WriteString(parts[0])
		builder.WriteString(snFlags.CmdFlags.Ext)
		return builder.String()
	}

	mimeStr = mime.GetExtensionFromMIME(mimeStr)
	builder.WriteString(parts[0])
	builder.WriteString(mimeStr)

	return builder.String()
}

func fileToSn(path string) string {
	// Splits filepath to later append the .sn extension to the file name
	snExtension := ".sn"
	parts := strings.Split(path, ".")
	return parts[0] + snExtension
}

func WriteEncryption(user *snowUser.User) error {
	bytes, err := encrypt(user)
	if err != nil {
		return err
	}

	newPath := fileToSn(user.GetPath())

	err = os.Rename(user.GetPath(), newPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(newPath, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
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

	newPath := snToFile(user.GetPath())

	err = os.Rename(user.GetPath(), newPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(newPath, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
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
