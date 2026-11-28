package repositories

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/datastore/database"
	"github.com/smtdfc/nagare/core/datastore/database/models"
	"github.com/smtdfc/nagare/core/datastore/filesystem"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModelProviderRepository(t *testing.T) {
	db, err := database.InitDatabase(&filesystem.DataPaths{
		MainDatabaseFile: filepath.Join(t.TempDir(), "main.db"),
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	repository := NewModelProviderRepository(db)

	require.NotNil(t, repository)
	assert.Same(t, db, repository.db)
}

func TestModelProviderRepository_Create(t *testing.T) {
	db, err := database.InitDatabase(&filesystem.DataPaths{
		MainDatabaseFile: filepath.Join(t.TempDir(), "main.db"),
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	repository := NewModelProviderRepository(db)
	entity := &domains.ModelProvider{
		Name:       "OpenAI",
		BaseURL:    "https://api.openai.com/v1",
		ApiKey:     "secret",
		Compatible: domains.OpenAI,
	}

	require.NoError(t, repository.Create(context.Background(), entity))
	_, err = uuid.Parse(entity.ID)
	assert.NoError(t, err)

	var stored models.ModelProvider
	require.NoError(t, db.NewSelect().Model(&stored).Where("id = ?", entity.ID).Scan(context.Background()))
	assert.Equal(t, entity.ID, stored.ID)
	assert.Equal(t, entity.Name, stored.Name)
	assert.Equal(t, entity.BaseURL, stored.BaseURL)
	assert.Equal(t, entity.ApiKey, stored.ApiKey)
	assert.Equal(t, models.OpenAI, stored.Compatible)
	assert.NotZero(t, stored.CreatedAt)
	assert.NotZero(t, stored.UpdatedAt)
}

func TestModelProviderRepository_MapToModel(t *testing.T) {
	repository := &ModelProviderRepository{}
	entity := &domains.ModelProvider{
		ID:         "provider-1",
		Name:       "OpenAI",
		BaseURL:    "https://api.openai.com/v1",
		ApiKey:     "secret",
		Compatible: domains.OpenAI,
	}

	mapped := repository.mapToModel(entity)

	assert.Equal(t, entity.ID, mapped.ID)
	assert.Equal(t, entity.Name, mapped.Name)
	assert.Equal(t, entity.BaseURL, mapped.BaseURL)
	assert.Equal(t, entity.ApiKey, mapped.ApiKey)
	assert.Equal(t, models.OpenAI, mapped.Compatible)
}
