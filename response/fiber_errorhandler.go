package response

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// FiberErrorHandler is a custom error handler for Fiber framework that processes errors
// and returns internationalized error responses.
//
// It handles different HTTP status codes and returns appropriate i18n error responses:
//   - 404 (Not Found): Returns NotFoundI18n response
//   - 400 (Bad Request): Returns BadRequestI18n response
//   - Other status codes: Returns ErrorI18n response with the corresponding status code
//
// If the error is a *fiber.Error, it uses the error's status code.
// Otherwise, it defaults to 500 (Internal Server Error).
//
// Parameters:
//   - ctx: The Fiber context containing the request/response data
//   - err: The error to be handled and formatted
//
// Returns:
//   - error: An internationalized error response based on the status code
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
