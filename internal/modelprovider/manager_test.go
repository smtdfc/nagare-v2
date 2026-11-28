package modelprovider

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModelProviderManager(t *testing.T) {
	log := logger.NewLogger()

	manager := NewModelProviderManager(log)

	require.NotNil(t, manager)
	assert.Same(t, log, manager.logger)
}

func TestNewModelProviderManager_NilLogger(t *testing.T) {
	manager := NewModelProviderManager(nil)

	require.NotNil(t, manager)
	assert.Nil(t, manager.logger)
}
