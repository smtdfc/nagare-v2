package utils

import "github.com/gofiber/fiber/v3"

func ParseBody[T any](ctx fiber.Ctx) (*T, error) {
	var payload T
	if err := ctx.Bind().Body(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
