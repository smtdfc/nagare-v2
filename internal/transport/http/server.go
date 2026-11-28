package http

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/core/configuration"
	"github.com/smtdfc/nagare/pkgs/logger"
)

type Server struct {
	app    *fiber.App
	config *configuration.Configuration
	logger *logger.Logger
}

func (s *Server) Listen() {
	err := s.app.Listen(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port))
	if err != nil {
		s.logger.Error("Failed to start server", "err", err)
	}
}

// @Injectable
func ProvideFiberApp() *fiber.App {
	app := fiber.New()
	return app
}

// @Injectable
func NewServer(app *fiber.App, config *configuration.Configuration, logger *logger.Logger) *Server {
	return &Server{
		app:    app,
		config: config,
		logger: logger,
	}
}
