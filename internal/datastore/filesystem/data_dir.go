package filesystem

import (
	"os"
	"path/filepath"

	"github.com/smtdfc/nagare/core/configuration"
)

type DataDirInitializer struct {
	isInit bool
	config *configuration.Configuration

	Paths *DataPaths
}

func (d *DataDirInitializer) Init() error {
	if d.isInit {
		return nil
	}

	paths := &DataPaths{
		DataDir: d.config.DataDir,
	}

	paths.DatabasesDir = filepath.Join(paths.DataDir, "databases")
	paths.MainDatabaseFile = filepath.Join(paths.DatabasesDir, "main.db")
	dirs := []string{
		paths.DataDir,
		paths.DatabasesDir,
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	d.Paths = paths
	d.isInit = true
	return nil
}

// @Injectable
func NewDataDirInitializer(config *configuration.Configuration) *DataDirInitializer {
	return &DataDirInitializer{
		config: config,
		isInit: false,
	}
}
