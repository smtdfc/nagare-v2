package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ModelProviderCompatible int

const (
	OpenAI ModelProviderCompatible = 1
)

type ModelProvider struct {
	bun.BaseModel `bun:"table:model_providers,alias:mp"`

	ID         string                  `bun:",pk,type:uuid"`
	Name       string                  `bun:",notnull"`
	BaseURL    string                  `bun:",notnull"`
	ApiKey     string                  `bun:",default:''"`
	Compatible ModelProviderCompatible `bun:",default:1"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
