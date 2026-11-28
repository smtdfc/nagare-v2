package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModelProviderController(t *testing.T) {
	controller := NewModelProviderController()

	require.NotNil(t, controller)
	assert.NotNil(t, controller)
}
