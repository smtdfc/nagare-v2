package utils

import (
	"github.com/gofiber/fiber/v3"
	"github.com/smtdfc/nagare/pkgs/dtos/rest"
)

func ResponseError(ctx fiber.Ctx, message string, code int) error {
	return ctx.Status(code).JSON(&rest.ApiResponse[any]{
		IsSuccess: false,
		Data:      nil,
		Error: rest.ApiResponseError{
			Message: message,
		},
	})
}
