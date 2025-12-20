package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

// HeaderAuthConfig defines the configuration for Header-based API Key Authentication.
type HeaderAuthConfig struct {
	// KeyProvider is the source of valid API keys.
	KeyProvider BaseKey

	// The name of the header to look for the API key.
	// Default is "X-API-Key".
	HeaderName string

	// function called if the key is valid
	SuccessHandler *func(c *fiber.Ctx, token string) error

	// function called if the key is invalid or missing
	ErrorHandler fiber.ErrorHandler
}

// HeaderAuth provides Header-based API Key Authentication middleware for Fiber.
type HeaderAuth struct {
	config HeaderAuthConfig
}

// NewHeaderAuth creates a new instance of HeaderAuth with the provided configuration.
func NewHeaderAuth(config HeaderAuthConfig) *HeaderAuth {
	// Set default header name if not provided
	if config.HeaderName == "" {
		config.HeaderName = "X-API-Key"
	}

	return &HeaderAuth{
		config: config,
	}
}

// GetHeaderName returns the configured header name.
func (ha *HeaderAuth) GetHeaderName() string {
	return ha.config.HeaderName
}

// SetHeaderName sets the header name for API key lookup.
func (ha *HeaderAuth) SetHeaderName(name string) {
	ha.config.HeaderName = name
}

// Middleware returns the Fiber middleware handler for Header-based API Key Authentication.
func (ha *HeaderAuth) Middleware() fiber.Handler {
	return keyauth.New(keyauth.Config{
		// Define where to look for the key: "header:X-API-Key"
		KeyLookup: "header:" + ha.config.HeaderName,

		// Define the function to validate the extracted key
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			if ha.config.KeyProvider.IsExists(key) {
				if ha.config.SuccessHandler != nil {
					// Call the custom success handler
					if err := (*ha.config.SuccessHandler)(c, key); err != nil {
						return false, err
					}
				}
				return true, nil
			}

			// Key is invalid
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},

		// Optional: Error handler for invalid/missing keys
		ErrorHandler: ha.config.ErrorHandler,
	})
}
