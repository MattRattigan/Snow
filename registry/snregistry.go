package registry

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

type SN struct {
	fileExtension string
	iconPath      string
}

func Create() SN {
	// Set the file extension and the associated icon path
	return SN{
		fileExtension: ".sn",
		iconPath:      "./registry/sn.ico",
	}
}

func (s *SN) AddToRegistry() error {
	var extKey registry.Key
	var regError error

	// Create a key in the Registry for the .registry file extension
	// located at HKEY_CLASSES_ROOT\.registry
	extKey, _, regError = registry.CreateKey(registry.CLASSES_ROOT, s.fileExtension, registry.SET_VALUE)
	if regError != nil {
		regError = fmt.Errorf("Error with file key extension creation: %s ", regError)
	}

	// Sets the default value
	regError = extKey.SetStringValue("", "snfile")
	if regError != nil {
		regError = fmt.Errorf("error with set snfile string value in registry: %s", regError)
	}
	extKey.Close() // close extKey

	key, err := s.addFileType()
	if err != nil {
		regError = fmt.Errorf("error for file type: %s", err)
	}

	key, err = s.createKeyForFile(key)
	if err != nil {
		regError = fmt.Errorf("error for DefaultIcon key: %s ", err)
	}
	err = s.addDefaultIcon(key)
	if err != nil {
		regError = fmt.Errorf("error for setting Icon path in registry %s", err)
	}

	return regError

}

func (s *SN) addFileType() (*registry.Key, error) {
	fileTypeKey, _, err := registry.CreateKey(registry.CLASSES_ROOT, "snfile", registry.SET_VALUE)
	if err != nil {
		return nil, err
	}
	return &fileTypeKey, err
}

func (s *SN) createKeyForFile(fileTypeKey *registry.Key) (*registry.Key, error) {
	// Create a DefaultIcon subkey located at HKEY_CLASSES_ROOT\snfile\DefaultIcon
	defaultIconKey, _, err := registry.CreateKey(*fileTypeKey, "DefaultIcon", registry.SET_VALUE)
	if err != nil {
		return nil, err
	}
	defer fileTypeKey.Close()
	return &defaultIconKey, err
}

func (s *SN) addDefaultIcon(defaultIconKey *registry.Key) (err error) {
	// Set the default icon path
	err = defaultIconKey.SetStringValue("", s.iconPath)
	if err != nil {
		return err
	}

	defer defaultIconKey.Close()
	return err
}
