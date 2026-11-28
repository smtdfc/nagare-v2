package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/core/transport/http/controllers"
)

type RouteInitializer struct{}

// @Injectable
func NewRouterInitializer(app *fiber.App, modelProviderController *controllers.ModelProviderController) *RouteInitializer {
	return &RouteInitializer{}
}
