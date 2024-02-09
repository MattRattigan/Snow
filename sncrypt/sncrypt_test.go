package sncrypt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func testFile() (*os.File, string, string) {
	text := "I like chocolate milk"
	temp, _ := os.CreateTemp("", "testing.txt")
	path, _ := os.Getwd()
	path = filepath.Join(path, "testing.txt")
	return temp, text, path
}

func TestDeriveKey(t *testing.T) {
	t.Parallel()
	password := []byte("testPassword")
	salt := []byte("testSalt")

	key := deriveKey(password, salt)

	if len(key) != 32 {
		t.Errorf("Expected key length of 32, got %d", len(key))
	}
}

func TestEncrypt(t *testing.T) {
	t.Parallel()
	//uuid := snowUser.CreateFileUUID()
	file, text, path := testFile()
	txt, _ := file.WriteString(text)

	file.Close()
	defer os.Remove(file.Name())
	fmt.Println(path, txt)
	//fileDir := &snowUser.UserFile{UUID: uuid, FilePath: path}
	//user := &snowUser.User{
	//	Username:   "user1",
	//	Passpharse: []byte("testPassword"),
	//	FileDir:    fileDir,
	//}
	//
	//secret, err := encrypt(user)
	//if err != nil {
	//	t.Errorf("encrypt() with valid user returned an error: %v", err)
	//}
	//if len(secret) == 0 {
	//	t.Errorf("encrypt() returned empty data")
	//}
}

func TestDecrypt(t *testing.T) {

}
