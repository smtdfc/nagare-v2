package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/datastore/filesystem"
	error_utils "github.com/smtdfc/nagare/core/error-utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitDatabase(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "databases", "main.db")
	require.NoError(t, os.MkdirAll(filepath.Dir(dbPath), 0o755))

	db, err := InitDatabase(&filesystem.DataPaths{MainDatabaseFile: dbPath})
	require.NoError(t, err)
	require.NotNil(t, db)
	t.Cleanup(func() { _ = db.Close() })

	require.NoError(t, db.Ping())
	assert.NoError(t, db.QueryRow("SELECT 1").Scan(new(int)))

	_, statErr := os.Stat(dbPath)
	assert.NoError(t, statErr, "sqlite file should be created at the configured path")
}

func TestRunMigration_CreatesModelProviderTable(t *testing.T) {
	db, err := InitDatabase(&filesystem.DataPaths{
		MainDatabaseFile: filepath.Join(t.TempDir(), "main.db"),
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	var count int
	err = db.NewRaw(
		"SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?",
		"model_providers",
	).Scan(t.Context(), &count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestInitDatabase_PathWithSpaces(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "my data", "main.db")
	require.NoError(t, os.MkdirAll(filepath.Dir(dbPath), 0o755))

	db, err := InitDatabase(&filesystem.DataPaths{MainDatabaseFile: dbPath})
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	assert.NoError(t, db.Ping())
	_, statErr := os.Stat(dbPath)
	assert.NoError(t, statErr)
}

func TestInitDatabase_ReopenSameFile(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "main.db")
	paths := &filesystem.DataPaths{MainDatabaseFile: dbPath}

	first, err := InitDatabase(paths)
	require.NoError(t, err)
	require.NoError(t, first.Close())

	second, err := InitDatabase(paths)
	require.NoError(t, err)
	t.Cleanup(func() { _ = second.Close() })

	assert.NoError(t, second.Ping())
}

func TestInitDatabase_RejectsEmptyPath(t *testing.T) {
	tests := []struct {
		name  string
		paths *filesystem.DataPaths
	}{
		{name: "nil paths", paths: nil},
		{name: "empty struct", paths: &filesystem.DataPaths{}},
		{name: "empty file path", paths: &filesystem.DataPaths{DataDir: "/tmp", DatabasesDir: "/tmp/databases"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := InitDatabase(tt.paths)

			assert.ErrorIs(t, err, error_utils.ErrDatabasePathEmpty)
			assert.Nil(t, db)
		})
	}
}
