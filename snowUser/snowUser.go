package snowUser

import (
	"Snow/snFlags"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/term"
	"log"
	"os"
)

type SnUUID uuid.UUID

type User struct {
	Username   string
	Passpharse []byte
	FileDir
}

type UserFile struct {
	UUID     SnUUID
	Data     []byte
	FileExt  string
	FilePath string
}

type UserDir struct {
	UUID    SnUUID
	Data    []byte
	FileExt string
	DirPath string
}

type FileDir interface {
	GetUUID() SnUUID
	GetPath() string
	GetData() []byte
}

func SetUsername(username string) string {
	if username == "" {
		fmt.Println("Username is required")
		os.Exit(1)
	}

	return username
}

func SetPassword() ([]byte, error) {
	fmt.Println("Enter Password: ")
	minPasswordLength := 8
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		fmt.Printf("error reading password: %v\n", err)
		os.Exit(1)
	}

	if len(passwordBytes) == 0 {
		fmt.Printf("password is required\n")
		os.Exit(1)
	}

	if len(passwordBytes) < minPasswordLength {
		fmt.Printf("password must be at least %d characters long\n", minPasswordLength)
		os.Exit(1)
	}

	return passwordBytes, nil
}

func (u UserFile) GetData() []byte {
	data, err := os.ReadFile(u.GetPath())
	if err != nil {
		log.Fatal(err)
	}
	u.Data = data

	return u.Data
}

func (u UserDir) GetData() []byte {
	return u.Data
}

func CreateFileUUID() SnUUID {
	return SnUUID{}.CreateUUID()
}

func BuildFileDir() (FileDir, error) {
	var flags = snFlags.CmdFlags

	if flags.DirPath == "" && flags.FilePath == "" {
		return nil, errors.New("file or directory was not provided")
	}

	if flags.FilePath != "" {
		return &UserFile{FilePath: flags.FilePath}, nil
	}

	if flags.DirPath != "" {
		return &UserDir{DirPath: flags.DirPath}, nil
	}

	return nil, errors.New("unexpected error in BuildFileDir")
}

func AppendUUID(f FileDir, uuid SnUUID) FileDir {
	path := f.GetPath()

	switch f.(type) {
	case *UserFile:
		return &UserFile{
			UUID:     uuid,
			FileExt:  "",
			FilePath: path,
		}
	case *UserDir:
		return &UserDir{
			UUID:    uuid,
			DirPath: path,
		}
	default:
		fmt.Println("Unknown type")
		return nil
	}
}

func (u UserFile) GetUUID() SnUUID {
	return u.UUID
}

func (u UserFile) GetPath() string {
	return u.FilePath
}

func (u UserDir) GetUUID() SnUUID {
	return u.UUID
}

func (u UserDir) GetPath() string {
	return u.DirPath
}

// FromBytes16 converts a [16]byte array to SnUUID
func (u SnUUID) FromBytes16(b [16]byte) (SnUUID, error) {
	id, err := uuid.FromBytes(b[:])
	if err != nil {
		return SnUUID{}, err
	}

	return SnUUID(id), nil
}

// FromBytes converts a byte slice to SnUUID
func (u SnUUID) FromBytes(b []byte) (SnUUID, error) {
	if len(b) != 16 {
		return SnUUID{}, errors.New("byte slice must be 16 bytes long")
	}

	id, err := uuid.FromBytes(b)
	if err != nil {
		return SnUUID{}, err
	}

	return SnUUID(id), nil
}

func (u SnUUID) CreateUUID() SnUUID {
	return SnUUID(uuid.New())
}

func (u SnUUID) String() string {
	return uuid.UUID(u).String()
}

func (u SnUUID) ToBytesInplace(b []byte) ([]byte, error) {
	newUuid := uuid.UUID(u)
	b, err := newUuid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return b, err
}
