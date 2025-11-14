package i18n

import (
	"strings"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// ============================================================================
// I18nManager Creation Tests
// ============================================================================

func TestNewI18nManager(t *testing.T) {
	t.Run("basic_initialization", func(t *testing.T) {
		config := I18nConfig{
			DefaultLanguage: language.English,
			SupportedLangs:  []string{"en", "id"},
			LocalesPath:     "../locales",
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
	})

	t.Run("with_chinese_support", func(t *testing.T) {
		config := I18nConfig{
			DefaultLanguage: language.English,
			SupportedLangs:  []string{"en", "id", "zh"},
			LocalesPath:     "../locales",
		}

		manager, err := NewI18nManager(config)
		if err != nil {
			t.Fatalf("Failed to create I18nManager with Chinese: %v", err)
		}

		if manager.Bundle == nil {
			t.Error("Bundle should not be nil")
		}
	})

	t.Run("default_locales_path", func(t *testing.T) {
		// Test with empty LocalesPath - should default to "locales"
		// But since we don't have locales/ in root, this should panic
		config := I18nConfig{
			DefaultLanguage: language.English,
			SupportedLangs:  []string{"en", "id"},
			LocalesPath:     "",
		}

		// We expect this to panic since locales/ doesn't exist in project root
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when locale files not found")
			}
		}()

		NewI18nManager(config)
	})

	t.Run("indonesian_default_language", func(t *testing.T) {
		config := I18nConfig{
			DefaultLanguage: language.Indonesian,
			SupportedLangs:  []string{"en", "id"},
			LocalesPath:     "../locales",
		}

		manager, err := NewI18nManager(config)
		if err != nil {
			t.Fatalf("Failed to create I18nManager with Indonesian default: %v", err)
		}

		if manager.DefaultLanguage != "id" {
			t.Errorf("Expected default language 'id', got '%s'", manager.DefaultLanguage)
		}
	})
}

// ============================================================================
// Translation Tests
// ============================================================================

func TestTranslate(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id", "zh"},
		LocalesPath:     "../locales",
	}

	manager, err := NewI18nManager(config)
	if err != nil {
		t.Fatalf("Failed to create I18nManager: %v", err)
	}

	t.Run("simple_translation_english", func(t *testing.T) {
		result := manager.Translate("en", "welcome", nil)
		expected := "Welcome to our application!"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("simple_translation_indonesian", func(t *testing.T) {
		result := manager.Translate("id", "welcome", nil)
		expected := "Selamat datang di aplikasi kami!"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("simple_translation_chinese", func(t *testing.T) {
		result := manager.Translate("zh", "welcome", nil)
		expected := "欢迎"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("template_english", func(t *testing.T) {
		result := manager.Translate("en", "hello_name", map[string]interface{}{"Name": "Test"})
		expected := "Hello, Test!"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("template_indonesian", func(t *testing.T) {
		result := manager.Translate("id", "hello_name", map[string]interface{}{"Name": "Test"})
		expected := "Halo, Test!"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("template_chinese", func(t *testing.T) {
		result := manager.Translate("zh", "hello_name", map[string]interface{}{"Name": "Test"})
		expected := "你好，Test!"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("validator_messages", func(t *testing.T) {
		// Test validator messages added in v1.0.1
		result := manager.Translate("en", "validator.required", nil)
		if result == "" {
			t.Error("Expected validator.required translation")
		}

		result = manager.Translate("id", "validator.email", nil)
		if result == "" {
			t.Error("Expected validator.email translation")
		}
	})

	t.Run("nil_template", func(t *testing.T) {
		result := manager.Translate("en", "welcome", nil)
		if result == "" {
			t.Error("Should return translation even with nil template")
		}
	})
}

// ============================================================================
// Fallback Tests
// ============================================================================

func TestTranslateFallback(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id"},
		LocalesPath:     "../locales",
	}

	manager, err := NewI18nManager(config)
	if err != nil {
		t.Fatalf("Failed to create I18nManager: %v", err)
	}

	t.Run("missing_key_in_default_language", func(t *testing.T) {
		result := manager.Translate("en", "selamat_pagi", nil)
		expected := "Missing translation for en: selamat_pagi"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("existing_key_in_specific_language", func(t *testing.T) {
		result := manager.Translate("id", "selamat_pagi", nil)
		expected := "Selamat pagi"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("missing_key_fallback_to_default", func(t *testing.T) {
		// Test with a key that doesn't exist in any language
		result := manager.Translate("id", "nonexistent_key_xyz", nil)
		// Should fallback to default (English) which also doesn't have it
		if result == "" {
			t.Error("Should return missing translation message")
		}
		// Should contain missing translation message
		if !strings.Contains(result, "Missing translation") {
			t.Errorf("Expected missing translation message, got '%s'", result)
		}
	})

	t.Run("completely_missing_key", func(t *testing.T) {
		result := manager.Translate("id", "nonexistent_key_12345", nil)
		// Should fallback to default and return missing message
		if result == "" {
			t.Error("Should return missing translation message")
		}
	})
}

// ============================================================================
// TranslateWithConfig Tests
// ============================================================================

func TestTranslateWithConfig(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id", "zh"},
		LocalesPath:     "../locales",
	}

	manager, err := NewI18nManager(config)
	if err != nil {
		t.Fatalf("Failed to create I18nManager: %v", err)
	}

	t.Run("with_message_id", func(t *testing.T) {
		result := manager.TranslateWithConfig("en", &i18n.LocalizeConfig{
			MessageID: "welcome",
		})
		expected := "Welcome to our application!"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("with_template_data", func(t *testing.T) {
		result := manager.TranslateWithConfig("id", &i18n.LocalizeConfig{
			MessageID: "hello_name",
			TemplateData: map[string]interface{}{
				"Name": "Budiman",
			},
		})
		expected := "Halo, Budiman!"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("with_default_message", func(t *testing.T) {
		result := manager.TranslateWithConfig("en", &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "welcome",
			},
		})
		expected := "Welcome to our application!"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("missing_translation_with_default_message", func(t *testing.T) {
		result := manager.TranslateWithConfig("en", &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "nonexistent_message",
			},
		})
		// Should return missing translation error
		if result == "" {
			t.Error("Should return missing translation message")
		}
	})

	t.Run("chinese_with_template", func(t *testing.T) {
		result := manager.TranslateWithConfig("zh", &i18n.LocalizeConfig{
			MessageID: "hello_name",
			TemplateData: map[string]interface{}{
				"Name": "张三",
			},
		})
		expected := "你好，张三!"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}

// ============================================================================
// Localizer Caching Tests
// ============================================================================

func TestLocalizerCaching(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id"},
		LocalesPath:     "../locales",
	}

	manager, err := NewI18nManager(config)
	if err != nil {
		t.Fatalf("Failed to create I18nManager: %v", err)
	}

	t.Run("localizer_created_on_first_use", func(t *testing.T) {
		// Initially, localizer map should be empty
		if len(manager.Localizer) != 0 {
			t.Error("Localizer map should be empty initially")
		}

		// First translation should create localizer
		manager.Translate("en", "welcome", nil)

		if len(manager.Localizer) != 1 {
			t.Errorf("Expected 1 localizer, got %d", len(manager.Localizer))
		}

		if manager.Localizer["en"] == nil {
			t.Error("English localizer should be created")
		}
	})

	t.Run("localizer_reused_on_subsequent_calls", func(t *testing.T) {
		manager.Translate("en", "welcome", nil)
		firstLocalizer := manager.Localizer["en"]

		manager.Translate("en", "hello_name", map[string]interface{}{"Name": "Test"})
		secondLocalizer := manager.Localizer["en"]

		if firstLocalizer != secondLocalizer {
			t.Error("Localizer should be reused, not recreated")
		}
	})

	t.Run("multiple_language_localizers", func(t *testing.T) {
		manager.Translate("en", "welcome", nil)
		manager.Translate("id", "welcome", nil)

		if len(manager.Localizer) != 2 {
			t.Errorf("Expected 2 localizers, got %d", len(manager.Localizer))
		}

		if manager.Localizer["en"] == nil {
			t.Error("English localizer should exist")
		}

		if manager.Localizer["id"] == nil {
			t.Error("Indonesian localizer should exist")
		}
	})
}
