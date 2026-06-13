package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
)

type QueryStringAuthConfig struct {
	// KeyProvider is the source of valid API keys.
	KeyProvider BaseKey

	// The name of the query string parameter to look for the API key.
	// Default is "access-token".
	ParamName string

	// function called if the key is valid
	SuccessHandler *func(c fiber.Ctx, token string) error

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
		Extractor: extractors.FromQuery(qsa.config.ParamName),

		Validator: func(c fiber.Ctx, key string) (bool, error) {
			if qsa.config.KeyProvider.IsExists(key) {
				if qsa.config.SuccessHandler != nil {
					if err := (*qsa.config.SuccessHandler)(c, key); err != nil {
						return false, err
					}
				}
				return true, nil
			}
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},

		ErrorHandler: qsa.config.ErrorHandler,
	})
}
