package auth

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestNewDefaultQueryStringAuth(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)
	if qsa == nil {
		t.Fatal("Expected NewDefaultQueryStringAuth to return a non-nil instance")
	}

	if qsa.config.KeyProvider == nil {
		t.Error("Expected KeyProvider to be set")
	}

	if qsa.config.ParamName != "access-token" {
		t.Errorf("Expected ParamName to be 'access-token', got '%s'", qsa.config.ParamName)
	}
}

func TestQueryStringAuth_GetParamName(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "api_key",
	}

	qsa := NewDefaultQueryStringAuth(config)
	paramName := qsa.GetParamName()

	if paramName != "api_key" {
		t.Errorf("Expected ParamName to be 'api_key', got '%s'", paramName)
	}
}

func TestQueryStringAuth_SetParamName(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)
	qsa.SetParamName("new_param")

	paramName := qsa.GetParamName()
	if paramName != "new_param" {
		t.Errorf("Expected ParamName to be 'new_param', got '%s'", paramName)
	}
}

func TestQueryStringAuth_Middleware_Success(t *testing.T) {
	// Setup key provider with valid keys
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-key-123")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with valid key in query string
	req := httptest.NewRequest("GET", "/test?access-token=valid-key-123", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Success" {
		t.Errorf("Expected body 'Success', got '%s'", string(body))
	}
}

func TestQueryStringAuth_Middleware_InvalidKey(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-key-123")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with invalid key
	req := httptest.NewRequest("GET", "/test?access-token=invalid-key", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestQueryStringAuth_Middleware_MissingKey(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-key-123")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request without key
	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestQueryStringAuth_Middleware_CustomParamName(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("custom-key-456")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "api_key",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with custom param name
	req := httptest.NewRequest("GET", "/test?api_key=custom-key-456", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestQueryStringAuth_Middleware_WrongParamName(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-key-123")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with wrong param name
	req := httptest.NewRequest("GET", "/test?wrong_param=valid-key-123", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 when using wrong param name, got %d", resp.StatusCode)
	}
}

func TestQueryStringAuth_Middleware_SuccessHandler(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-key-123")

	successHandlerCalled := false
	var capturedToken string

	successHandler := func(c *fiber.Ctx, token string) error {
		successHandlerCalled = true
		capturedToken = token
		c.Locals("user_id", "user-123")
		return nil
	}

	config := QueryStringAuthConfig{
		KeyProvider:    keyProvider,
		ParamName:      "access-token",
		SuccessHandler: &successHandler,
	}

	qsa := NewDefaultQueryStringAuth(config)

	var userID string

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		userID = c.Locals("user_id").(string)
		return c.SendString("Success")
	})

	// Create request with valid key
	req := httptest.NewRequest("GET", "/test?access-token=valid-key-123", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if !successHandlerCalled {
		t.Error("Expected success handler to be called")
	}

	if capturedToken != "valid-key-123" {
		t.Errorf("Expected token 'valid-key-123', got '%s'", capturedToken)
	}

	if userID != "user-123" {
		t.Errorf("Expected user_id 'user-123', got '%s'", userID)
	}
}

func TestQueryStringAuth_Middleware_SuccessHandlerError(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-key-123")

	successHandler := func(c *fiber.Ctx, token string) error {
		return errors.New("custom error from success handler")
	}

	config := QueryStringAuthConfig{
		KeyProvider:    keyProvider,
		ParamName:      "access-token",
		SuccessHandler: &successHandler,
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with valid key
	req := httptest.NewRequest("GET", "/test?access-token=valid-key-123", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	// Success handler error should result in unauthorized
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 when success handler returns error, got %d", resp.StatusCode)
	}
}

func TestQueryStringAuth_Middleware_CustomErrorHandler(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-key-123")

	errorHandlerCalled := false
	customErrorHandler := func(c *fiber.Ctx, err error) error {
		errorHandlerCalled = true
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Custom error message",
		})
	}

	config := QueryStringAuthConfig{
		KeyProvider:  keyProvider,
		ParamName:    "access-token",
		ErrorHandler: customErrorHandler,
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with invalid key
	req := httptest.NewRequest("GET", "/test?access-token=invalid-key", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusForbidden {
		t.Errorf("Expected status 403, got %d", resp.StatusCode)
	}

	if !errorHandlerCalled {
		t.Error("Expected custom error handler to be called")
	}
}

func TestQueryStringAuth_Middleware_MultipleKeys(t *testing.T) {
	// Setup key provider with multiple keys
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("key-1")
	keyProvider.Add("key-2")
	keyProvider.Add("key-3")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "token",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	testCases := []struct {
		key      string
		expected int
	}{
		{"key-1", fiber.StatusOK},
		{"key-2", fiber.StatusOK},
		{"key-3", fiber.StatusOK},
		{"invalid-key", fiber.StatusUnauthorized},
		{"", fiber.StatusUnauthorized},
	}

	for _, tc := range testCases {
		var url string
		if tc.key == "" {
			url = "/test"
		} else {
			url = "/test?token=" + tc.key
		}

		req := httptest.NewRequest("GET", url, nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request for key '%s': %v", tc.key, err)
		}

		if resp.StatusCode != tc.expected {
			t.Errorf("For key '%s', expected status %d, got %d", tc.key, tc.expected, resp.StatusCode)
		}
	}
}

func TestQueryStringAuth_Middleware_EmptyKey(t *testing.T) {
	// Setup key provider with empty key
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with empty key value
	// Note: Empty query param value is treated as missing by keyauth middleware
	req := httptest.NewRequest("GET", "/test?access-token=", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	// Empty query string value is treated as missing/invalid
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 for empty query param value, got %d", resp.StatusCode)
	}
}

func TestQueryStringAuth_Middleware_SpecialCharactersInKey(t *testing.T) {
	// Setup key provider with special characters (URL safe)
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("key-with-dashes_and_underscores.dots")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with URL-safe special characters in key
	req := httptest.NewRequest("GET", "/test?access-token=key-with-dashes_and_underscores.dots", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestQueryStringAuth_Middleware_UpdateKeys(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("old-key")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Test with old key
	req := httptest.NewRequest("GET", "/test?access-token=old-key", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200 with old key, got %d", resp.StatusCode)
	}

	// Remove old key and add new key
	keyProvider.Remove("old-key")
	keyProvider.Add("new-key")

	// Test with old key (should fail)
	req = httptest.NewRequest("GET", "/test?access-token=old-key", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 with removed key, got %d", resp.StatusCode)
	}

	// Test with new key (should succeed)
	req = httptest.NewRequest("GET", "/test?access-token=new-key", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200 with new key, got %d", resp.StatusCode)
	}
}

func TestQueryStringAuth_Middleware_MultipleQueryParams(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-key-123")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with multiple query parameters
	req := httptest.NewRequest("GET", "/test?page=1&limit=10&access-token=valid-key-123&sort=desc", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200 with multiple query params, got %d", resp.StatusCode)
	}
}

func TestQueryStringAuth_Middleware_DifferentHTTPMethods(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-key-123")

	config := QueryStringAuthConfig{
		KeyProvider: keyProvider,
		ParamName:   "access-token",
	}

	qsa := NewDefaultQueryStringAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(qsa.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("GET Success")
	})
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("POST Success")
	})
	app.Put("/test", func(c *fiber.Ctx) error {
		return c.SendString("PUT Success")
	})
	app.Delete("/test", func(c *fiber.Ctx) error {
		return c.SendString("DELETE Success")
	})

	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/test?access-token=valid-key-123", nil)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make %s request: %v", method, err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200 for %s method, got %d", method, resp.StatusCode)
		}
	}
}
