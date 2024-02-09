//go:build windows

package registry

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
)

type SN struct {
	fileExtension string
	iconPath      string
}

func Create() (*SN, error) {
	iconPath, err := getIconPath()
	if err != nil {
		return nil, err
	}
	return &SN{
		fileExtension: ".sn",
		iconPath:      iconPath,
	}, nil
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

func getIconPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("with executable path: %s", err)
	}
	execDir := filepath.Dir(execPath)

	// sn.ico is located in the same directory as the executable
	iconPath := filepath.Join(execDir, "sn.ico")
	return iconPath, nil
}
