//go:build windows

package main

import (
	"Snow/registry"
	"fmt"
)

func setupPlatformSpecific() <-chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		reg, err := registry.Create()
		if err != nil {
			ch <- fmt.Errorf("error with registry creation: %s", err)
			return
		}
		if ok, feErr := reg.DoesFileExtensionExist(); !ok {
			if feErr != nil {
				err = reg.CreateRegistry()
				fmt.Println("Created .sn extension")
				if err != nil {
					ch <- fmt.Errorf("%s: %s", feErr, err)
				}
			}
		}
	}()
	return ch
}
