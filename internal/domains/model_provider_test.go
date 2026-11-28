package domains

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelProviderCompatible_OpenAI(t *testing.T) {
	assert.Equal(t, ModelProviderCompatible("openai"), OpenAI)
}

func TestModelProvider_StoresFields(t *testing.T) {
	provider := &ModelProvider{
		ID:         "provider-1",
		Name:       "OpenAI",
		BaseURL:    "https://api.openai.com/v1",
		ApiKey:     "secret",
		Compatible: OpenAI,
	}

	assert.Equal(t, "provider-1", provider.ID)
	assert.Equal(t, "OpenAI", provider.Name)
	assert.Equal(t, "https://api.openai.com/v1", provider.BaseURL)
	assert.Equal(t, "secret", provider.ApiKey)
	assert.Equal(t, OpenAI, provider.Compatible)
}
