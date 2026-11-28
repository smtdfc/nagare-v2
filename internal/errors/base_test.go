package apperrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	err := NewError("Error", 500)
	assert.NotNil(t, err)
	assert.NotNil(t, err.error)
	assert.Equal(t, "Error", err.Error())
	assert.Equal(t, 500, err.Code)
}
