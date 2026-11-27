package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataDirInitializer(t *testing.T) {
	config := &configuration.Configuration{DataDir: "/tmp/nagare"}
	initializer := NewDataDirInitializer(config)

	require.NotNil(t, initializer)
	assert.Equal(t, config, initializer.config)
	assert.False(t, initializer.isInit)
	assert.Nil(t, initializer.Paths)
}

func TestDataDirInitializer_Init_CreatesDirsAndPaths(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "nagare")
	initializer := NewDataDirInitializer(&configuration.Configuration{DataDir: dir})

	require.NoError(t, initializer.Init())

	for _, want := range []string{
		dir,
		filepath.Join(dir, "databases"),
	} {
		info, err := os.Stat(want)
		require.NoError(t, err, want)
		assert.True(t, info.IsDir(), want)
	}

	require.NotNil(t, initializer.Paths)
	assert.Equal(t, dir, initializer.Paths.DataDir)
	assert.Equal(t, filepath.Join(dir, "databases"), initializer.Paths.DatabasesDir)
	assert.Equal(t, filepath.Join(dir, "databases", "main.db"), initializer.Paths.MainDatabaseFile)
	assert.True(t, initializer.isInit)
}

func TestDataDirInitializer_Init_IsIdempotent(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nagare")
	initializer := NewDataDirInitializer(&configuration.Configuration{DataDir: dir})

	require.NoError(t, initializer.Init())
	first := initializer.Paths

	require.NoError(t, initializer.Init())

	assert.Same(t, first, initializer.Paths)
}

func TestDataDirInitializer_Init_ReturnsErrorWhenParentIsFile(t *testing.T) {
	base := t.TempDir()
	blocker := filepath.Join(base, "blocker")
	require.NoError(t, os.WriteFile(blocker, []byte("x"), 0o644))

	initializer := NewDataDirInitializer(&configuration.Configuration{
		DataDir: filepath.Join(blocker, "nagare"),
	})

	assert.Error(t, initializer.Init())
	assert.Nil(t, initializer.Paths)
	assert.False(t, initializer.isInit)
}

func TestDataDirInitializer_Init_ReturnsErrorWhenDataDirEmpty(t *testing.T) {
	initializer := NewDataDirInitializer(&configuration.Configuration{})

	assert.Error(t, initializer.Init())
	assert.Nil(t, initializer.Paths)
}

func TestGetDataPaths(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nagare")
	initializer := NewDataDirInitializer(&configuration.Configuration{DataDir: dir})

	paths, err := GetDataPaths(initializer)
	require.NoError(t, err)
	require.NotNil(t, paths)
	assert.Equal(t, filepath.Join(dir, "databases", "main.db"), paths.MainDatabaseFile)
}

func TestGetDataPaths_ReturnsErrorOnInitFailure(t *testing.T) {
	initializer := NewDataDirInitializer(&configuration.Configuration{})

	paths, err := GetDataPaths(initializer)

	assert.Error(t, err)
	assert.Nil(t, paths)
}
