package http

import (
	"testing"

	"github.com/smtdfc/nagare/core/configuration"
	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	config := &configuration.Configuration{Host: "127.0.0.1", Port: "0"}
	log := logger.NewLogger()

	app := ProvideFiberApp()
	server := NewServer(app, config, log)

	require.NotNil(t, server)
	assert.Same(t, app, server.app)
	assert.Same(t, config, server.config)
	assert.Same(t, log, server.logger)
}

func TestNewServer_NilDependencies(t *testing.T) {
	server := NewServer(nil, nil, nil)

	require.NotNil(t, server)
	assert.Nil(t, server.app)
	assert.Nil(t, server.config)
	assert.Nil(t, server.logger)
}
