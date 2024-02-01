package dbSnow

import (
	"Snow/snowUser"
	"database/sql"
	"fmt"
	"log"
	"sync"
)

var (
	once            sync.Once
	dbStoreInstance *DbStore
)

type DbStore struct {
	*sql.DB
}

func InitDB(filepath string) *DbStore {
	once.Do(func() {
		dbSnow, err := sql.Open("sqlite3", filepath)
		if err != nil {
			log.Fatal(err)
		}

		createTableSQL := `
			CREATE TABLE IF NOT EXISTS users (
				"id" BLOB PRIMARY KEY,
				"username" TEXT NOT NULL UNIQUE,
				"password_hash" BLOB NOT NULL
			);`

		_, err = dbSnow.Exec(createTableSQL)
		if err != nil {
			log.Fatal(err)
		}

		dbStoreInstance = &DbStore{dbSnow}
	})
	return dbStoreInstance
}

func (db *DbStore) GetUser(username string, password []byte) (snowUser.User, error) {
	// Get ID and password from database, check password and populate the struct

	// If user is in database procede to getUser if not create user
	if ok, err := db.selectUsernameByUsername(username); err != nil {
		return snowUser.User{}, err
	} else if !ok {
		err = db.createUser(username, password)
		if err != nil {
			return snowUser.User{}, err
		}
	}

	id, err := db.selectUserIDByUsername(username)

	if err != nil {
		return snowUser.User{}, nil
	}

	// Get User password from database
	storedPwd, err := db.selectPasswordHashByID(id)

	if err != nil {
		return snowUser.User{}, err
	}

	// validate users password
	err = isPassword(password, storedPwd)
	if err != nil {
		return snowUser.User{}, err
	}

	// Grab fileDir interface to build SnowUser struct
	fileDir, err := snowUser.BuildFileDir()
	if err != nil {
		return snowUser.User{}, err
	}

	// appending uuid to fileDir type
	fileDir = snowUser.AppendUUID(fileDir, *id)

	return snowUser.User{
		Username:   username,
		Passpharse: storedPwd,
		FileDir:    fileDir,
	}, nil

}

func (db *DbStore) createUser(username string, password []byte) error {
	var err error
	var uuid, createdHash []byte

	// Create UUID for the User
	uuid, err = snowUser.CreateFileUUID().ToBytesInplace(uuid)

	if err != nil {
		return err
	}

	// Create the Hash
	createdHash, err = createHash(password)

	if err != nil {
		return err
	}

	// Insert into the database the username, password "hashed" and UUID as id
	insertUserSQL := `INSERT INTO users (username, password_hash, id) VALUES (?, ?, ?);`

	statement, err := db.Prepare(insertUserSQL)
	if err != nil {
		return fmt.Errorf("at insertion and err has occured: %s", err) // TODO: Check for error
	}

	defer statement.Close()

	_, err = statement.Exec(username, createdHash, uuid)

	return err
}
