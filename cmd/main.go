package main

import (
	db "Snow/dbSnow"
	"Snow/snFlags"
	"Snow/sncrypt"
	"Snow/snowUser"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

func main() {
	path := func() (string, error) {
		execPath, err := os.Executable()
		if err != nil {
			return "", err
		}
		execDir := filepath.Dir(execPath)
		dbPath := filepath.Join(execDir, "dbSnow/data/snow_database.db")
		return dbPath, nil
	}

	dbpath, err := path()
	if err != nil {
		log.Fatal(err)
	}

	// Platform specific function calls main_windows or main_linux depending on environment
	if err = setupPlatformSpecific(); err != nil {
		log.Fatal(err)
	}

	if err = db.ExtractDatabase(dbpath); err != nil {
		log.Fatalf("Failed to extract database: %v", err)
	}

	dbstore, err := db.InitDB(dbpath)

	if err != nil {
		log.Fatalf("error in database init: %v\n", err)
	}

	cmdFlags := snFlags.CmdFlags
	username := snowUser.SetUsername(cmdFlags.Username)
	password, err := snowUser.SetPassword()

	if err != nil {
		os.Exit(1)
	}

	snUser, err := dbstore.GetUser(username, password)
	if err != nil {
		log.Fatal(err)
	}

	if snFlags.CmdFlags.Encrypt == false && snFlags.CmdFlags.Decrypt == false {
		log.Fatal("no encryption option was chosen")
	} else if snFlags.CmdFlags.Ext == "" && snFlags.CmdFlags.Decrypt == true {
		log.Fatal("no file extension was given with -d flag")
	} else if snFlags.CmdFlags.Encrypt {
		err = sncrypt.WriteEncryption(&snUser)
		if err != nil {
			log.Fatal(err)
		}
	} else if snFlags.CmdFlags.Decrypt {
		err = sncrypt.WriteDecryption(&snUser)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("\nFinished!!!")

}
