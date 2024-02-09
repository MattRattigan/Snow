//go:build windows

package registry

import (
	"errors"
	"golang.org/x/sys/windows/registry"
)

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
		return false, errors.New("error opening the key")
	}

	// The key exists
	return true, nil
}
