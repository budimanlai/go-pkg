package response

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	pkg_i18n "github.com/budimanlai/go-pkg/i18n"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

func TestSuccess(t *testing.T) {
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
}

func TestError(t *testing.T) {
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
}

func TestBadRequest(t *testing.T) {
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
}

func TestNotFound(t *testing.T) {
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
}

func TestSuccessI18n(t *testing.T) {
	// Setup i18n
	i18nConfig := pkg_i18n.I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id"},
		LocalesPath:     "/Users/budiman/Documents/development/my_github/go-pkg/locales",
	}

	i18nManager, err := pkg_i18n.NewI18nManager(i18nConfig)
	if err != nil {
		t.Fatal(err)
	}
	SetI18nManager(i18nManager)

	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
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
}

func TestErrorI18n(t *testing.T) {
	// Setup i18n
	i18nConfig := pkg_i18n.I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id"},
		LocalesPath:     "/Users/budiman/Documents/development/my_github/go-pkg/locales",
	}

	i18nManager, err := pkg_i18n.NewI18nManager(i18nConfig)
	if err != nil {
		t.Fatal(err)
	}
	SetI18nManager(i18nManager)

	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
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
}
