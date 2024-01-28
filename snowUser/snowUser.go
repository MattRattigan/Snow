package snowUser

import (
	"Snow/dbSnow"
	"Snow/snFlags"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
	"log"
	"os"
)

type User struct {
	Username   string
	Passpharse []byte
	UserFile
}

type UserFile struct {
	UUID [16]byte
	data []byte
}

func CreateUser(store *dbSnow.DbStore) (*User, error) {
	cmdFlags := snFlags.Flags()

	username := SetUsername(cmdFlags["username"])
	password_hash, err := SetPasspharse()

	if err != nil {
		log.Fatal(err)
	}

	snUser, err := store.GetUser(username, password_hash)
	if err != nil {
		err = store.CreateUser(username, password_hash)
		if err != nil {
			log.Fatal(err)
		}
	}
	return snUser, err
}

func SetUsername(username string) string {
	if username == "" {
		fmt.Println("Username is required")
		os.Exit(1)
	}

	return username
}

func (u *UserFile) setFileUUID() uuid.UUID {
	return uuid.New()
}

func SetPasspharse() ([]byte, error) {
	fmt.Println("Enter Password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		return nil, err
	}

	password_hash, err := func(bytePwd []byte) ([]byte, error) {
		hashPassword, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		return hashPassword, nil
	}(passwordBytes)

	if err != nil {
		return nil, err
	}

	return password_hash, nil
}

func CheckPasswordHash([]byte, error) {

}
