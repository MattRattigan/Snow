package snFlags

import "flag"

func Flags() map[string]string {
	cmdFlags := make(map[string]string)

	//username flag
	username := *flag.String("username", "", "Username for login")

	// file flag
	filePath := *flag.String("filepath", "", "Path to the file")

	// directory flag
	dirPath := *flag.String("dirpath", "", "Path to the directory")

	// parse snFlags
	flag.Parse()

	cmdFlags["username"] = username
	cmdFlags["filePath"] = filePath
	cmdFlags["dirPath"] = dirPath

	return cmdFlags
}
