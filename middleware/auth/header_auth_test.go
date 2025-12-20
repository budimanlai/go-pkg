package auth

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestNewHeaderAuth(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
		HeaderName:  "X-API-Key",
	}

	headerAuth := NewHeaderAuth(config)
	if headerAuth == nil {
		t.Fatal("Expected NewHeaderAuth to return a non-nil instance")
	}

	if headerAuth.config.KeyProvider == nil {
		t.Error("Expected KeyProvider to be set")
	}

	if headerAuth.config.HeaderName != "X-API-Key" {
		t.Errorf("Expected HeaderName to be 'X-API-Key', got '%s'", headerAuth.config.HeaderName)
	}
}

func TestNewHeaderAuth_DefaultHeaderName(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
	}

	headerAuth := NewHeaderAuth(config)
	if headerAuth.config.HeaderName != "X-API-Key" {
		t.Errorf("Expected default HeaderName to be 'X-API-Key', got '%s'", headerAuth.config.HeaderName)
	}
}

func TestHeaderAuth_GetHeaderName(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
		HeaderName:  "X-Custom-Key",
	}

	headerAuth := NewHeaderAuth(config)
	headerName := headerAuth.GetHeaderName()

	if headerName != "X-Custom-Key" {
		t.Errorf("Expected HeaderName to be 'X-Custom-Key', got '%s'", headerName)
	}
}

func TestHeaderAuth_SetHeaderName(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
		HeaderName:  "X-API-Key",
	}

	headerAuth := NewHeaderAuth(config)
	headerAuth.SetHeaderName("X-New-Key")

	headerName := headerAuth.GetHeaderName()
	if headerName != "X-New-Key" {
		t.Errorf("Expected HeaderName to be 'X-New-Key', got '%s'", headerName)
	}
}

func TestHeaderAuth_Middleware_Success(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-api-key-123")

	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
		HeaderName:  "X-API-Key",
	}

	headerAuth := NewHeaderAuth(config)

	app := fiber.New()
	app.Use(headerAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "valid-api-key-123")

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

func TestHeaderAuth_Middleware_InvalidKey(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-api-key-123")

	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
		HeaderName:  "X-API-Key",
	}

	headerAuth := NewHeaderAuth(config)

	app := fiber.New()
	app.Use(headerAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "invalid-key")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestHeaderAuth_Middleware_MissingKey(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-api-key-123")

	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
		HeaderName:  "X-API-Key",
	}

	headerAuth := NewHeaderAuth(config)

	app := fiber.New()
	app.Use(headerAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestHeaderAuth_Middleware_CustomHeaderName(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("custom-key-456")

	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
		HeaderName:  "X-Custom-API-Key",
	}

	headerAuth := NewHeaderAuth(config)

	app := fiber.New()
	app.Use(headerAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Custom-API-Key", "custom-key-456")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHeaderAuth_Middleware_SuccessHandler(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-api-key-123")

	successHandlerCalled := false
	var capturedToken string

	successHandler := func(c *fiber.Ctx, token string) error {
		successHandlerCalled = true
		capturedToken = token
		c.Locals("user_id", "user-123")
		return nil
	}

	config := HeaderAuthConfig{
		KeyProvider:    keyProvider,
		HeaderName:     "X-API-Key",
		SuccessHandler: &successHandler,
	}

	headerAuth := NewHeaderAuth(config)

	var userID string

	app := fiber.New()
	app.Use(headerAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		userID = c.Locals("user_id").(string)
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "valid-api-key-123")

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

	if capturedToken != "valid-api-key-123" {
		t.Errorf("Expected token 'valid-api-key-123', got '%s'", capturedToken)
	}

	if userID != "user-123" {
		t.Errorf("Expected user_id 'user-123', got '%s'", userID)
	}
}

func TestHeaderAuth_Middleware_SuccessHandlerError(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-api-key-123")

	successHandler := func(c *fiber.Ctx, token string) error {
		return errors.New("custom error from success handler")
	}

	config := HeaderAuthConfig{
		KeyProvider:    keyProvider,
		HeaderName:     "X-API-Key",
		SuccessHandler: &successHandler,
	}

	headerAuth := NewHeaderAuth(config)

	app := fiber.New()
	app.Use(headerAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "valid-api-key-123")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 when success handler returns error, got %d", resp.StatusCode)
	}
}

func TestHeaderAuth_Middleware_CustomErrorHandler(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-api-key-123")

	errorHandlerCalled := false
	customErrorHandler := func(c *fiber.Ctx, err error) error {
		errorHandlerCalled = true
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Custom error message",
		})
	}

	config := HeaderAuthConfig{
		KeyProvider:  keyProvider,
		HeaderName:   "X-API-Key",
		ErrorHandler: customErrorHandler,
	}

	headerAuth := NewHeaderAuth(config)

	app := fiber.New()
	app.Use(headerAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "invalid-key")

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

func TestHeaderAuth_Middleware_MultipleKeys(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("key-1")
	keyProvider.Add("key-2")
	keyProvider.Add("key-3")

	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
		HeaderName:  "X-API-Key",
	}

	headerAuth := NewHeaderAuth(config)

	app := fiber.New()
	app.Use(headerAuth.Middleware())
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
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", tc.key)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request for key '%s': %v", tc.key, err)
		}

		if resp.StatusCode != tc.expected {
			t.Errorf("For key '%s', expected status %d, got %d", tc.key, tc.expected, resp.StatusCode)
		}
	}
}

func TestHeaderAuth_Middleware_CaseInsensitiveHeader(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-api-key-123")

	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
		HeaderName:  "X-API-Key",
	}

	headerAuth := NewHeaderAuth(config)

	app := fiber.New()
	app.Use(headerAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	testCases := []string{
		"X-API-Key",
		"x-api-key",
		"X-Api-Key",
	}

	for _, headerName := range testCases {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set(headerName, "valid-api-key-123")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request with header '%s': %v", headerName, err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200 for header '%s', got %d", headerName, resp.StatusCode)
		}
	}
}

func TestHeaderAuth_Middleware_DifferentHTTPMethods(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	keyProvider.Add("valid-api-key-123")

	config := HeaderAuthConfig{
		KeyProvider: keyProvider,
		HeaderName:  "X-API-Key",
	}

	headerAuth := NewHeaderAuth(config)

	app := fiber.New()
	app.Use(headerAuth.Middleware())
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
		req := httptest.NewRequest(method, "/test", nil)
		req.Header.Set("X-API-Key", "valid-api-key-123")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make %s request: %v", method, err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200 for %s method, got %d", method, resp.StatusCode)
		}
	}
}
