# Validator Package

The `validator` package provides struct validation with user-friendly error messages and multi-language support through i18n integration.

## Overview

This package wraps the popular `go-playground/validator/v10` library and enhances it with:
- Automatic field name extraction from JSON tags
- Internationalized error messages
- Custom error type with field-level error details
- Seamless integration with Fiber web framework
- Support for multiple validation rules per field

## Features

- **Struct validation** with comprehensive validation tags
- **Internationalization (i18n)** support for error messages
- **Custom ValidationError** type with field-level error details
- **JSON tag integration** - uses JSON field names in error messages
- **Context-aware validation** - automatic language detection from Fiber context
- **Fallback messages** - English defaults when i18n is not configured
- **Field-specific errors** - map of field names to error messages

## Validation Error Format

The `ValidationError` type provides structured error information:

```go
type ValidationError struct {
    Messages []string            // All error messages
    Errors   map[string][]string // Field name -> error messages mapping
}
```

Methods:
- `Error()` - Returns all messages joined by semicolon
- `First()` - Returns the first error message
- `All()` - Returns all error messages as a slice
- `GetFieldErrors()` - Returns map of field names to error messages

## Installation

```bash
go get github.com/budimanlai/go-pkg/validator
```

## Quick Start

### Basic Usage

```go
import (
    "github.com/budimanlai/go-pkg/validator"
    "fmt"
)

type User struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Age      int    `json:"age" validate:"required,gte=18"`
}

func main() {
    user := User{
        Email:    "invalid-email",
        Password: "123",
        Age:      15,
    }
    
    if err := validator.ValidateStruct(user); err != nil {
        if verr, ok := err.(*validator.ValidationError); ok {
            fmt.Println(verr.First())
            // Output: Email must be a valid email address
            
            for field, errs := range verr.GetFieldErrors() {
                fmt.Printf("%s: %v\n", field, errs)
            }
        }
    }
}
```

### With Fiber and I18n

```go
import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/validator"
    "github.com/budimanlai/go-pkg/response"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Initialize i18n
    i18nMgr, _ := i18n.NewI18nManager(i18n.Config{
        LocalesPath:     "./locales",
        DefaultLanguage: "en",
    })
    
    // Set i18n manager for validator
    validator.SetI18nManager(i18nMgr)
    response.SetI18nManager(i18nMgr)
    
    app := fiber.New()
    app.Use(i18n.I18nMiddleware(i18nMgr))
    
    app.Post("/users", func(c *fiber.Ctx) error {
        var user User
        
        if err := c.BodyParser(&user); err != nil {
            return response.BadRequestI18n(c, "invalid_request_body", nil)
        }
        
        // Validate with automatic language detection
        if err := validator.ValidateStructWithContext(c, user); err != nil {
            return response.ValidationErrorI18n(c, err)
        }
        
        return response.SuccessI18n(c, "user_created", user)
    })
    
    app.Listen(":3000")
}
```

## Documentation Structure

- **[Validation Tags](validation-tags.md)** - Comprehensive list of validation rules
- **[Error Handling](error-handling.md)** - ValidationError type and error handling patterns
- **[I18n Integration](i18n-integration.md)** - Multilingual validation messages
- **[Examples](examples.md)** - Practical usage examples and patterns

## API Reference

### Validation Functions

| Function | Description |
|----------|-------------|
| `ValidateStruct(s)` | Validates struct with default language |
| `ValidateStructWithLang(s, lang)` | Validates struct with specified language |
| `ValidateStructWithContext(c, s)` | Validates struct with language from Fiber context |

### Setup Functions

| Function | Description |
|----------|-------------|
| `SetI18nManager(manager)` | Configure i18n manager for translations |

### ValidationError Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `Error()` | `string` | All messages joined by semicolon |
| `First()` | `string` | First error message |
| `All()` | `[]string` | All error messages |
| `GetFieldErrors()` | `map[string][]string` | Field names to error messages |

## Common Validation Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field is required | `validate:"required"` |
| `email` | Valid email address | `validate:"email"` |
| `min` | Minimum length/value | `validate:"min=8"` |
| `max` | Maximum length/value | `validate:"max=100"` |
| `gte` | Greater than or equal | `validate:"gte=18"` |
| `lte` | Less than or equal | `validate:"lte=100"` |
| `len` | Exact length | `validate:"len=10"` |
| `numeric` | Numeric characters only | `validate:"numeric"` |
| `alphanum` | Alphanumeric only | `validate:"alphanum"` |

See [Validation Tags](validation-tags.md) for complete list.

## Response Format

When using with `response.ValidationErrorI18n()`:

```json
{
  "meta": {
    "success": false,
    "message": "Email is required",
    "errors": {
      "email": ["Email is required"],
      "password": ["Password must be at least 8 characters"],
      "age": ["Age must be greater than or equal to 18"]
    }
  },
  "data": null
}
```

## Best Practices

1. **Use JSON tags** - Define JSON tags for consistent field names in errors
2. **Combine validation rules** - Use multiple tags separated by commas
3. **Set up i18n early** - Configure i18n manager during application initialization
4. **Use context validation** - Use `ValidateStructWithContext` in Fiber handlers
5. **Handle errors properly** - Always check for `ValidationError` type
6. **Provide translations** - Create locale files for all supported languages

## Requirements

- Go 1.18 or higher
- github.com/go-playground/validator/v10
- github.com/budimanlai/go-pkg/i18n (for i18n features)
- github.com/gofiber/fiber/v2 (for Fiber integration)

## Related Packages

- [I18n Package](../i18n.md) - Internationalization support
- [Response Package](../response/README.md) - HTTP response helpers with validation error formatting

## License

MIT License - see the LICENSE file for details
