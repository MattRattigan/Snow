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
	usr "os/user"
	"path/filepath"
	"runtime"
)

func main() {
	// Platform specific function calls main_windows or main_linux depending on environment
	go setupPlatformSpecific()

	ch := make(chan string)

	go func() {
		switch runtime.GOOS {
		case "windows":
			ch <- os.Getenv("LOCALAPPDATA")
		case "linux":
			userPath, err := usr.Current()
			if err != nil {
				log.Fatal(err)
			}
			ch <- filepath.Join(userPath.HomeDir, ".local", "bin")
		}
	}()

	path := func() (string, error) {
		defer close(ch)
		dbpath := filepath.Join(<-ch, "Snow")
		if _, err := os.Stat(dbpath); os.IsNotExist(err) {
			// MKdirAll makes the Snow directory
			err = os.MkdirAll(dbpath, 0700)
			if err != nil {
				return "", err
			}
		}
		return dbpath, nil
	}

	dbpath, err := path()
	if err != nil {
		log.Fatal(err)
	}

	// File path for snow.db is appended to the created or existing database
	dbpath = filepath.Join(dbpath, "Snow.db")
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
