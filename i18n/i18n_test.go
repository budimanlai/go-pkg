package i18n

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestNewI18nManager(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id"},
		LocalesPath:     "/Users/budiman/Documents/development/my_github/go-pkg/locales",
	}

	manager, err := NewI18nManager(config)
	if err != nil {
		t.Fatalf("Failed to create I18nManager: %v", err)
	}

	if manager.Bundle == nil {
		t.Error("Bundle should not be nil")
	}

	if manager.Localizer == nil {
		t.Error("Localizer map should not be nil")
	}

	if manager.DefaultLanguage != "en" {
		t.Errorf("Expected default language 'en', got '%s'", manager.DefaultLanguage)
	}
}

func TestTranslate(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id"},
		LocalesPath:     "/Users/budiman/Documents/development/my_github/go-pkg/locales",
	}

	manager, err := NewI18nManager(config)
	if err != nil {
		t.Fatalf("Failed to create I18nManager: %v", err)
	}

	// Test existing translation
	result := manager.Translate("en", "welcome", nil)
	expected := "Welcome to our application!"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	result = manager.Translate("id", "welcome", nil)
	expected = "Selamat datang di aplikasi kami!"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test with template
	result = manager.Translate("en", "hello_name", map[string]interface{}{"Name": "Test"})
	expected = "Hello, Test!"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	result = manager.Translate("id", "hello_name", map[string]interface{}{"Name": "Test"})
	expected = "Halo, Test!"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestTranslateFallback(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id"},
		LocalesPath:     "/Users/budiman/Documents/development/my_github/go-pkg/locales",
	}

	manager, err := NewI18nManager(config)
	if err != nil {
		t.Fatalf("Failed to create I18nManager: %v", err)
	}

	// Test fallback to default language for missing translation
	result := manager.Translate("en", "selamat_pagi", nil)
	expected := "Missing translation for en: selamat_pagi"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// For id, it should find the translation
	result = manager.Translate("id", "selamat_pagi", nil)
	expected = "Selamat pagi"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestTranslateWithConfig(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id"},
		LocalesPath:     "/Users/budiman/Documents/development/my_github/go-pkg/locales",
	}

	manager, err := NewI18nManager(config)
	if err != nil {
		t.Fatalf("Failed to create I18nManager: %v", err)
	}

	// Test with LocalizeConfig
	result := manager.TranslateWithConfig("en", &i18n.LocalizeConfig{
		MessageID: "welcome",
	})
	expected := "Welcome to our application!"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
