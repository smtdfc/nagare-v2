package entrypoint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/configuration"
	"github.com/smtdfc/nagare/core/datastore/database"
	"github.com/smtdfc/nagare/core/datastore/filesystem"
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

	entry := Run(initializer, db)
	require.NotNil(t, entry)
	return entry
}

func TestRun(t *testing.T) {
	entry := newTestEntrypoint(t)

	assert.NotNil(t, entry.dataDirInitializer)
	assert.NotNil(t, entry.db)
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
