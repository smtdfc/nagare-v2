package entrypoint

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/smtdfc/nagare/core/datastore/filesystem"
	"github.com/smtdfc/nagare/core/transport"
	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/uptrace/bun"
)

type Entrypoint struct {
	gateway            *transport.Gateway
	dataDirInitializer *filesystem.DataDirInitializer
	db                 *bun.DB
	Logger             *logger.Logger
}

func (e *Entrypoint) Start() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	e.Logger.Info("Initializing database")
	err := e.dataDirInitializer.Init()
	if err != nil {
		return err
	}

	e.Logger.Info("Starting gateway")
	e.gateway.Start()

	<-signalChan
	e.Logger.Info("Received termination signal, shutting down...")
	return e.Shutdown()
}

func (e *Entrypoint) Shutdown() error {
	e.Logger.Info("Disconnecting from the database")
	err := e.db.Close()
	if err != nil {
		return err
	}

	return nil
}

// @Injectable
// @Root
func Run(dataDirInitializer *filesystem.DataDirInitializer, db *bun.DB, logger *logger.Logger, gateway *transport.Gateway) *Entrypoint {
	return &Entrypoint{
		dataDirInitializer: dataDirInitializer,
		db:                 db,
		Logger:             logger,
		gateway:            gateway,
	}
}
