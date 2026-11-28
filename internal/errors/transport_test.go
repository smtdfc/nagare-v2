package apperrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransportErrors(t *testing.T) {
	assert.Equal(t, "Invalid request body", ErrRequestBodyInvalid.Message)
	assert.Equal(t, 400, ErrRequestBodyInvalid.Code)
	assert.Equal(t, "Internal error", ErrInternal.Message)
	assert.Equal(t, 500, ErrInternal.Code)
}
