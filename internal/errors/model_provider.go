package apperrors

var (
	ErrAdapterNotFound       = NewError("Agent adapter not found", 400)
	ErrModelNotFound         = NewError("Model not found", 404)
	ErrModelQuotaExceed      = NewError("You have reached the usage limit for this model. Please check with your model provider or contact them.", 429)
	ErrModelProviderResponse = NewError("The model provider returned an error", 400)
)
