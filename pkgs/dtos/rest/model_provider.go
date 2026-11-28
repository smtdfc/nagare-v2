package rest

type ModelProviderCompatible string

const (
	OpenAI ModelProviderCompatible = "openai"
)

type ModelProvider struct {
	ID         string                  `json:"id"`
	Name       string                  `json:"name"`
	BaseURL    string                  `json:"base_url"`
	Compatible ModelProviderCompatible `json:"compatible"`
}

// Add model provider
var AddModelProviderEndpoint = "/api/v1/model-provider"
var AddModelProviderMethod = POST

type AddModelProviderRequest struct {
	ID         string                  `json:"id"`
	Name       string                  `json:"name"`
	BaseURL    string                  `json:"base_url"`
	Compatible ModelProviderCompatible `json:"compatible"`
	ApiKey     string                  `json:"api_key"`
}

type AddModelProviderResponse struct {
	ModelProvider *ModelProvider `json:"model_provider"`
}
