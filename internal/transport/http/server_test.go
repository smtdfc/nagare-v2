package http

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/core/configuration"
	apperrors "github.com/smtdfc/nagare/core/errors"
	"github.com/smtdfc/nagare/pkgs/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	config := &configuration.Configuration{Host: "127.0.0.1", Port: "0"}
	log := logger.NewLogger()

	app := ProvideFiberApp(config)
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

func TestHandleError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantMessage string
	}{
		{name: "app error", err: apperrors.ErrRequestBodyInvalid, wantStatus: fiber.StatusBadRequest, wantMessage: "Invalid request body"},
		{name: "fiber error", err: fiber.ErrUnauthorized, wantStatus: fiber.StatusUnauthorized, wantMessage: "Unauthorized"},
		{name: "deadline exceeded", err: context.DeadlineExceeded, wantStatus: fiber.StatusGatewayTimeout, wantMessage: "request timeout"},
		{name: "unknown error", err: errors.New("boom"), wantStatus: fiber.StatusInternalServerError, wantMessage: "Internal error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{ErrorHandler: handleError})
			app.Get("/", func(ctx fiber.Ctx) error { return tt.err })

			response, err := app.Test(httptest.NewRequest("GET", "/", nil))
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), tt.wantMessage)
		})
	}
}
