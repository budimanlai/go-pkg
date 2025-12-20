package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

type QueryStringAuthConfig struct {
	// KeyProvider is the source of valid API keys.
	KeyProvider BaseKey

	// The name of the query string parameter to look for the API key.
	// Default is "access-token".
	ParamName string

	// function called if the key is valid
	SuccessHandler *func(c *fiber.Ctx, token string) error

	// function called if the key is invalid or missing
	ErrorHandler fiber.ErrorHandler
}

type QueryStringAuth struct {
	config QueryStringAuthConfig
}

// NewDefaultQueryStringAuth returns a QueryStringAuth with default values.
func NewDefaultQueryStringAuth(config QueryStringAuthConfig) *QueryStringAuth {
	return &QueryStringAuth{
		config: config,
	}
}

func (qsa *QueryStringAuth) GetParamName() string {
	return qsa.config.ParamName
}

func (qsa *QueryStringAuth) SetParamName(name string) {
	qsa.config.ParamName = name
}

// Middleware returns the Fiber middleware handler for Query String Authentication.
func (qsa *QueryStringAuth) Middleware() fiber.Handler {
	return keyauth.New(keyauth.Config{
		// Define where to look for the key: "query:access-token" looks for ?access-token=...
		KeyLookup: "query:" + qsa.config.ParamName,

		// Define the function to validate the extracted key
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			if qsa.config.KeyProvider.IsExists(key) {
				if qsa.config.SuccessHandler != nil {
					// Call the custom valid function
					if err := (*qsa.config.SuccessHandler)(c, key); err != nil {
						return false, err
					}
				}
				return true, nil
			}

			// Key is invalid
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},

		// Optional: Error handler for invalid/missing keys
		ErrorHandler: qsa.config.ErrorHandler,
	})
}
