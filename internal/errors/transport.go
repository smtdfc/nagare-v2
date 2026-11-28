package apperrors

var (
	ErrRequestBodyInvalid = NewError("Invalid request body", 400)
	ErrInternal           = NewError("Internal error", 500)
)
