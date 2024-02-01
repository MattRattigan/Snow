package main

import (
	db "Snow/dbSnow"
	"Snow/registry"
	"Snow/snFlags"
	"Snow/sncrypt"
	"Snow/snowUser"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	var dbstore = db.InitDB("data/snow_database.sqlite")
	reg := registry.Create()

	if ok, err := reg.DoesFileExtensionExist(); !ok {
		err = reg.CreateRegistry()
		fmt.Println("Created .sn extension")
		if err != nil {
			log.Fatal(err)
		}
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
		fmt.Println("no encryption option was chosen")
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
