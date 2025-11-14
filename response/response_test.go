package response

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	pkg_i18n "github.com/budimanlai/go-pkg/i18n"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

// ============================================================================
// Basic Response Tests
// ============================================================================

func TestSuccess(t *testing.T) {
	t.Run("basic_success_response", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return Success(c, "Success message", map[string]string{"key": "value"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		meta, ok := result["meta"].(map[string]interface{})
		if !ok {
			t.Error("Meta not found")
		}

		if meta["success"] != true {
			t.Error("Success should be true")
		}

		if meta["message"] != "Success message" {
			t.Errorf("Expected message 'Success message', got %v", meta["message"])
		}

		data, ok := result["data"].(map[string]interface{})
		if !ok {
			t.Error("Data not found")
		}

		if data["key"] != "value" {
			t.Errorf("Expected data key 'value', got %v", data["key"])
		}
	})

	t.Run("success_with_nil_data", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return Success(c, "No data", nil)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		if result["data"] != nil {
			t.Error("Data should be nil")
		}
	})

	t.Run("success_with_array_data", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return Success(c, "List of items", []string{"item1", "item2", "item3"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		data, ok := result["data"].([]interface{})
		if !ok {
			t.Error("Data should be array")
		}

		if len(data) != 3 {
			t.Errorf("Expected 3 items, got %d", len(data))
		}
	})
}

func TestError(t *testing.T) {
	t.Run("basic_error_response", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return Error(c, 400, "Error message")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		meta, ok := result["meta"].(map[string]interface{})
		if !ok {
			t.Error("Meta not found")
		}

		if meta["success"] != false {
			t.Error("Success should be false")
		}

		if meta["message"] != "Error message" {
			t.Errorf("Expected message 'Error message', got %v", meta["message"])
		}

		if result["data"] != nil {
			t.Error("Data should be nil")
		}
	})

	t.Run("error_with_500_status", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return Error(c, 500, "Internal server error")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 500 {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	})

	t.Run("error_with_custom_status", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return Error(c, 403, "Forbidden")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 403 {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})
}

func TestBadRequest(t *testing.T) {
	t.Run("basic_bad_request", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return BadRequest(c, "Bad request message")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		meta, ok := result["meta"].(map[string]interface{})
		if !ok {
			t.Error("Meta not found")
		}

		if meta["success"] != false {
			t.Error("Success should be false")
		}

		if meta["message"] != "Bad request message" {
			t.Errorf("Expected message 'Bad request message', got %v", meta["message"])
		}
	})

	t.Run("bad_request_validation_message", func(t *testing.T) {
		app := fiber.New()

		app.Post("/test", func(c *fiber.Ctx) error {
			return BadRequest(c, "Email is required")
		})

		req := httptest.NewRequest("POST", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		if meta["message"] != "Email is required" {
			t.Error("Expected validation message")
		}
	})
}

func TestNotFound(t *testing.T) {
	t.Run("basic_not_found", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return NotFound(c, "Not found message")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 404 {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		meta, ok := result["meta"].(map[string]interface{})
		if !ok {
			t.Error("Meta not found")
		}

		if meta["success"] != false {
			t.Error("Success should be false")
		}

		if meta["message"] != "Not found message" {
			t.Errorf("Expected message 'Not found message', got %v", meta["message"])
		}
	})

	t.Run("resource_not_found", func(t *testing.T) {
		app := fiber.New()

		app.Get("/users/:id", func(c *fiber.Ctx) error {
			return NotFound(c, "User not found")
		})

		req := httptest.NewRequest("GET", "/users/999", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})
}

// ============================================================================
// I18n Response Tests
// ============================================================================

func setupI18n(t *testing.T) {
	i18nConfig := pkg_i18n.I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id", "zh"},
		LocalesPath:     "../locales",
	}

	i18nManager, err := pkg_i18n.NewI18nManager(i18nConfig)
	if err != nil {
		t.Fatal(err)
	}
	SetI18nManager(i18nManager)
}

func TestSuccessI18n(t *testing.T) {
	setupI18n(t)

	t.Run("success_with_english", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "en")
			return SuccessI18n(c, "welcome", map[string]string{"key": "value"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		meta, ok := result["meta"].(map[string]interface{})
		if !ok {
			t.Error("Meta not found")
		}

		if meta["success"] != true {
			t.Error("Success should be true")
		}

		expectedMessage := "Welcome to our application!"
		if meta["message"] != expectedMessage {
			t.Errorf("Expected message '%s', got %v", expectedMessage, meta["message"])
		}
	})

	t.Run("success_with_indonesian", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "id")
			return SuccessI18n(c, "welcome", nil)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		expectedMessage := "Selamat datang di aplikasi kami!"
		if meta["message"] != expectedMessage {
			t.Errorf("Expected Indonesian message, got %v", meta["message"])
		}
	})

	t.Run("success_with_chinese", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "zh")
			return SuccessI18n(c, "welcome", nil)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		expectedMessage := "欢迎"
		if meta["message"] != expectedMessage {
			t.Errorf("Expected Chinese message, got %v", meta["message"])
		}
	})
}

func TestErrorI18n(t *testing.T) {
	setupI18n(t)

	t.Run("error_with_english", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "en")
			return ErrorI18n(c, 500, "welcome", nil)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 500 {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}

		meta, ok := result["meta"].(map[string]interface{})
		if !ok {
			t.Error("Meta not found")
		}

		if meta["success"] != false {
			t.Error("Success should be false")
		}

		expectedMessage := "Welcome to our application!"
		if meta["message"] != expectedMessage {
			t.Errorf("Expected message '%s', got %v", expectedMessage, meta["message"])
		}
	})

	t.Run("error_with_indonesian", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "id")
			return ErrorI18n(c, 400, "welcome", nil)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		expectedMessage := "Selamat datang di aplikasi kami!"
		if meta["message"] != expectedMessage {
			t.Errorf("Expected Indonesian error message, got %v", meta["message"])
		}
	})

	t.Run("error_with_template_data", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "en")
			return ErrorI18n(c, 404, "hello_name", map[string]interface{}{"Name": "John"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		if meta["message"] != "Hello, John!" {
			t.Errorf("Expected templated message, got %v", meta["message"])
		}
	})
}

func TestBadRequestI18n(t *testing.T) {
	setupI18n(t)

	t.Run("bad_request_with_english", func(t *testing.T) {
		app := fiber.New()

		app.Post("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "en")
			return BadRequestI18n(c, "welcome", nil)
		})

		req := httptest.NewRequest("POST", "/test", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 400 {
			t.Errorf("Expected 400, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		expectedMessage := "Welcome to our application!"
		if meta["message"] != expectedMessage {
			t.Errorf("Expected message '%s', got %v", expectedMessage, meta["message"])
		}
	})

	t.Run("bad_request_with_template", func(t *testing.T) {
		app := fiber.New()

		app.Post("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "id")
			return BadRequestI18n(c, "hello_name", map[string]interface{}{"Name": "Budi"})
		})

		req := httptest.NewRequest("POST", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		if meta["message"] != "Halo, Budi!" {
			t.Errorf("Expected 'Halo, Budi!', got %v", meta["message"])
		}
	})
}

func TestNotFoundI18n(t *testing.T) {
	setupI18n(t)

	t.Run("not_found_with_english", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "en")
			return NotFoundI18n(c, "welcome")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		expectedMessage := "Welcome to our application!"
		if meta["message"] != expectedMessage {
			t.Errorf("Expected message '%s', got %v", expectedMessage, meta["message"])
		}
	})

	t.Run("not_found_with_chinese", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "zh")
			return NotFoundI18n(c, "welcome")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		expectedMessage := "欢迎"
		if meta["message"] != expectedMessage {
			t.Errorf("Expected Chinese message, got %v", meta["message"])
		}
	})
}

// ============================================================================
// ValidationErrorI18n Tests
// ============================================================================

// Mock ValidationError for testing
type mockValidationError struct {
	firstMsg    string
	fieldErrors map[string][]string
}

func (m *mockValidationError) Error() string {
	return m.firstMsg
}

func (m *mockValidationError) First() string {
	return m.firstMsg
}

func (m *mockValidationError) GetFieldErrors() map[string][]string {
	return m.fieldErrors
}

func TestValidationErrorI18n(t *testing.T) {
	t.Run("validation_error_with_field_errors", func(t *testing.T) {
		app := fiber.New()

		app.Post("/test", func(c *fiber.Ctx) error {
			mockErr := &mockValidationError{
				firstMsg: "Email is required",
				fieldErrors: map[string][]string{
					"email":    {"Email is required", "Email must be valid"},
					"password": {"Password must be at least 8 characters"},
				},
			}
			return ValidationErrorI18n(c, mockErr)
		})

		req := httptest.NewRequest("POST", "/test", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 400 {
			t.Errorf("Expected 400, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		if meta["message"] != "Email is required" {
			t.Errorf("Expected first error message, got %v", meta["message"])
		}

		errors, ok := meta["errors"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected errors field")
		}

		emailErrors, ok := errors["email"].([]interface{})
		if !ok {
			t.Fatal("Expected email errors array")
		}

		if len(emailErrors) != 2 {
			t.Errorf("Expected 2 email errors, got %d", len(emailErrors))
		}
	})

	t.Run("non_validation_error_fallback", func(t *testing.T) {
		app := fiber.New()

		app.Post("/test", func(c *fiber.Ctx) error {
			return ValidationErrorI18n(c, errors.New("Some generic error"))
		})

		req := httptest.NewRequest("POST", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		if meta["message"] != "Some generic error" {
			t.Errorf("Expected fallback error message, got %v", meta["message"])
		}

		if _, hasErrors := meta["errors"]; hasErrors {
			t.Error("Should not have errors field for generic error")
		}
	})

	t.Run("empty_field_errors", func(t *testing.T) {
		app := fiber.New()

		app.Post("/test", func(c *fiber.Ctx) error {
			mockErr := &mockValidationError{
				firstMsg:    "Validation failed",
				fieldErrors: map[string][]string{},
			}
			return ValidationErrorI18n(c, mockErr)
		})

		req := httptest.NewRequest("POST", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		errors, ok := meta["errors"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected errors field")
		}

		if len(errors) != 0 {
			t.Errorf("Expected empty errors map, got %d items", len(errors))
		}
	})
}

// ============================================================================
// FiberErrorHandler Tests
// ============================================================================

func TestFiberErrorHandler(t *testing.T) {
	setupI18n(t)

	t.Run("handle_404_error", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: FiberErrorHandler,
		})

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "en")
			return fiber.NewError(404, "Page not found")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("handle_400_error", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: FiberErrorHandler,
		})

		app.Post("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "en")
			return fiber.NewError(400, "Invalid request")
		})

		req := httptest.NewRequest("POST", "/test", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 400 {
			t.Errorf("Expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("handle_500_error", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: FiberErrorHandler,
		})

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "en")
			return fiber.NewError(500, "Internal server error")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 500 {
			t.Errorf("Expected 500, got %d", resp.StatusCode)
		}
	})

	t.Run("handle_generic_error", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: FiberErrorHandler,
		})

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "en")
			return errors.New("Generic error")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		// Generic errors default to 500
		if resp.StatusCode != 500 {
			t.Errorf("Expected 500 for generic error, got %d", resp.StatusCode)
		}
	})

	t.Run("handle_custom_status_code", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: FiberErrorHandler,
		})

		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "id")
			return fiber.NewError(403, "Forbidden")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 403 {
			t.Errorf("Expected 403, got %d", resp.StatusCode)
		}
	})
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestGetLanguageFromContext(t *testing.T) {
	setupI18n(t)

	t.Run("language_in_context", func(t *testing.T) {
		app := fiber.New()

		var capturedLang string
		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "zh")
			capturedLang = getLanguageFromContext(c)
			return c.SendString("ok")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		app.Test(req)

		if capturedLang != "zh" {
			t.Errorf("Expected 'zh', got '%s'", capturedLang)
		}
	})

	t.Run("fallback_to_default", func(t *testing.T) {
		app := fiber.New()

		var capturedLang string
		app.Get("/test", func(c *fiber.Ctx) error {
			capturedLang = getLanguageFromContext(c)
			return c.SendString("ok")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		app.Test(req)

		if capturedLang != "en" {
			t.Errorf("Expected default 'en', got '%s'", capturedLang)
		}
	})
}

func TestSetI18nManager(t *testing.T) {
	t.Run("set_i18n_manager", func(t *testing.T) {
		i18nConfig := pkg_i18n.I18nConfig{
			DefaultLanguage: language.Indonesian,
			SupportedLangs:  []string{"en", "id"},
			LocalesPath:     "../locales",
		}

		manager, _ := pkg_i18n.NewI18nManager(i18nConfig)
		SetI18nManager(manager)

		if i18nManager == nil {
			t.Error("i18nManager should not be nil after SetI18nManager")
		}

		if i18nManager.DefaultLanguage != "id" {
			t.Errorf("Expected default language 'id', got '%s'", i18nManager.DefaultLanguage)
		}
	})
}

// ============================================================================
// I18n Fallback Tests (when i18nManager is nil)
// ============================================================================

func TestI18nFallbackWhenManagerNil(t *testing.T) {
	// Save current manager and restore after test
	savedManager := i18nManager
	defer func() { i18nManager = savedManager }()

	// Set manager to nil
	i18nManager = nil

	t.Run("success_i18n_fallback", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return SuccessI18n(c, "raw_message", nil)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		// Should use messageID as message
		if meta["message"] != "raw_message" {
			t.Errorf("Expected 'raw_message', got %v", meta["message"])
		}
	})

	t.Run("error_i18n_fallback", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return ErrorI18n(c, 500, "error_key", nil)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		if meta["message"] != "error_key" {
			t.Errorf("Expected 'error_key', got %v", meta["message"])
		}
	})

	t.Run("bad_request_i18n_fallback", func(t *testing.T) {
		app := fiber.New()

		app.Post("/test", func(c *fiber.Ctx) error {
			return BadRequestI18n(c, "bad_request_key", nil)
		})

		req := httptest.NewRequest("POST", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		if meta["message"] != "bad_request_key" {
			t.Errorf("Expected 'bad_request_key', got %v", meta["message"])
		}
	})

	t.Run("not_found_i18n_fallback", func(t *testing.T) {
		app := fiber.New()

		app.Get("/test", func(c *fiber.Ctx) error {
			return NotFoundI18n(c, "not_found_key")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		meta := result["meta"].(map[string]interface{})

		if meta["message"] != "not_found_key" {
			t.Errorf("Expected 'not_found_key', got %v", meta["message"])
		}
	})
}
