package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v3/middleware/logger"
	fiberrecover "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/configuration"
	apperrors "github.com/smtdfc/nagare/core/errors"
	"github.com/smtdfc/nagare/core/transport/http/utils"
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

func handleError(ctx fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := apperrors.ErrInternal.Message

	var e *fiber.Error
	var appErr *apperrors.BaseError

	switch {
	case errors.As(err, &appErr):
		code = appErr.Code
		message = appErr.Message
	case errors.As(err, &e):
		code = e.Code
		message = e.Message
	case errors.Is(err, context.DeadlineExceeded):
		code = fiber.StatusGatewayTimeout
		message = "request timeout"
	}

	return utils.ResponseError(ctx, message, code)
}

// @Injectable
func ProvideFiberApp(config *configuration.Configuration) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: handleError,
	})

	app.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return uuid.NewString()
		},
	}))

	app.Use(fiberrecover.New())
	app.Use(fiberlogger.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.AllowOrigins,
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Minute,
	}))
	app.Use(compress.New())
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
