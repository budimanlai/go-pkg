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

type I18nConfig struct {
	DefaultLanguage language.Tag
	SupportedLangs  []string
	LocalesPath     string
	Modules         []string
}

type I18nManager struct {
	Bundle          *i18n.Bundle
	Localizer       map[string]*i18n.Localizer
	DefaultLanguage string
}

func NewI18nManagerWithFiber(app *fiber.App, i18nConfig I18nConfig) (*I18nManager, error) {
	i18nManager, err := NewI18nManager(i18nConfig)
	if err != nil {
		return nil, errors.New("failed to initialize i18n")
	}

	// Add i18n middleware
	app.Use(I18nMiddleware(i18nConfig))
	return i18nManager, nil
}

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
