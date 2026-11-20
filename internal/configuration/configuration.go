package configuration

import (
	"os"
	"path/filepath"
)

type Configuration struct {
	DataDir string
}

var DefaultConfiguration = Configuration{}

func init() {
	var userHomeDir, _ = os.UserHomeDir()
	DefaultConfiguration.DataDir = filepath.Join(userHomeDir, ".nagare-v2")
}
