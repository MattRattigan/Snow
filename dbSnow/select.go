package dbSnow

import (
	"Snow/snowUser"
	"database/sql"
	"errors"
)

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

func (db *DbStore) selectPasswordHashByID(id *snowUser.SnUUID) ([]byte, error) {
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

func (db *DbStore) selectUserIDByUsername(username string) (*snowUser.SnUUID, error) {
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
