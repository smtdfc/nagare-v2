package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/datastore/database/models"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/uptrace/bun"
)

type ModelProviderRepository struct {
	db *bun.DB
}

func (m *ModelProviderRepository) mapToModel(entity *domains.ModelProvider) *models.ModelProvider {
	var comp models.ModelProviderCompatible
	switch entity.Compatible {
	case domains.OpenAI:
		comp = models.OpenAI
	}

	return &models.ModelProvider{
		ID:         entity.ID,
		Name:       entity.Name,
		BaseURL:    entity.BaseURL,
		ApiKey:     entity.ApiKey,
		Compatible: comp,
	}
}

func (m *ModelProviderRepository) Create(ctx context.Context, entity *domains.ModelProvider) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	model := m.mapToModel(entity)
	model.ID = id.String()
	if _, err := m.db.NewInsert().Model(model).Exec(ctx); err != nil {
		return err
	}

	entity.ID = model.ID
	return nil
}

// @Injectable
func NewModelProviderRepository(db *bun.DB) *ModelProviderRepository {
	return &ModelProviderRepository{db: db}
}
