package rest

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelProviderJSON(t *testing.T) {
	provider := &ModelProvider{
		ID:         "provider-1",
		Name:       "OpenAI",
		BaseURL:    "https://api.openai.com/v1",
		Compatible: OpenAI,
	}

	payload, err := json.Marshal(provider)
	require.NoError(t, err)
	assert.JSONEq(t, `{"id":"provider-1","name":"OpenAI","base_url":"https://api.openai.com/v1","compatible":"openai"}`, string(payload))
}

func TestAddModelProviderRequestJSON(t *testing.T) {
	request := &AddModelProviderRequest{
		ID:         "provider-1",
		Name:       "OpenAI",
		BaseURL:    "https://api.openai.com/v1",
		Compatible: OpenAI,
		ApiKey:     "secret",
	}

	payload, err := json.Marshal(request)
	require.NoError(t, err)
	assert.JSONEq(t, `{"id":"provider-1","name":"OpenAI","base_url":"https://api.openai.com/v1","compatible":"openai","api_key":"secret"}`, string(payload))
}

func TestApiResponseErrorJSON(t *testing.T) {
	response := &ApiResponse[any]{
		IsSuccess: false,
		Data:      nil,
		Error:     ApiResponseError{Message: "bad request"},
	}

	payload, err := json.Marshal(response)
	require.NoError(t, err)
	assert.JSONEq(t, `{"is_success":false,"data":null,"error":{"message":"bad request"}}`, string(payload))
}
