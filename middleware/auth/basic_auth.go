package auth

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

// BasicAuthConfig defines the configuration for BasicAuth middleware.
type BasicAuthConfig struct {
	KeyProvider     BaseKey
	Unauthorized    fiber.Handler
	ContextUsername string
	ContextPassword string
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
		Users: nil,
		Authorizer: func(user, pass string) bool {
			// retrieve password from KeyProvider
			// use the provided username as the key
			storedPass, err := b.config.KeyProvider.GetValue(user)
			if err != nil {
				return false
			}
			if subtle.ConstantTimeCompare([]byte(pass), []byte(storedPass)) == 1 {
				return true
			}
			return false
		},
		Unauthorized:    b.config.Unauthorized,
		ContextUsername: b.config.ContextUsername,
		ContextPassword: b.config.ContextPassword,
	})
}
