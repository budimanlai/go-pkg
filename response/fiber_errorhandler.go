package response

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	switch code {
	case fiber.StatusNotFound:
		return NotFoundI18n(ctx, err.Error())
	case fiber.StatusBadRequest:
		return BadRequestI18n(ctx, err.Error(), nil)
	default:
		return ErrorI18n(ctx, code, err.Error(), nil)
	}
}
