package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModelProviderService(t *testing.T) {
	service := NewModelProviderService()

	require.NotNil(t, service)
	assert.NotNil(t, service)
}
