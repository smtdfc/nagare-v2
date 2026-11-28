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
func NewServer(config *configuration.Configuration, logger *logger.Logger) *Server {
	app := fiber.New()
	return &Server{
		app:    app,
		config: config,
		logger: logger,
	}
}
