//go:build linux

package main

func setupPlatformSpecific() <-chan error {
	ch := make(chan error)
	defer close(ch)
	return ch
}
