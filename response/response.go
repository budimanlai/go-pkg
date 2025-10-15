package response

import (
	"github.com/budimanlai/go-pkg/i18n"
	"github.com/gofiber/fiber/v2"
)

var (
	i18nManager *i18n.I18nManager
)

func SetI18nManager(manager *i18n.I18nManager) {
	i18nManager = manager
}

func getLanguageFromContext(c *fiber.Ctx) string {
	if lang, ok := c.Locals("language").(string); ok {
		return lang
	}

	return i18nManager.DefaultLanguage // fallback to default language
}

func NotFoundI18n(c *fiber.Ctx, messageID string) error {
	if i18nManager == nil {
		return NotFound(c, messageID)
	}
	message := i18nManager.Translate(getLanguageFromContext(c), messageID, nil)
	return NotFound(c, message)
}

func ErrorI18n(c *fiber.Ctx, code int, messageID string, data interface{}) error {
	if i18nManager == nil {
		return Error(c, code, messageID)
	}
	message := i18nManager.Translate(getLanguageFromContext(c), messageID, data)
	return Error(c, code, message)

}

func BadRequestI18n(c *fiber.Ctx, messageID string, data interface{}) error {
	if i18nManager == nil {
		return BadRequest(c, messageID)
	}
	message := i18nManager.Translate(getLanguageFromContext(c), messageID, data)
	return BadRequest(c, message)
}

func SuccessI18n(c *fiber.Ctx, messageID string, data interface{}) error {
	if i18nManager == nil {
		return Success(c, messageID, data)
	}
	message := i18nManager.Translate(getLanguageFromContext(c), messageID, nil)
	return Success(c, message, data)
}

func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}

func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"meta": fiber.Map{
			"success": false,
			"message": message,
		},
		"data": nil,
	})
}

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"meta": fiber.Map{
			"success": false,
			"message": message,
		},
		"data": nil,
	})
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"meta": fiber.Map{
			"success": true,
			"message": message,
		},
		"data": data,
	})
}
