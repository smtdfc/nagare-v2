package entrypoint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/configuration"
	"github.com/smtdfc/nagare/core/datastore/database"
	"github.com/smtdfc/nagare/core/datastore/filesystem"
	"github.com/smtdfc/nagare/core/transport"
	httptransport "github.com/smtdfc/nagare/core/transport/http"
	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestEntrypoint(t *testing.T) *Entrypoint {
	t.Helper()

	dir := t.TempDir()
	initializer := filesystem.NewDataDirInitializer(&configuration.Configuration{
		DataDir: filepath.Join(dir, "data"),
	})
	require.NoError(t, initializer.Init())

	db, err := database.InitDatabase(initializer.Paths)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	log := logger.NewLogger()
	server := httptransport.NewServer(
		httptransport.ProvideFiberApp(),
		&configuration.Configuration{Host: "127.0.0.1", Port: "0"},
		log,
	)
	gateway := transport.NewGateway(server)
	entry := Run(initializer, db, log, gateway)
	require.NotNil(t, entry)
	return entry
}

func TestRun(t *testing.T) {
	entry := newTestEntrypoint(t)

	assert.NotNil(t, entry.dataDirInitializer)
	assert.NotNil(t, entry.db)
	assert.NotNil(t, entry.Logger)
	assert.NotNil(t, entry.gateway)
}

func TestShutdown_ClosesDatabase(t *testing.T) {
	entry := newTestEntrypoint(t)

	require.NoError(t, entry.Shutdown())
	assert.Error(t, entry.db.Ping(), "database should be closed after shutdown")
}

func TestStart_ReturnsErrorWhenDataDirInitFails(t *testing.T) {
	entry := newTestEntrypoint(t)

	base := t.TempDir()
	blocker := filepath.Join(base, "blocker")
	require.NoError(t, os.WriteFile(blocker, []byte("x"), 0o644))
	entry.dataDirInitializer = filesystem.NewDataDirInitializer(&configuration.Configuration{
		DataDir: filepath.Join(blocker, "nagare"),
	})

	// Start must return before waiting for a signal when initialization fails.
	assert.Error(t, entry.Start())
}
