package main

import (
	"Snow/dbSnow"
	"Snow/registry"
	"Snow/snowUser"
	"fmt"
	"os"
)

func main() {
	db := dbSnow.DbStore{}

	//snUser := &snowUser.User{}
	//var filePath string
	//var err error

	//cmdFlags := snFlags()
	//username := cmdFlags["flags1"]
	//passpharse, err := snUser.SetUsername(*username).SetPasspharse()

	//if err != nil {
	//	_, _ = fmt.Fprintln(os.Stderr, err)
	//}
	//
	//db := dbSnow.InitDB("users.dbSnow")
	//defer db.Close()
	//
	//// TODO: Example of create (Delete or change later)
	//err = db.CreateUser(snUser) // TODO: Create file path method
	//
	//if err != nil {
	//	log.Fatal("Failed to create user:", err)
	//}
	//
	//// TODO: Example of get (Delete of change later)
	//user, err := db.GetUser("alice")
	//
	//if err != nil {
	//	log.Fatal("Failed to retrieve user: ", err)
	//}
	//
	//_ = filePath
	//fmt.Println(user)
	//
	///// demonstration
	//potato := sncrypt.Encrypt()

	su, err := snowUser.CreateUser(&db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "We have an error!")
	}

	_ = su

}

func registryCode() {
	sn := registry.Create()

	err := sn.AddToRegistry()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Additional steps may be required to refresh the icon cache
}
