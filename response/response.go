package response

import (
	"github.com/budimanlai/go-pkg/i18n"
	"github.com/gofiber/fiber/v2"
)

type PaginationResult struct {
	Data      any   `json:"data"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
}

var (
	// i18nManager holds the global I18nManager instance for response translations
	i18nManager *i18n.I18nManager
)

// SetI18nManager sets the global I18nManager instance to be used by i18n response functions.
// This function must be called during application initialization before using any i18n response methods.
//
// Parameters:
//   - manager: *i18n.I18nManager - The initialized i18n manager instance
//
// Example:
//
//	i18nMgr, _ := i18n.NewI18nManager(config)
//	response.SetI18nManager(i18nMgr)
func SetI18nManager(manager *i18n.I18nManager) {
	i18nManager = manager
}

// getLanguageFromContext retrieves the language code from the Fiber context.
// It attempts to get the language set by I18nMiddleware from context locals.
// If not found, it falls back to the default language from i18nManager.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//
// Returns:
//   - string: Language code (e.g., "en", "id", "zh")
func getLanguageFromContext(c *fiber.Ctx) string {
	if lang, ok := c.Locals("language").(string); ok {
		return lang
	}

	return i18nManager.DefaultLanguage // fallback to default language
}

// NotFoundI18n returns a 404 Not Found response with a translated message.
// The message is translated based on the language from the request context.
// If i18nManager is not set, it falls back to using the messageID as the message.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//   - messageID: Message identifier to translate
//
// Returns:
//   - error: Fiber error for response handling
//
// Example:
//
//	return response.NotFoundI18n(c, "user_not_found")
func NotFoundI18n(c *fiber.Ctx, messageID string) error {
	if i18nManager == nil {
		return NotFound(c, messageID)
	}
	message := i18nManager.Translate(getLanguageFromContext(c), messageID, nil)
	return NotFound(c, message)
}

// ErrorI18n returns an error response with a translated message and custom status code.
// The message is translated based on the language from the request context.
// Supports template data for dynamic message interpolation.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//   - code: HTTP status code
//   - messageID: Message identifier to translate
//   - data: Template data for message interpolation (can be nil)
//
// Returns:
//   - error: Fiber error for response handling
//
// Example:
//
//	return response.ErrorI18n(c, 500, "database_error", map[string]string{
//	    "Table": "users",
//	})
func ErrorI18n(c *fiber.Ctx, code int, messageID string, data interface{}) error {
	if i18nManager == nil {
		return Error(c, code, messageID)
	}
	message := i18nManager.Translate(getLanguageFromContext(c), messageID, data)
	return Error(c, code, message)

}

// BadRequestI18n returns a 400 Bad Request response with a translated message.
// The message is translated based on the language from the request context.
// Supports template data for dynamic message interpolation.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//   - messageID: Message identifier to translate
//   - data: Template data for message interpolation (can be nil)
//
// Returns:
//   - error: Fiber error for response handling
//
// Example:
//
//	return response.BadRequestI18n(c, "invalid_email", map[string]string{
//	    "Email": "invalid@",
//	})
func BadRequestI18n(c *fiber.Ctx, messageID string, data interface{}) error {
	if i18nManager == nil {
		return BadRequest(c, messageID)
	}
	message := i18nManager.Translate(getLanguageFromContext(c), messageID, data)
	return BadRequest(c, message)
}

// SuccessI18n returns a 200 OK response with a translated message and optional data.
// The message is translated based on the language from the request context.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//   - messageID: Message identifier to translate
//   - data: Response data to include in the response body (can be nil)
//
// Returns:
//   - error: Fiber error for response handling
//
// Example:
//
//	return response.SuccessI18n(c, "user_created", fiber.Map{
//	    "id": 123,
//	    "name": "John Doe",
//	})
func SuccessI18n(c *fiber.Ctx, messageID string, data interface{}) error {
	if i18nManager == nil {
		return Success(c, messageID, data)
	}
	message := i18nManager.Translate(getLanguageFromContext(c), messageID, nil)
	return Success(c, message, data)
}

func SuccessWithPaginationI18n(c *fiber.Ctx, messageID string, data PaginationResult) error {
	if i18nManager == nil {
		return Success(c, messageID, data)
	}
	message := i18nManager.Translate(getLanguageFromContext(c), messageID, nil)
	return SuccessWithPagination(c, message, data)
}

// ValidationErrorI18n returns a 400 Bad Request response with validation error details.
// It extracts field-specific errors from the ValidationError and formats them in a JSON response.
// If the error is not a ValidationError, it falls back to a generic bad request response.
//
// Response format:
//
//	{
//	  "meta": {
//	    "success": false,
//	    "message": "First validation error message",
//	    "errors": {
//	      "Email": ["Email is required", "Email must be valid"],
//	      "Password": ["Password must be at least 8 characters"]
//	    }
//	  },
//	  "data": null
//	}
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//   - err: error - The validation error (should be *validator.ValidationError)
//
// Returns:
//   - error: Fiber error for response handling
//
// Example:
//
//	if err := validator.ValidateStructWithContext(c, user); err != nil {
//	    return response.ValidationErrorI18n(c, err)
//	}
func ValidationErrorI18n(c *fiber.Ctx, err error) error {
	// Type assertion to get ValidationError
	type validationError interface {
		First() string
		GetFieldErrors() map[string][]string
	}

	if verr, ok := err.(validationError); ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"meta": fiber.Map{
				"success": false,
				"message": verr.First(),
				"errors":  verr.GetFieldErrors(),
			},
			"data": nil,
		})
	}

	// Fallback if not a validation error
	return BadRequest(c, err.Error())
}

// NotFound returns a 404 Not Found JSON response with the specified message.
//
// Response format:
//
//	{
//	  "meta": {
//	    "success": false,
//	    "message": "Resource not found"
//	  },
//	  "data": null
//	}
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//   - message: Error message to include in response
//
// Returns:
//   - error: Fiber error for response handling
//
// Example:
//
//	return response.NotFound(c, "User not found")
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}

// Error returns a JSON error response with the specified status code and message.
//
// Response format:
//
//	{
//	  "meta": {
//	    "success": false,
//	    "message": "Error message"
//	  },
//	  "data": null
//	}
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//   - code: HTTP status code (e.g., 400, 404, 500)
//   - message: Error message to include in response
//
// Returns:
//   - error: Fiber error for response handling
//
// Example:
//
//	return response.Error(c, 500, "Internal server error")
func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"meta": fiber.Map{
			"success": false,
			"message": message,
		},
		"data": nil,
	})
}

// BadRequest returns a 400 Bad Request JSON response with the specified message.
//
// Response format:
//
//	{
//	  "meta": {
//	    "success": false,
//	    "message": "Invalid request"
//	  },
//	  "data": null
//	}
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//   - message: Error message to include in response
//
// Returns:
//   - error: Fiber error for response handling
//
// Example:
//
//	return response.BadRequest(c, "Invalid email format")
func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"meta": fiber.Map{
			"success": false,
			"message": message,
		},
		"data": nil,
	})
}

// Success returns a 200 OK JSON response with the specified message and data.
//
// Response format:
//
//	{
//	  "meta": {
//	    "success": true,
//	    "message": "Success message"
//	  },
//	  "data": {
//	    // your data here
//	  }
//	}
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//   - message: Success message to include in response
//   - data: Response data (can be nil, struct, map, slice, etc.)
//
// Returns:
//   - error: Fiber error for response handling
//
// Example:
//
//	return response.Success(c, "User created successfully", fiber.Map{
//	    "id": 123,
//	    "name": "John Doe",
//	})
func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"meta": fiber.Map{
			"success": true,
			"message": message,
		},
		"data": data,
	})
}

func SuccessWithPagination(c *fiber.Ctx, message string, data PaginationResult) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"meta": fiber.Map{
			"success":    true,
			"message":    message,
			"total":      data.Total,
			"total_page": data.TotalPage,
			"page":       data.Page,
			"limit":      data.Limit,
		},
		"data": data.Data,
	})
}
