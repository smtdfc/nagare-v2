package utils

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBody_InvalidJSON(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(ctx fiber.Ctx) error {
		payload, err := ParseBody[map[string]any](ctx)
		assert.Nil(t, payload)
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	})

	request := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func TestParseBody_Success(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(ctx fiber.Ctx) error {
		payload, err := ParseBody[map[string]any](ctx)
		require.NoError(t, err)
		assert.Equal(t, "nagare", (*payload)["name"])
		return ctx.SendStatus(fiber.StatusOK)
	})

	request := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"nagare"}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}
