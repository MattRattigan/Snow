package dbSnow

import (
	"Snow/snowUser"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
)

type DbStore struct {
	*sql.DB
}

func InitDB(filepath string) *DbStore {
	dbSnow, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "username" TEXT NOT NULL UNIQUE,
        "password_hash" BLOB NOT NULL,
        "UUID" TEXT 
    );`

	_, err = dbSnow.Exec(createTableSQL)

	if err != nil {
		log.Fatal(err)
	}

	return &DbStore{dbSnow}
}

func (db *DbStore) CreateUser(username string, passwordHash []byte) error {
	idUsername, exists, err := db.checkUsername(username)

	if err != nil && exists == false {
		log.Fatal("An error has occurred checking through the database: ", err)
	}

	if err == nil && exists == false {
		log.Fatal("Username was not found")
	}

	isTrue, err := db.checkPasswordHash(idUsername, passwordHash)

	if isTrue == false && err == nil {
		log.Fatal("User ID not found")
	}

	if isTrue == false && err != nil {
		log.Fatal("An error has occurred checking through the database: ")
	}

	if isTrue == false {
		log.Fatal("Incorrect Password!")
	}

	insertUserSQL := `INSERT INTO users (username, password_hash) VALUES (?, ?)`

	statement, err := db.Prepare(insertUserSQL)
	if err != nil {
		return fmt.Errorf("at insertion and err has occured: %s", err) // TODO: Check for error
	}
	defer statement.Close()

	_, err = statement.Exec(username, passwordHash) // TODO: Create function for fileIndicator
	return err
}

func (db *DbStore) GetUser(username string, passwordHash []byte) (*snowUser.User, error) {
	queryUserSQL := `SELECT username, password_hash FROM users WHERE username = ?`

	// Need to add check for passpharse before populating snowUser struct
	row := db.QueryRow(queryUserSQL, username)

	// passing pointers, updates username, password variables
	err := row.Scan(&username, &passwordHash)
	if err != nil {
		return &snowUser.User{}, err // TODO: Check for error when username or password does not match
	}

	// build SnowUser struct
	return &snowUser.User{
		Username:   username,
		Passpharse: passwordHash, // TODO: FileIndicator or something needs to of os.File
	}, nil
}

func (db *DbStore) checkUsername(username string) (int, bool, error) {
	var id int

	if username == "" {
		fmt.Println("Username is required")
		os.Exit(1)
	}

	// Query for the username
	row := db.QueryRow("SELECT id FROM users WHERE username = ?", username)

	// Scan the result into the id variable
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Username not found
			return 0, false, nil
		}
		// Some other error occurred
		return 0, false, err
	}

	// Username found
	return id, true, nil
}

func (db *DbStore) checkPasswordHash(userID int, providedHash []byte) (bool, error) {
	var storedHash []byte

	// Query for the password hash
	err := db.QueryRow("SELECT password_hash FROM users WHERE id = ?", userID).Scan(&storedHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User ID not found
			return false, nil
		}
		// Some other error occurred
		return false, err
	}

	// Compare the provided hash with the stored hash
	return bytes.Equal(storedHash, providedHash), nil
}
