package dbSnow

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func isPassword(pwd, storedPwd []byte) error {
	if err := bcrypt.CompareHashAndPassword(storedPwd, pwd); err != nil {
		err = fmt.Errorf("were not in!: %v", err)
		return err
	}
	return nil

}

func createHash(password []byte) ([]byte, error) {
	var err error
	var createdHash []byte

	createdHash, err = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return createdHash, err
}
