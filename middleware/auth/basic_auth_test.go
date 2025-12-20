package auth

import (
	"encoding/base64"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestNewBasicAuth(t *testing.T) {
	keyProvider := NewBaseKeyProvider()
	config := BasicAuthConfig{
		KeyProvider: keyProvider,
	}

	basicAuth := NewBasicAuth(config)
	if basicAuth == nil {
		t.Fatal("Expected NewBasicAuth to return a non-nil instance")
	}

	if basicAuth.config.KeyProvider == nil {
		t.Error("Expected KeyProvider to be set")
	}
}

func TestBasicAuth_Middleware_Success(t *testing.T) {
	// Setup key provider with test credentials
	keyProvider := NewBaseKeyProvider()
	keyProvider.AddKeyValue("admin", "secret123")

	config := BasicAuthConfig{
		KeyProvider: keyProvider,
	}

	basicAuth := NewBasicAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(basicAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with valid credentials
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:secret123")))

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

func TestBasicAuth_Middleware_InvalidPassword(t *testing.T) {
	// Setup key provider with test credentials
	keyProvider := NewBaseKeyProvider()
	keyProvider.AddKeyValue("admin", "secret123")

	config := BasicAuthConfig{
		KeyProvider: keyProvider,
	}

	basicAuth := NewBasicAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(basicAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with invalid password
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:wrongpassword")))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestBasicAuth_Middleware_UserNotFound(t *testing.T) {
	// Setup key provider without the user
	keyProvider := NewBaseKeyProvider()
	keyProvider.AddKeyValue("admin", "secret123")

	config := BasicAuthConfig{
		KeyProvider: keyProvider,
	}

	basicAuth := NewBasicAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(basicAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with non-existent user
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("nonexistent:password")))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestBasicAuth_Middleware_NoCredentials(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.AddKeyValue("admin", "secret123")

	config := BasicAuthConfig{
		KeyProvider: keyProvider,
	}

	basicAuth := NewBasicAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(basicAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request without credentials
	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestBasicAuth_Middleware_CustomUnauthorized(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.AddKeyValue("admin", "secret123")

	customUnauthorizedCalled := false
	config := BasicAuthConfig{
		KeyProvider: keyProvider,
		Unauthorized: func(c *fiber.Ctx) error {
			customUnauthorizedCalled = true
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Custom unauthorized message",
			})
		},
	}

	basicAuth := NewBasicAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(basicAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with invalid credentials
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:wrongpassword")))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusForbidden {
		t.Errorf("Expected status 403, got %d", resp.StatusCode)
	}

	if !customUnauthorizedCalled {
		t.Error("Expected custom unauthorized handler to be called")
	}
}

func TestBasicAuth_Middleware_ContextValues(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.AddKeyValue("testuser", "testpass")

	config := BasicAuthConfig{
		KeyProvider:     keyProvider,
		ContextUsername: "custom_username",
		ContextPassword: "custom_password",
	}

	basicAuth := NewBasicAuth(config)

	var capturedUsername, capturedPassword string

	// Create Fiber app
	app := fiber.New()
	app.Use(basicAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		capturedUsername = c.Locals("custom_username").(string)
		capturedPassword = c.Locals("custom_password").(string)
		return c.SendString("Success")
	})

	// Create request with valid credentials
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testuser:testpass")))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if capturedUsername != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", capturedUsername)
	}

	if capturedPassword != "testpass" {
		t.Errorf("Expected password 'testpass', got '%s'", capturedPassword)
	}
}

func TestBasicAuth_Middleware_MultipleUsers(t *testing.T) {
	// Setup key provider with multiple users
	keyProvider := NewBaseKeyProvider()
	keyProvider.AddKeyValue("admin", "admin123")
	keyProvider.AddKeyValue("user1", "pass1")
	keyProvider.AddKeyValue("user2", "pass2")

	config := BasicAuthConfig{
		KeyProvider: keyProvider,
	}

	basicAuth := NewBasicAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(basicAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	testCases := []struct {
		username string
		password string
		expected int
	}{
		{"admin", "admin123", fiber.StatusOK},
		{"user1", "pass1", fiber.StatusOK},
		{"user2", "pass2", fiber.StatusOK},
		{"admin", "wrongpass", fiber.StatusUnauthorized},
		{"user1", "wrongpass", fiber.StatusUnauthorized},
		{"nonexistent", "pass", fiber.StatusUnauthorized},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(tc.username+":"+tc.password)))

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request for %s: %v", tc.username, err)
		}

		if resp.StatusCode != tc.expected {
			t.Errorf("For user '%s' with password '%s', expected status %d, got %d",
				tc.username, tc.password, tc.expected, resp.StatusCode)
		}
	}
}

func TestBasicAuth_Middleware_EmptyCredentials(t *testing.T) {
	// Setup key provider with empty username/password
	keyProvider := NewBaseKeyProvider()
	keyProvider.AddKeyValue("", "")

	config := BasicAuthConfig{
		KeyProvider: keyProvider,
	}

	basicAuth := NewBasicAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(basicAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with empty credentials
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":")))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200 for empty credentials that exist, got %d", resp.StatusCode)
	}
}

func TestBasicAuth_Middleware_SpecialCharactersInPassword(t *testing.T) {
	// Setup key provider with special characters in password
	keyProvider := NewBaseKeyProvider()
	keyProvider.AddKeyValue("user", "p@ss:w0rd!#$%")

	config := BasicAuthConfig{
		KeyProvider: keyProvider,
	}

	basicAuth := NewBasicAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(basicAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Create request with special characters in password
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("user:p@ss:w0rd!#$%")))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestBasicAuth_Middleware_UpdateCredentials(t *testing.T) {
	// Setup key provider
	keyProvider := NewBaseKeyProvider()
	keyProvider.AddKeyValue("admin", "oldpassword")

	config := BasicAuthConfig{
		KeyProvider: keyProvider,
	}

	basicAuth := NewBasicAuth(config)

	// Create Fiber app
	app := fiber.New()
	app.Use(basicAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	// Test with old password
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:oldpassword")))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200 with old password, got %d", resp.StatusCode)
	}

	// Update password
	keyProvider.AddKeyValue("admin", "newpassword")

	// Test with old password (should fail)
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:oldpassword")))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 with old password after update, got %d", resp.StatusCode)
	}

	// Test with new password (should succeed)
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:newpassword")))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200 with new password, got %d", resp.StatusCode)
	}
}
