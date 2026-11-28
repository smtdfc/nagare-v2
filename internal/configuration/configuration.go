package configuration

import (
	"os"
	"path/filepath"
)

type Configuration struct {
	DataDir string
	Port    string
	Host    string
}

var DefaultConfiguration = Configuration{}

func init() {
	var userHomeDir, _ = os.UserHomeDir()
	DefaultConfiguration.DataDir = filepath.Join(userHomeDir, ".nagare")
	DefaultConfiguration.Port = "9832"
}
