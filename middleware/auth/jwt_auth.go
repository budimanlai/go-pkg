package auth

import (
	"errors"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig defines the configuration for JWT middleware.
type JWTConfig struct {
	// SecretKey is used to sign and validate JWT tokens
	SecretKey string

	// SigningMethod defines the signing method (default: HS256)
	SigningMethod string

	// TokenLookup defines where to look for the JWT token
	// Format: "<source>:<name>"
	// Possible values:
	// - "header:Authorization"
	// - "query:token"
	// - "cookie:jwt"
	// Default: "header:Authorization"
	TokenLookup string

	// AuthScheme defines the authorization scheme (default: "Bearer")
	AuthScheme string

	// ContextKey is the key used to store user claims in context
	// Default: "user"
	ContextKey string

	// SuccessHandler is called after successful JWT validation
	SuccessHandler func(c *fiber.Ctx, claims jwt.MapClaims) error

	// ErrorHandler is called when JWT validation fails
	ErrorHandler fiber.ErrorHandler

	// Claims is a custom claims struct that implements jwt.Claims interface
	// If not provided, jwt.MapClaims will be used
	Claims jwt.Claims
}

// JWTAuth provides JWT Authentication middleware for Fiber.
type JWTAuth struct {
	config JWTConfig
	mu     sync.RWMutex
}

var (
	// ErrJWTMissing indicates that JWT token is missing
	ErrJWTMissing = errors.New("missing or malformed JWT")

	// ErrJWTInvalid indicates that JWT token is invalid
	ErrJWTInvalid = errors.New("invalid or expired JWT")
)

// NewJWTAuth creates a new instance of JWTAuth middleware.
func NewJWTAuth(config JWTConfig) *JWTAuth {
	// Set defaults
	if config.SigningMethod == "" {
		config.SigningMethod = "HS256"
	}
	if config.TokenLookup == "" {
		config.TokenLookup = "header:Authorization"
	}
	if config.AuthScheme == "" {
		config.AuthScheme = "Bearer"
	}
	if config.ContextKey == "" {
		config.ContextKey = "user"
	}

	return &JWTAuth{
		config: config,
	}
}

// Middleware returns the Fiber middleware handler for JWT Authentication.
func (j *JWTAuth) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from request
		tokenString, err := j.extractToken(c)
		if err != nil {
			if j.config.ErrorHandler != nil {
				return j.config.ErrorHandler(c, err)
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
		}

		// Parse and validate token with read lock
		j.mu.RLock()
		token, err := j.parseToken(tokenString)
		j.mu.RUnlock()

		if err != nil || !token.Valid {
			if j.config.ErrorHandler != nil {
				return j.config.ErrorHandler(c, ErrJWTInvalid)
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid or expired JWT",
			})
		}

		// Extract claims
		var claims jwt.MapClaims
		if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
			claims = mapClaims
		} else {
			claims = jwt.MapClaims{}
		}

		// Store claims in context
		c.Locals(j.config.ContextKey, claims)

		// Call success handler if provided
		if j.config.SuccessHandler != nil {
			if err := j.config.SuccessHandler(c, claims); err != nil {
				if j.config.ErrorHandler != nil {
					return j.config.ErrorHandler(c, err)
				}
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   "Unauthorized",
					"message": err.Error(),
				})
			}
		}

		return c.Next()
	}
}

// extractToken extracts JWT token from request based on TokenLookup configuration
func (j *JWTAuth) extractToken(c *fiber.Ctx) (string, error) {
	parts := strings.Split(j.config.TokenLookup, ":")
	if len(parts) != 2 {
		return "", ErrJWTMissing
	}

	source := parts[0]
	name := parts[1]

	var tokenString string

	switch source {
	case "header":
		authHeader := c.Get(name)
		if authHeader == "" {
			return "", ErrJWTMissing
		}

		// Check for Bearer scheme
		if j.config.AuthScheme != "" {
			prefix := j.config.AuthScheme + " "
			if !strings.HasPrefix(authHeader, prefix) {
				return "", ErrJWTMissing
			}
			tokenString = strings.TrimPrefix(authHeader, prefix)
		} else {
			tokenString = authHeader
		}

	case "query":
		tokenString = c.Query(name)
		if tokenString == "" {
			return "", ErrJWTMissing
		}

	case "cookie":
		tokenString = c.Cookies(name)
		if tokenString == "" {
			return "", ErrJWTMissing
		}

	default:
		return "", ErrJWTMissing
	}

	return tokenString, nil
}

// parseToken parses and validates the JWT token
func (j *JWTAuth) parseToken(tokenString string) (*jwt.Token, error) {
	// Determine claims type
	var claims jwt.Claims
	if j.config.Claims != nil {
		claims = j.config.Claims
	} else {
		claims = jwt.MapClaims{}
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if token.Method.Alg() != j.config.SigningMethod {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.config.SecretKey), nil
	})

	return token, err
}

// GetSecretKey returns the secret key used for JWT signing
func (j *JWTAuth) GetSecretKey() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.config.SecretKey
}

// SetSecretKey sets a new secret key for JWT signing dynamically
func (j *JWTAuth) SetSecretKey(secretKey string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.config.SecretKey = secretKey
}

// GetSigningMethod returns the signing method
func (j *JWTAuth) GetSigningMethod() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.config.SigningMethod
}

// GetContextKey returns the context key for storing claims
func (j *JWTAuth) GetContextKey() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.config.ContextKey
}
