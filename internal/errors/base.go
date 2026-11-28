package apperrors

import (
	"errors"
)

type BaseError struct {
	error
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewError(message string, code int) *BaseError {
	return &BaseError{
		error:   errors.New(message),
		Message: message,
		Code:    code,
	}
}
