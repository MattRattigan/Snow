//go:build windows

package main

import (
	"Snow/registry"
	"fmt"
	"log"
)

func setupPlatformSpecific() {
	reg, err := registry.Create()
	if err != nil {
		log.Fatalf("error with registry creation: %s", err)
		return
	}
	if ok, feErr := reg.DoesFileExtensionExist(); !ok {
		if feErr != nil {
			err = reg.CreateRegistry()
			fmt.Println("Created .sn extension")
			if err != nil {
				log.Fatalf("%s: %s", feErr, err)
			}
		}
	}

}
