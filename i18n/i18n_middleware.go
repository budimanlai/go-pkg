package i18n

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// I18nMiddleware creates a Fiber middleware handler that extracts the language preference
// from incoming requests and stores it in the context for use in downstream handlers.
// The language is extracted from multiple sources with a priority order:
// 1. Query parameter (?lang=id)
// 2. Accept-Language HTTP header
// 3. Default language from config
//
// Parameters:
//   - config: I18nConfig containing default language and supported languages list
//
// Returns:
//   - fiber.Handler: Middleware function to be used with Fiber app
//
// Example:
//
//	app := fiber.New()
//	config := I18nConfig{
//	    DefaultLanguage: English,
//	    SupportedLangs:  []string{"en", "id", "zh"},
//	}
//	app.Use(I18nMiddleware(config))
func I18nMiddleware(config I18nConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to get language from various sources in order of priority
		lang := extractLanguage(c, config)

		// Set language in context for use in handlers
		c.Locals("language", lang)

		return c.Next()
	}
}

// extractLanguage extracts the preferred language from an HTTP request following a priority order:
// 1. Query parameter ?lang=id (highest priority)
// 2. Accept-Language HTTP header
// 3. Default language from configuration (fallback)
//
// Only languages listed in config.SupportedLangs are accepted. If the requested
// language is not supported, it falls back to the next source in the priority chain.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context containing the HTTP request
//   - config: I18nConfig - Configuration with default and supported languages
//
// Returns:
//   - string: The selected language code (e.g., "en", "id", "zh")
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

// parseAcceptLanguage parses the Accept-Language HTTP header and extracts language codes.
// It handles quality values (e.g., en-US;q=0.9) by removing them and extracts the
// primary language code from locale-specific tags (e.g., en-US -> en).
//
// Parameters:
//   - header: Accept-Language header value (e.g., "en-US,en;q=0.9,id;q=0.8")
//
// Returns:
//   - []string: Ordered list of language codes extracted from the header
//
// Example:
//
//	langs := parseAcceptLanguage("en-US,en;q=0.9,id;q=0.8")
//	// Returns: ["en", "en", "id"]
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

// isSupported checks if a given language code is in the list of supported languages.
//
// Parameters:
//   - lang: Language code to check (e.g., "en", "id", "zh")
//   - supported: Slice of supported language codes
//
// Returns:
//   - bool: true if the language is supported, false otherwise
//
// Example:
//
//	supported := []string{"en", "id", "zh"}
//	isSupported("id", supported)  // Returns: true
//	isSupported("fr", supported)  // Returns: false
func isSupported(lang string, supported []string) bool {
	for _, supportedLang := range supported {
		if lang == supportedLang {
			return true
		}
	}
	return false
}

// GetLanguage retrieves the language code stored in the Fiber context by I18nMiddleware.
// This function should be called in handlers after I18nMiddleware has been applied.
// If no language is found in context, it returns "en" as a fallback.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context containing the stored language
//
// Returns:
//   - string: The language code (e.g., "en", "id", "zh"), defaults to "en" if not found
//
// Example:
//
//	app.Get("/hello", func(c *fiber.Ctx) error {
//	    lang := GetLanguage(c)
//	    message := translate(lang, "greeting")
//	    return c.SendString(message)
//	})
func GetLanguage(c *fiber.Ctx) string {
	if lang, ok := c.Locals("language").(string); ok {
		return lang
	}
	return "en" // fallback to English
}
