# Response Package

The `response` package provides standardized HTTP response helpers for Fiber applications with built-in internationalization (i18n) support.

## Overview

This package offers a consistent way to return JSON responses from your Fiber web applications. It includes both standard and internationalized (i18n) response functions, making it easy to build multilingual APIs.

## Features

- **Standardized JSON response format** with `meta` and `data` structure
- **Internationalization (i18n) support** for multilingual error messages
- **Pre-built response helpers** for common HTTP status codes
- **Validation error formatting** with field-level error details
- **Custom Fiber error handler** with automatic i18n integration
- **Type-safe responses** with consistent structure

## Response Format

All responses follow a standardized JSON structure:

```json
{
  "meta": {
    "success": true,
    "message": "Success message",
    "errors": null
  },
  "data": {
    // your data here
  }
}
```

For error responses:

```json
{
  "meta": {
    "success": false,
    "message": "Error message",
    "errors": {
      "Email": ["Email is required", "Email must be valid"],
      "Password": ["Password is too short"]
    }
  },
  "data": null
}
```

## Installation

```bash
go get github.com/budimanlai/go-pkg/response
```

## Quick Start

### Basic Usage

```go
import (
    "github.com/budimanlai/go-pkg/response"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    // Success response
    app.Get("/users/:id", func(c *fiber.Ctx) error {
        user := getUserByID(c.Params("id"))
        return response.Success(c, "User retrieved successfully", user)
    })

    // Error response
    app.Get("/users/:id", func(c *fiber.Ctx) error {
        return response.NotFound(c, "User not found")
    })

    app.Listen(":3000")
}
```

### With Internationalization

```go
import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/response"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Initialize i18n
    i18nMgr, _ := i18n.NewI18nManager(i18n.Config{
        LocalesPath:     "./locales",
        DefaultLanguage: "en",
    })

    // Set i18n manager for response package
    response.SetI18nManager(i18nMgr)

    app := fiber.New()

    // Use i18n middleware
    app.Use(i18n.I18nMiddleware(i18nMgr))

    // I18n responses automatically use the language from request
    app.Post("/users", func(c *fiber.Ctx) error {
        var user User
        if err := c.BodyParser(&user); err != nil {
            return response.BadRequestI18n(c, "invalid_request", nil)
        }

        // Language will be determined from Accept-Language header
        return response.SuccessI18n(c, "user_created", user)
    })

    app.Listen(":3000")
}
```

## Documentation Structure

- **[Standard Responses](standard-responses.md)** - Basic response functions without i18n
- **[I18n Responses](i18n-responses.md)** - Internationalized response functions
- **[Error Handler](error-handler.md)** - Custom Fiber error handler
- **[Examples](examples.md)** - Practical usage examples and patterns

## API Reference

### Standard Response Functions

| Function | HTTP Status | Description |
|----------|-------------|-------------|
| `Success(c, message, data)` | 200 OK | Success response with data |
| `Error(c, code, message)` | Custom | Generic error response |
| `BadRequest(c, message)` | 400 | Bad request error |
| `NotFound(c, message)` | 404 | Resource not found |

### I18n Response Functions

| Function | HTTP Status | Description |
|----------|-------------|-------------|
| `SuccessI18n(c, messageID, data)` | 200 OK | Translated success response |
| `ErrorI18n(c, code, messageID, data)` | Custom | Translated error response |
| `BadRequestI18n(c, messageID, data)` | 400 | Translated bad request |
| `NotFoundI18n(c, messageID)` | 404 | Translated not found |
| `ValidationErrorI18n(c, err)` | 400 | Validation errors with field details |

### Setup Functions

| Function | Description |
|----------|-------------|
| `SetI18nManager(manager)` | Configure i18n manager for translations |
| `FiberErrorHandler(ctx, err)` | Custom error handler for Fiber app |

## Best Practices

1. **Use I18n for Production** - Always use i18n responses in production applications
2. **Consistent Message IDs** - Use consistent message IDs across your application
3. **Error Handler** - Set up the custom error handler for automatic i18n error handling
4. **Validation Errors** - Use `ValidationErrorI18n` for detailed field-level error reporting
5. **Status Codes** - Use appropriate HTTP status codes for different scenarios

## Requirements

- Go 1.18 or higher
- Fiber v2
- github.com/budimanlai/go-pkg/i18n (for i18n features)

## Related Packages

- [I18n Package](../i18n.md) - Internationalization support
- [Validator Package](../validator/README.md) - Struct validation with i18n

## License

MIT License - see the LICENSE file for details
