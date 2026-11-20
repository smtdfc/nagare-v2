package entrypoint

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/smtdfc/nagare/core/datastore/filesystem"
	"github.com/uptrace/bun"
)

type Entrypoint struct {
	dataDirInitializer *filesystem.DataDirInitializer
	db                 *bun.DB
}

func (e *Entrypoint) Start() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	err := e.dataDirInitializer.Init()
	if err != nil {
		return err
	}

	<-signalChan
	fmt.Println("Received termination signal, shutting down...")
	return e.Shutdown()
}

func (e *Entrypoint) Shutdown() error {
	err := e.db.Close()
	if err != nil {
		return err
	}

	return nil
}

// @Injectable
// @Root
func Run(dataDirInitializer *filesystem.DataDirInitializer, db *bun.DB) *Entrypoint {
	return &Entrypoint{
		dataDirInitializer: dataDirInitializer,
		db:                 db,
	}
}
