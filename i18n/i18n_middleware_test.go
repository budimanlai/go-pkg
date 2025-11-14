package i18n

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

// ============================================================================
// I18nMiddleware Tests
// ============================================================================

func TestI18nMiddleware(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id", "zh"},
	}

	t.Run("language_from_query_parameter", func(t *testing.T) {
		app := fiber.New()
		app.Use(I18nMiddleware(config))
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := c.Locals("language").(string)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test?lang=id", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "id" {
			t.Errorf("Expected 'id', got '%s'", string(body))
		}
	})

	t.Run("language_from_accept_language_header", func(t *testing.T) {
		app := fiber.New()
		app.Use(I18nMiddleware(config))
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := c.Locals("language").(string)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Accept-Language", "id-ID,id;q=0.9")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "id" {
			t.Errorf("Expected 'id', got '%s'", string(body))
		}
	})

	t.Run("default_language_when_no_preference", func(t *testing.T) {
		app := fiber.New()
		app.Use(I18nMiddleware(config))
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := c.Locals("language").(string)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "en" {
			t.Errorf("Expected default 'en', got '%s'", string(body))
		}
	})

	t.Run("query_parameter_overrides_header", func(t *testing.T) {
		app := fiber.New()
		app.Use(I18nMiddleware(config))
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := c.Locals("language").(string)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test?lang=zh", nil)
		req.Header.Set("Accept-Language", "id-ID")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "zh" {
			t.Errorf("Expected 'zh' from query param, got '%s'", string(body))
		}
	})

	t.Run("unsupported_language_falls_back", func(t *testing.T) {
		app := fiber.New()
		app.Use(I18nMiddleware(config))
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := c.Locals("language").(string)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test?lang=fr", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "en" {
			t.Errorf("Expected fallback to 'en', got '%s'", string(body))
		}
	})

	t.Run("unsupported_header_falls_back_to_default", func(t *testing.T) {
		app := fiber.New()
		app.Use(I18nMiddleware(config))
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := c.Locals("language").(string)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Accept-Language", "fr-FR,de;q=0.9")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "en" {
			t.Errorf("Expected fallback to 'en', got '%s'", string(body))
		}
	})

	t.Run("chinese_language_support", func(t *testing.T) {
		app := fiber.New()
		app.Use(I18nMiddleware(config))
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := c.Locals("language").(string)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test?lang=zh", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "zh" {
			t.Errorf("Expected 'zh', got '%s'", string(body))
		}
	})

	t.Run("empty_query_parameter", func(t *testing.T) {
		app := fiber.New()
		app.Use(I18nMiddleware(config))
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := c.Locals("language").(string)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test?lang=", nil)
		req.Header.Set("Accept-Language", "id")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		body, _ := io.ReadAll(resp.Body)
		// Should fall back to header since query param is empty
		if string(body) != "id" {
			t.Errorf("Expected 'id' from header, got '%s'", string(body))
		}
	})
}

// ============================================================================
// extractLanguage Tests
// ============================================================================

func TestExtractLanguage(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id", "zh"},
	}

	t.Run("priority_query_param", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := extractLanguage(c, config)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test?lang=zh", nil)
		req.Header.Set("Accept-Language", "id")
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)

		if string(body) != "zh" {
			t.Errorf("Expected 'zh' from query param, got '%s'", string(body))
		}
	})

	t.Run("fallback_to_default", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := extractLanguage(c, config)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)

		if string(body) != "en" {
			t.Errorf("Expected default 'en', got '%s'", string(body))
		}
	})
}

// ============================================================================
// parseAcceptLanguage Tests
// ============================================================================

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected []string
	}{
		{
			name:     "simple_single_language",
			header:   "en",
			expected: []string{"en"},
		},
		{
			name:     "multiple_languages",
			header:   "en,id,zh",
			expected: []string{"en", "id", "zh"},
		},
		{
			name:     "with_quality_values",
			header:   "en-US,en;q=0.9,id;q=0.8",
			expected: []string{"en", "en", "id"},
		},
		{
			name:     "locale_specific",
			header:   "id-ID",
			expected: []string{"id"},
		},
		{
			name:     "complex_header",
			header:   "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
			expected: []string{"zh", "zh", "en", "en"},
		},
		{
			name:     "with_spaces",
			header:   "en-US, en;q=0.9, id;q=0.8",
			expected: []string{"en", "en", "id"},
		},
		{
			name:     "empty_header",
			header:   "",
			expected: []string{},
		},
		{
			name:     "single_locale",
			header:   "en-GB",
			expected: []string{"en"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAcceptLanguage(tt.header)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d languages, got %d", len(tt.expected), len(result))
				return
			}

			for i, lang := range result {
				if lang != tt.expected[i] {
					t.Errorf("Expected language[%d] = '%s', got '%s'", i, tt.expected[i], lang)
				}
			}
		})
	}
}

// ============================================================================
// isSupported Tests
// ============================================================================

func TestIsSupported(t *testing.T) {
	supported := []string{"en", "id", "zh"}

	tests := []struct {
		name     string
		lang     string
		expected bool
	}{
		{"english_supported", "en", true},
		{"indonesian_supported", "id", true},
		{"chinese_supported", "zh", true},
		{"french_not_supported", "fr", false},
		{"german_not_supported", "de", false},
		{"japanese_not_supported", "ja", false},
		{"empty_string", "", false},
		{"case_sensitive_EN", "EN", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSupported(tt.lang, supported)
			if result != tt.expected {
				t.Errorf("isSupported(%s) = %v, expected %v", tt.lang, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// GetLanguage Tests
// ============================================================================

func TestGetLanguage(t *testing.T) {
	t.Run("language_exists_in_context", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "id")
			lang := GetLanguage(c)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)

		if string(body) != "id" {
			t.Errorf("Expected 'id', got '%s'", string(body))
		}
	})

	t.Run("fallback_to_english_when_not_set", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			lang := GetLanguage(c)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)

		if string(body) != "en" {
			t.Errorf("Expected fallback to 'en', got '%s'", string(body))
		}
	})

	t.Run("chinese_language_in_context", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "zh")
			lang := GetLanguage(c)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)

		if string(body) != "zh" {
			t.Errorf("Expected 'zh', got '%s'", string(body))
		}
	})

	t.Run("wrong_type_in_context", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", 123) // Wrong type
			lang := GetLanguage(c)
			return c.SendString(lang)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)

		if string(body) != "en" {
			t.Errorf("Expected fallback to 'en', got '%s'", string(body))
		}
	})
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestI18nMiddlewareIntegration(t *testing.T) {
	config := I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id", "zh"},
	}

	t.Run("full_workflow_with_getlanguage", func(t *testing.T) {
		app := fiber.New()
		app.Use(I18nMiddleware(config))
		app.Get("/api/greeting", func(c *fiber.Ctx) error {
			lang := GetLanguage(c)
			greetings := map[string]string{
				"en": "Hello",
				"id": "Halo",
				"zh": "你好",
			}
			return c.SendString(greetings[lang])
		})

		// Test with query parameter
		req := httptest.NewRequest("GET", "/api/greeting?lang=id", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		if string(body) != "Halo" {
			t.Errorf("Expected 'Halo', got '%s'", string(body))
		}

		// Test with header
		req = httptest.NewRequest("GET", "/api/greeting", nil)
		req.Header.Set("Accept-Language", "zh-CN")
		resp, _ = app.Test(req)
		body, _ = io.ReadAll(resp.Body)
		if string(body) != "你好" {
			t.Errorf("Expected '你好', got '%s'", string(body))
		}

		// Test default
		req = httptest.NewRequest("GET", "/api/greeting", nil)
		resp, _ = app.Test(req)
		body, _ = io.ReadAll(resp.Body)
		if string(body) != "Hello" {
			t.Errorf("Expected 'Hello', got '%s'", string(body))
		}
	})

	t.Run("multiple_routes_same_middleware", func(t *testing.T) {
		app := fiber.New()
		app.Use(I18nMiddleware(config))

		app.Get("/route1", func(c *fiber.Ctx) error {
			return c.SendString(GetLanguage(c))
		})

		app.Get("/route2", func(c *fiber.Ctx) error {
			return c.SendString(GetLanguage(c))
		})

		req := httptest.NewRequest("GET", "/route1?lang=zh", nil)
		resp, _ := app.Test(req)
		body, _ := io.ReadAll(resp.Body)
		if string(body) != "zh" {
			t.Errorf("Route1: Expected 'zh', got '%s'", string(body))
		}

		req = httptest.NewRequest("GET", "/route2?lang=id", nil)
		resp, _ = app.Test(req)
		body, _ = io.ReadAll(resp.Body)
		if string(body) != "id" {
			t.Errorf("Route2: Expected 'id', got '%s'", string(body))
		}
	})
}
