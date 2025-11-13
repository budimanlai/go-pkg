package i18n

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// I18nConfig holds the configuration for internationalization (i18n) setup.
// It defines the supported languages, default language, locale file paths, and optional modules.
//
// Fields:
//   - DefaultLanguage: The default language tag used as fallback
//   - SupportedLangs: List of supported language codes (e.g., ["en", "id", "zh"])
//   - LocalesPath: Path to the directory containing locale files (default: "locales")
//   - Modules: Optional list of module names for modular locale files
//
// Example:
//
//	config := I18nConfig{
//	    DefaultLanguage: language.English,
//	    SupportedLangs:  []string{"en", "id", "zh"},
//	    LocalesPath:     "locales",
//	    Modules:         []string{"auth", "user", "product"},
//	}
type I18nConfig struct {
	DefaultLanguage language.Tag
	SupportedLangs  []string
	LocalesPath     string
	Modules         []string
}

// I18nManager manages internationalization operations including translation bundles and localizers.
// It maintains a cache of localizers for each language to improve performance.
//
// Fields:
//   - Bundle: The i18n bundle containing all loaded message files
//   - Localizer: Map of language codes to their respective localizer instances
//   - DefaultLanguage: The default language code as string
type I18nManager struct {
	Bundle          *i18n.Bundle
	Localizer       map[string]*i18n.Localizer
	DefaultLanguage string
}

// NewI18nManagerWithFiber creates a new I18nManager and automatically registers
// the I18nMiddleware with the provided Fiber application. This is a convenience
// function that combines i18n initialization and middleware setup in one call.
//
// Parameters:
//   - app: *fiber.App - The Fiber application instance
//   - i18nConfig: I18nConfig - Configuration for i18n setup
//
// Returns:
//   - *I18nManager: Initialized I18nManager instance
//   - error: Error if initialization fails
//
// Example:
//
//	app := fiber.New()
//	config := I18nConfig{
//	    DefaultLanguage: language.English,
//	    SupportedLangs:  []string{"en", "id"},
//	    LocalesPath:     "locales",
//	}
//	i18nManager, err := NewI18nManagerWithFiber(app, config)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewI18nManagerWithFiber(app *fiber.App, i18nConfig I18nConfig) (*I18nManager, error) {
	i18nManager, err := NewI18nManager(i18nConfig)
	if err != nil {
		return nil, errors.New("failed to initialize i18n")
	}

	// Add i18n middleware
	app.Use(I18nMiddleware(i18nConfig))
	return i18nManager, nil
}

// NewI18nManager creates and initializes a new I18nManager with the provided configuration.
// It loads all locale files from the specified path and sets up the translation bundle.
//
// The function supports two loading modes:
// 1. Flat structure (no modules): Loads files from locales/{lang}.json
// 2. Modular structure: Loads files from locales/{lang}/{module}.json
//
// Parameters:
//   - config: I18nConfig - Configuration with locale paths and supported languages
//
// Returns:
//   - *I18nManager: Initialized I18nManager instance
//   - error: Error if locale files cannot be loaded
//
// Example (Flat structure):
//
//	config := I18nConfig{
//	    DefaultLanguage: language.English,
//	    SupportedLangs:  []string{"en", "id"},
//	    LocalesPath:     "locales",
//	}
//	manager, err := NewI18nManager(config)
//
// Example (Modular structure):
//
//	config := I18nConfig{
//	    DefaultLanguage: language.English,
//	    SupportedLangs:  []string{"en", "id"},
//	    LocalesPath:     "locales",
//	    Modules:         []string{"auth", "user"},
//	}
//	manager, err := NewI18nManager(config)
func NewI18nManager(config I18nConfig) (*I18nManager, error) {
	bundle := i18n.NewBundle(config.DefaultLanguage)

	if config.LocalesPath == "" {
		config.LocalesPath = "locales"
	}

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	if len(config.Modules) == 0 {
		for _, lang := range config.SupportedLangs {
			bundle.MustLoadMessageFile(fmt.Sprintf("%s/%s.json", config.LocalesPath, lang))
		}
	} else {
		for _, lang := range config.SupportedLangs {
			for _, module := range config.Modules {
				bundle.MustLoadMessageFile(fmt.Sprintf("%s/%s/%s.json", config.LocalesPath, lang, module))
			}
		}
	}

	return &I18nManager{
		Bundle:          bundle,
		Localizer:       make(map[string]*i18n.Localizer),
		DefaultLanguage: config.DefaultLanguage.String(),
	}, nil
}

// TranslateWithConfig translates a message using the provided LocalizeConfig.
// It supports template data for dynamic message interpolation and falls back to
// the default language if translation is not found in the requested language.
//
// The function caches localizers for each language to improve performance on subsequent calls.
//
// Parameters:
//   - lang: Language code for translation (e.g., "en", "id", "zh")
//   - c: *i18n.LocalizeConfig - Configuration containing message ID and template data
//
// Returns:
//   - string: Translated message, or error message if translation is missing
//
// Example:
//
//	translation := manager.TranslateWithConfig("id", &i18n.LocalizeConfig{
//	    MessageID: "welcome_message",
//	    TemplateData: map[string]string{
//	        "Name": "John",
//	    },
//	})
func (m *I18nManager) TranslateWithConfig(lang string, c *i18n.LocalizeConfig) string {
	localizer, ok := m.Localizer[lang]
	if !ok {
		// Fallback to default language if specific language not found
		m.Localizer[lang] = i18n.NewLocalizer(m.Bundle, lang)
		localizer = m.Localizer[lang]
	}
	localized, err := localizer.Localize(c)
	if err != nil {
		if m.DefaultLanguage != lang {
			// pakai bahasa default
			return m.TranslateWithConfig(m.DefaultLanguage, c)
		} else {
			// get message id
			var msgId string
			if c.MessageID == "" {
				msgId = c.DefaultMessage.ID
			} else {
				msgId = c.MessageID
			}
			return fmt.Sprintf("Missing translation for %s: %s", lang, msgId)
		}
	}
	return localized
}

// Translate is a convenience method for translating messages with optional template data.
// It wraps TranslateWithConfig with a simpler interface for common use cases.
//
// Parameters:
//   - lang: Language code for translation (e.g., "en", "id", "zh")
//   - messageID: The message identifier to translate
//   - template: Optional template data for message interpolation (can be nil)
//
// Returns:
//   - string: Translated message
//
// Example:
//
//	// Simple translation without template
//	msg := manager.Translate("id", "greeting", nil)
//
//	// Translation with template data
//	msg := manager.Translate("id", "user_registered", map[string]string{
//	    "Email": "user@example.com",
//	})
func (m *I18nManager) Translate(lang, messageID string, template interface{}) string {
	cfg := &i18n.LocalizeConfig{
		MessageID:      messageID,
		DefaultMessage: &i18n.Message{ID: messageID},
	}

	if template != nil {
		cfg.TemplateData = template
	}

	return m.TranslateWithConfig(lang, cfg)
}

// Test demonstrates usage of the I18nManager translation methods.
// This method serves as an example and can be used for testing i18n functionality.
// It shows both TranslateWithConfig and Translate method usage with template data.
//
// Example output:
//   - Tests translation with LocalizeConfig
//   - Tests translation with Translate method
//   - Logs results to console
func (m *I18nManager) Test() {
	// Implement test logic
	emailAlreadyExists := m.TranslateWithConfig("id", &i18n.LocalizeConfig{
		MessageID: "email_already_exists",
		TemplateData: map[string]string{
			"Email": "budiman.lai@gmail.com",
		},
	})

	handphoneAlreadyExists := m.Translate("id", "handphone_already_exists", map[string]string{
		"Handphone": "08123456789",
	})

	log.Info(emailAlreadyExists)
	log.Info(handphoneAlreadyExists)
}
