package main

import (
	"log"

	"github.com/budimanlai/go-pkg/i18n"
	"github.com/budimanlai/go-pkg/response"
	"github.com/budimanlai/go-pkg/validator"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

type User struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Age      int    `json:"age" validate:"gte=18,lte=130"`
	Phone    string `json:"phone" validate:"required,numeric"`
}

func main() {
	// Setup I18n
	i18nConfig := i18n.I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id", "zh"},
		LocalesPath:     "locales",
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: response.FiberErrorHandler,
	})

	i18nManager, err := i18n.NewI18nManager(i18nConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Set i18n untuk response dan validator
	response.SetI18nManager(i18nManager)
	validator.SetI18nManager(i18nManager)

	// Use i18n middleware
	app.Use(i18n.I18nMiddleware(i18nConfig))

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, "API is running", fiber.Map{
			"version": "1.0.0",
		})
	})

	// Create user endpoint with validation
	app.Post("/users", func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err != nil {
			return response.BadRequest(c, "Invalid request body")
		}

		// Validate with automatic language detection and detailed field errors
		if err := validator.ValidateStructWithContext(c, user); err != nil {
			return response.ValidationErrorI18n(c, err)
		}

		// If validation passed, create user
		return response.Success(c, "User created successfully", user)
	})

	log.Println("Server running on http://localhost:3000")
	log.Println("\nTest dengan curl:")
	log.Println("  English:    curl -X POST 'http://localhost:3000/users?lang=en' -H 'Content-Type: application/json' -d '{\"name\":\"\",\"email\":\"invalid\",\"age\":15}'")
	log.Println("  Indonesian: curl -X POST 'http://localhost:3000/users?lang=id' -H 'Content-Type: application/json' -d '{\"name\":\"\",\"email\":\"invalid\",\"age\":15}'")
	log.Println("  Chinese:    curl -X POST 'http://localhost:3000/users?lang=zh' -H 'Content-Type: application/json' -d '{\"name\":\"\",\"email\":\"invalid\",\"age\":15}'")
	log.Println("\nValid request:")
	log.Println("  curl -X POST 'http://localhost:3000/users?lang=en' -H 'Content-Type: application/json' -d '{\"name\":\"John\",\"email\":\"john@example.com\",\"password\":\"password123\",\"age\":25,\"phone\":\"08123456789\"}'")

	app.Listen(":3000")
}
