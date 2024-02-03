package dbSnow

import (
	"Snow/snowUser"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

//go:embed data/snow_database.db
var embeddedDB embed.FS

var (
	once            sync.Once
	dbStoreInstance *DbStore
)

type DbStore struct {
	*sql.DB
}

func ExtractDatabase(targetPath string) error {
	// Check if the database file already exists
	if _, err := os.Stat(targetPath); err == nil {
		// File already exists, no need to extract
		return nil
	}

	// File doesn't exist, extract it from embedded resources
	data, readErr := fs.ReadFile(embeddedDB, "data/snow_database.db")
	if readErr != nil {
		return readErr
	}

	// Ensure the directory exists
	dir := filepath.Dir(targetPath)
	if mkErr := os.MkdirAll(dir, os.ModePerm); mkErr != nil {
		return mkErr
	}

	// Write the data to the target path with appropriate file permissions
	return os.WriteFile(targetPath, data, 0600)
}

func InitDB(filePath string) (*DbStore, error) {
	var err error
	once.Do(func() {
		// Open or create the database
		dbSnow, dbErr := sql.Open("sqlite3", filePath)
		if dbErr != nil {
			err = dbErr // Capture the error
			return      // Exit the once.Do block
		}

		createTableSQL := `
			CREATE TABLE IF NOT EXISTS users (
				"id" BLOB PRIMARY KEY,
				"username" TEXT NOT NULL UNIQUE,
				"password_hash" BLOB NOT NULL
			);`

		if _, dbErr = dbSnow.Exec(createTableSQL); dbErr != nil {
			err = dbErr
			return
		}

		dbStoreInstance = &DbStore{dbSnow}
	})

	if err != nil {
		return nil, err
	}
	return dbStoreInstance, nil
}

func (db *DbStore) GetUser(username string, password []byte) (snowUser.User, error) {
	// Get ID and password from database, check password and populate the struct

	// If user is in database proceed to getUser if not create user
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
