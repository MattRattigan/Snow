package dbSnow

import (
	"Snow/snowUser"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
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
		err = db.createUser2(username, password)
		if err != nil {
			return snowUser.User{}, err
		}
	}

	id, err := db.getUserIDByUsername(username)

	if err != nil {
		return snowUser.User{}, nil
	}

	// Get User password from database
	storedPwd, err := db.getPasswordHashByID(id)

	if err != nil {
		return snowUser.User{}, err
	}

	err = isPassword(password, storedPwd)
	if err != nil {
		return snowUser.User{}, err
	}

	// Grab fileDir interface to build SnowUser struct
	fileDir, err := snowUser.BuildFileDir()
	if err != nil {
		return snowUser.User{}, err
	}

	fileDir = snowUser.AppendUUID(fileDir, *id)

	return snowUser.User{
		Username:   username,
		Passpharse: storedPwd,
		FileDir:    fileDir,
	}, nil

}

func (db *DbStore) createUser2(username string, password []byte) error {
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

func (db *DbStore) selectUsernameByUsername(username string) (bool, error) {
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT username FROM users WHERE username = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.QueryRow(username).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No result, but not necessarily an error
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *DbStore) getPasswordHashByID(id *snowUser.SnUUID) ([]byte, error) {
	var passwordHash []byte
	var byteArr []byte

	byteArr, err := id.ToBytesInplace(byteArr)
	if err != nil {
		return nil, err
	}

	query := `SELECT password_hash FROM users WHERE id = ?`
	err = db.QueryRow(query, byteArr).Scan(&passwordHash)
	if err != nil {
		return nil, err
	}
	return passwordHash, nil
}

func (db *DbStore) getUserIDByUsername(username string) (*snowUser.SnUUID, error) {
	var id []byte
	var snUUID snowUser.SnUUID

	query := `SELECT id FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&id)
	if err != nil {
		return &snowUser.SnUUID{}, err
	}

	snUUID, err = snowUser.SnUUID{}.FromBytes(id)
	if err != nil {
		return &snowUser.SnUUID{}, err
	}

	return &snUUID, nil
}

func isPassword(pwd, storedPwd []byte) error {
	if err := bcrypt.CompareHashAndPassword(storedPwd, pwd); err != nil {
		// TODO: nil pointer reference if incorrect password, see if can make a more detailed error
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
