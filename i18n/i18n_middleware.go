package i18n

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// I18nMiddleware extracts language from request and sets it in context
func I18nMiddleware(config I18nConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to get language from various sources in order of priority
		lang := extractLanguage(c, config)

		// Set language in context for use in handlers
		c.Locals("language", lang)

		return c.Next()
	}
}

// extractLanguage extracts language from request in order of priority:
// 1. Query parameter ?lang=id
// 2. Header Accept-Language
// 3. Default language
func extractLanguage(c *fiber.Ctx, config I18nConfig) string {
	// 1. Check query parameter
	if lang := c.Query("lang"); lang != "" {
		if isSupported(lang, config.SupportedLangs) {
			return lang
		}
	}

	// 2. Check Accept-Language header
	acceptLang := c.Get("Accept-Language")
	if acceptLang != "" {
		// Parse Accept-Language header (simplified)
		langs := parseAcceptLanguage(acceptLang)
		for _, lang := range langs {
			if isSupported(lang, config.SupportedLangs) {
				return lang
			}
		}
	}

	// 3. Return default language
	return config.DefaultLanguage.String()
}

// parseAcceptLanguage parses Accept-Language header
func parseAcceptLanguage(header string) []string {
	var languages []string

	// Split by comma and extract language codes
	parts := strings.Split(header, ",")
	for _, part := range parts {
		// Remove quality values (e.g., en-US;q=0.9 -> en-US)
		lang := strings.TrimSpace(strings.Split(part, ";")[0])

		// Extract primary language (e.g., en-US -> en)
		if idx := strings.Index(lang, "-"); idx > 0 {
			lang = lang[:idx]
		}

		if lang != "" {
			languages = append(languages, lang)
		}
	}

	return languages
}

// isSupported checks if language is supported
func isSupported(lang string, supported []string) bool {
	for _, supportedLang := range supported {
		if lang == supportedLang {
			return true
		}
	}
	return false
}

// GetLanguage gets language from fiber context
func GetLanguage(c *fiber.Ctx) string {
	if lang, ok := c.Locals("language").(string); ok {
		return lang
	}
	return "en" // fallback to English
}
