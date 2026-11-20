package filesystem

import (
	"os"

	"github.com/smtdfc/nagare/core/configuration"
)

type DataDirInitializer struct {
	config *configuration.Configuration
}

func (d *DataDirInitializer) Init() error {
	dirs := []string{d.config.DataDir}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// @Injectable
func NewDataDirInitializer(config *configuration.Configuration) *DataDirInitializer {
	return &DataDirInitializer{
		config: config,
	}
}
