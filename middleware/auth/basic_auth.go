package auth

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
)

// BasicAuthConfig defines the configuration for BasicAuth middleware.
type BasicAuthConfig struct {
	KeyProvider  BaseKey
	Unauthorized fiber.Handler
}

// BasicAuth provides Basic Authentication middleware for Fiber.
type BasicAuth struct {
	config BasicAuthConfig
}

// NewBasicAuth creates a new instance of BasicAuth middleware with the provided configuration.
func NewBasicAuth(config BasicAuthConfig) *BasicAuth {
	return &BasicAuth{
		config: config,
	}
}

// Middleware returns the Fiber middleware handler for Basic Authentication.
func (b *BasicAuth) Middleware() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Authorizer: func(user, pass string, c fiber.Ctx) bool {
			storedPass, err := b.config.KeyProvider.GetValue(user)
			if err != nil {
				return false
			}
			return subtle.ConstantTimeCompare([]byte(pass), []byte(storedPass)) == 1
		},
		Unauthorized: b.config.Unauthorized,
	})
}
