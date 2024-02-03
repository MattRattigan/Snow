//go:build windows

package main

import (
	"Snow/registry"
	"fmt"
	"log"
)

func setupPlatformSpecific() error {
	reg, err := registry.Create()
	if err != nil {
		return err
	}
	if ok, err := reg.DoesFileExtensionExist(); !ok {
		err = reg.CreateRegistry()
		fmt.Println("Created .sn extension")
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
