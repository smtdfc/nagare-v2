package error_utils

import (
	"errors"
	"strconv"
)

type BaseError struct {
	error
	Message string `json:"message"`
	Code    string `json:"code"`
}

func NewError(message string, code int) *BaseError {
	return &BaseError{
		error:   errors.New(message),
		Message: message,
		Code:    strconv.Itoa(code),
	}
}
