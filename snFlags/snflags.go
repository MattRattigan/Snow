package snFlags

import "flag"

var CmdFlags = Flags()

type Config struct {
	Username string
	FilePath string
	DirPath  string
	Encrypt  bool
	Decrypt  bool
	Ext      string
}

func Flags() *Config {
	var cfg Config
	flag.StringVar(&cfg.Username, "username", "", "Username for login")
	flag.StringVar(&cfg.FilePath, "filepath", "", "Path to the file")
	flag.StringVar(&cfg.DirPath, "dirpath", "", "Path to the directory")
	flag.BoolVar(&cfg.Encrypt, "e", false, "Encrypt file") // Changed to BoolVar
	flag.BoolVar(&cfg.Decrypt, "d", false, "Decrypt file")
	flag.StringVar(&cfg.Ext, "ext", "", "name of file extension")

	flag.Parse()

	return &cfg
}
