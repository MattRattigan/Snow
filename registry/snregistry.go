package registry

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
)

type SN struct {
	fileExtension string
	iconPath      string
}

func Create() *SN {
	return &SN{
		fileExtension: ".sn",
		iconPath:      "./registry/sn.ico",
	}
}

func (s *SN) CreateRegistry() error {
	// Create a key in the Registry for the .registry file extension
	// located at HKEY_CLASSES_ROOT\.registry
	extKey, _, err := registry.CreateKey(registry.CLASSES_ROOT, s.fileExtension, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("error with file key extension creation, please run as Admin: %s ", err)
	}
	defer extKey.Close()

	// Sets the default value
	if err = extKey.SetStringValue("", "snfile"); err != nil {
		return fmt.Errorf("error with set snfile string value in registry: %s", err)
	}

	fileTypeKey, err := addFileType()
	if err != nil {
		return fmt.Errorf("error for file type: %s", err)
	}
	defer fileTypeKey.Close()

	defaultIconKey, err := createKeyForFile(fileTypeKey)
	if err != nil {
		return fmt.Errorf("error for DefaultIcon key: %s ", err)
	}
	defer defaultIconKey.Close()

	if err = addDefaultIcon(defaultIconKey, s); err != nil {
		return fmt.Errorf("error for setting Icon path in registry %s", err)
	}

	return nil
}

func addFileType() (*registry.Key, error) {
	fileTypeKey, _, err := registry.CreateKey(registry.CLASSES_ROOT, "snfile", registry.SET_VALUE)
	if err != nil {
		return nil, err
	}
	return &fileTypeKey, err
}

func createKeyForFile(fileTypeKey *registry.Key) (*registry.Key, error) {
	// Create a DefaultIcon subkey located at HKEY_CLASSES_ROOT\snfile\DefaultIcon
	defaultIconKey, _, err := registry.CreateKey(*fileTypeKey, "DefaultIcon", registry.SET_VALUE)
	if err != nil {
		return nil, err
	}
	return &defaultIconKey, err
}

func addDefaultIcon(defaultIconKey *registry.Key, s *SN) error {
	// Set the default icon path
	return defaultIconKey.SetStringValue("", s.iconPath)
}

func (s *SN) DoesFileExtensionExist() (bool, error) {
	var KeyNotFoundError = errors.New("key not found")

	key, err := registry.OpenKey(registry.CLASSES_ROOT, s.fileExtension, registry.QUERY_VALUE)
	defer key.Close()

	if err != nil {
		if errors.Is(err, KeyNotFoundError) {
			// The key does not exist
			return false, nil
		}
		// There was some other error opening the key
		return false, err
	}

	// The key exists
	return true, nil
}
