package domains

type ModelProviderCompatible string

const (
	OpenAI ModelProviderCompatible = "openai"
)

type ModelProvider struct {
	ID         string
	Name       string
	BaseURL    string
	ApiKey     string
	Compatible ModelProviderCompatible
}
