package auth

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Helper function to generate JWT token for testing
func generateTestToken(secretKey string, claims jwt.MapClaims, signingMethod string) string {
	var method jwt.SigningMethod
	switch signingMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	case "HS512":
		method = jwt.SigningMethodHS512
	default:
		method = jwt.SigningMethodHS256
	}

	token := jwt.NewWithClaims(method, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}

func TestNewJWTAuth(t *testing.T) {
	config := JWTConfig{
		SecretKey: "test-secret-key",
	}

	jwtAuth := NewJWTAuth(config)
	if jwtAuth == nil {
		t.Fatal("Expected NewJWTAuth to return a non-nil instance")
	}

	if jwtAuth.config.SigningMethod != "HS256" {
		t.Errorf("Expected default SigningMethod 'HS256', got '%s'", jwtAuth.config.SigningMethod)
	}

	if jwtAuth.config.TokenLookup != "header:Authorization" {
		t.Errorf("Expected default TokenLookup 'header:Authorization', got '%s'", jwtAuth.config.TokenLookup)
	}

	if jwtAuth.config.AuthScheme != "Bearer" {
		t.Errorf("Expected default AuthScheme 'Bearer', got '%s'", jwtAuth.config.AuthScheme)
	}

	if jwtAuth.config.ContextKey != "user" {
		t.Errorf("Expected default ContextKey 'user', got '%s'", jwtAuth.config.ContextKey)
	}
}

func TestJWTAuth_Middleware_Success(t *testing.T) {
	secretKey := "test-secret-key"
	config := JWTConfig{
		SecretKey: secretKey,
	}

	jwtAuth := NewJWTAuth(config)

	claims := jwt.MapClaims{
		"user_id": "123",
		"email":   "test@example.com",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := generateTestToken(secretKey, claims, "HS256")

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

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

func TestJWTAuth_Middleware_MissingToken(t *testing.T) {
	config := JWTConfig{
		SecretKey: "test-secret-key",
	}

	jwtAuth := NewJWTAuth(config)

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
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

func TestJWTAuth_Middleware_InvalidToken(t *testing.T) {
	config := JWTConfig{
		SecretKey: "test-secret-key",
	}

	jwtAuth := NewJWTAuth(config)

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestJWTAuth_Middleware_ExpiredToken(t *testing.T) {
	secretKey := "test-secret-key"
	config := JWTConfig{
		SecretKey: secretKey,
	}

	jwtAuth := NewJWTAuth(config)

	claims := jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}
	token := generateTestToken(secretKey, claims, "HS256")

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 for expired token, got %d", resp.StatusCode)
	}
}

func TestJWTAuth_Middleware_WrongSecret(t *testing.T) {
	config := JWTConfig{
		SecretKey: "correct-secret-key",
	}

	jwtAuth := NewJWTAuth(config)

	claims := jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := generateTestToken("wrong-secret-key", claims, "HS256")

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 for wrong secret, got %d", resp.StatusCode)
	}
}

func TestJWTAuth_Middleware_TokenFromQuery(t *testing.T) {
	secretKey := "test-secret-key"
	config := JWTConfig{
		SecretKey:   secretKey,
		TokenLookup: "query:token",
		AuthScheme:  "",
	}

	jwtAuth := NewJWTAuth(config)

	claims := jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := generateTestToken(secretKey, claims, "HS256")

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test?token="+token, nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestJWTAuth_Middleware_TokenFromCookie(t *testing.T) {
	secretKey := "test-secret-key"
	config := JWTConfig{
		SecretKey:   secretKey,
		TokenLookup: "cookie:jwt",
		AuthScheme:  "",
	}

	jwtAuth := NewJWTAuth(config)

	claims := jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := generateTestToken(secretKey, claims, "HS256")

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", "jwt="+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestJWTAuth_Middleware_SuccessHandler(t *testing.T) {
	secretKey := "test-secret-key"

	successHandlerCalled := false
	var capturedUserID string

	successHandler := func(c *fiber.Ctx, claims jwt.MapClaims) error {
		successHandlerCalled = true
		capturedUserID = claims["user_id"].(string)
		c.Locals("custom_data", "test-data")
		return nil
	}

	config := JWTConfig{
		SecretKey:      secretKey,
		SuccessHandler: successHandler,
	}

	jwtAuth := NewJWTAuth(config)

	claims := jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := generateTestToken(secretKey, claims, "HS256")

	var customData string

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		customData = c.Locals("custom_data").(string)
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

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

	if capturedUserID != "123" {
		t.Errorf("Expected user_id '123', got '%s'", capturedUserID)
	}

	if customData != "test-data" {
		t.Errorf("Expected custom_data 'test-data', got '%s'", customData)
	}
}

func TestJWTAuth_Middleware_CustomErrorHandler(t *testing.T) {
	errorHandlerCalled := false

	customErrorHandler := func(c *fiber.Ctx, err error) error {
		errorHandlerCalled = true
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Custom JWT Error",
		})
	}

	config := JWTConfig{
		SecretKey:    "test-secret-key",
		ErrorHandler: customErrorHandler,
	}

	jwtAuth := NewJWTAuth(config)

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

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

func TestJWTAuth_Middleware_ContextStorage(t *testing.T) {
	secretKey := "test-secret-key"
	config := JWTConfig{
		SecretKey:  secretKey,
		ContextKey: "jwt_claims",
	}

	jwtAuth := NewJWTAuth(config)

	claims := jwt.MapClaims{
		"user_id": "123",
		"email":   "test@example.com",
		"role":    "admin",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := generateTestToken(secretKey, claims, "HS256")

	var storedClaims jwt.MapClaims

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		storedClaims = c.Locals("jwt_claims").(jwt.MapClaims)
		return c.SendString("Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if storedClaims["user_id"] != "123" {
		t.Errorf("Expected user_id '123', got '%v'", storedClaims["user_id"])
	}

	if storedClaims["email"] != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%v'", storedClaims["email"])
	}

	if storedClaims["role"] != "admin" {
		t.Errorf("Expected role 'admin', got '%v'", storedClaims["role"])
	}
}

func TestJWTAuth_Middleware_DifferentHTTPMethods(t *testing.T) {
	secretKey := "test-secret-key"
	config := JWTConfig{
		SecretKey: secretKey,
	}

	jwtAuth := NewJWTAuth(config)

	claims := jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := generateTestToken(secretKey, claims, "HS256")

	app := fiber.New()
	app.Use(jwtAuth.Middleware())
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
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make %s request: %v", method, err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200 for %s method, got %d", method, resp.StatusCode)
		}
	}
}

func TestSetSecretKey(t *testing.T) {
	// Create JWT auth with initial secret key
	jwtAuth := NewJWTAuth(JWTConfig{
		SecretKey: "initial-secret",
	})

	// Verify initial secret key
	if jwtAuth.GetSecretKey() != "initial-secret" {
		t.Errorf("Expected initial secret key 'initial-secret', got '%s'", jwtAuth.GetSecretKey())
	}

	// Change secret key
	jwtAuth.SetSecretKey("new-secret-key")

	// Verify new secret key
	if jwtAuth.GetSecretKey() != "new-secret-key" {
		t.Errorf("Expected new secret key 'new-secret-key', got '%s'", jwtAuth.GetSecretKey())
	}
}

func TestSetSecretKeyDynamic(t *testing.T) {
	app := fiber.New()

	jwtAuth := NewJWTAuth(JWTConfig{
		SecretKey: "secret-v1",
	})

	app.Use(jwtAuth.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Generate token with secret-v1
	claims := jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token1 := generateTestToken("secret-v1", claims, "HS256")

	// Test with token1 (should succeed)
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.Header.Set("Authorization", "Bearer "+token1)
	resp1, _ := app.Test(req1)

	if resp1.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200 with secret-v1, got %d", resp1.StatusCode)
	}

	// Change secret key dynamically
	jwtAuth.SetSecretKey("secret-v2")

	// Test with old token (should fail)
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Authorization", "Bearer "+token1)
	resp2, _ := app.Test(req2)

	if resp2.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 with old token after secret change, got %d", resp2.StatusCode)
	}

	// Generate new token with secret-v2
	token2 := generateTestToken("secret-v2", claims, "HS256")

	// Test with new token (should succeed)
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.Header.Set("Authorization", "Bearer "+token2)
	resp3, _ := app.Test(req3)

	if resp3.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200 with secret-v2, got %d", resp3.StatusCode)
	}
}

func TestSetSecretKeyConcurrent(t *testing.T) {
	jwtAuth := NewJWTAuth(JWTConfig{
		SecretKey: "initial-secret",
	})

	done := make(chan bool)

	// Goroutine 1: Read secret key repeatedly
	go func() {
		for i := 0; i < 100; i++ {
			_ = jwtAuth.GetSecretKey()
		}
		done <- true
	}()

	// Goroutine 2: Write secret key repeatedly
	go func() {
		for i := 0; i < 100; i++ {
			jwtAuth.SetSecretKey("secret-" + string(rune(i)))
		}
		done <- true
	}()

	// Goroutine 3: Read secret key repeatedly
	go func() {
		for i := 0; i < 100; i++ {
			_ = jwtAuth.GetSecretKey()
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	<-done
	<-done
	<-done

	// Test should not panic or race
	t.Log("Concurrent access test passed")
}
