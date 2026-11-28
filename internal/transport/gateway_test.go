package transport

import (
	"testing"

	"github.com/smtdfc/nagare/core/configuration"
	httptransport "github.com/smtdfc/nagare/core/transport/http"
	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGateway(t *testing.T) {
	log := logger.NewLogger()
	server := httptransport.NewServer(
		httptransport.ProvideFiberApp(),
		&configuration.Configuration{Host: "127.0.0.1", Port: "0"},
		log,
	)

	gateway := NewGateway(server)

	require.NotNil(t, gateway)
	assert.Same(t, server, gateway.httpServer)
}

func TestNewGateway_NilServer(t *testing.T) {
	gateway := NewGateway(nil)

	require.NotNil(t, gateway)
	assert.Nil(t, gateway.httpServer)
}
