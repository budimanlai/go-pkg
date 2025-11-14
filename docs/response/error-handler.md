# Fiber Error Handler

The `FiberErrorHandler` provides a custom error handler for Fiber applications that automatically processes errors and returns internationalized error responses.

## Overview

This error handler integrates with Fiber's error handling mechanism to automatically translate error messages based on the user's language preference. It handles different HTTP status codes and returns appropriate i18n error responses.

## Signature

```go
func FiberErrorHandler(ctx *fiber.Ctx, err error) error
```

### Parameters

- `ctx` (*fiber.Ctx) - The Fiber context containing request/response data
- `err` (error) - The error to be handled and formatted

### Returns

- `error` - An internationalized error response based on the status code

## How It Works

The error handler processes errors based on their HTTP status code:

1. **404 Not Found** → Returns `NotFoundI18n` response
2. **400 Bad Request** → Returns `BadRequestI18n` response
3. **Other Status Codes** → Returns `ErrorI18n` response with the corresponding code

If the error is a `*fiber.Error`, it uses the error's status code. Otherwise, it defaults to 500 (Internal Server Error).

## Setup

Configure the error handler when creating your Fiber application:

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
    response.SetI18nManager(i18nMgr)

    // Create Fiber app with custom error handler
    app := fiber.New(fiber.Config{
        ErrorHandler: response.FiberErrorHandler,
    })

    // Add i18n middleware
    app.Use(i18n.I18nMiddleware(i18nMgr))

    // Your routes here
    app.Listen(":3000")
}
```

## Usage Examples

### Automatic Error Handling

Once configured, the error handler automatically processes all errors returned by your handlers:

```go
app.Get("/users/:id", func(c *fiber.Ctx) error {
    user := getUserByID(c.Params("id"))
    if user == nil {
        // This will be caught by FiberErrorHandler
        return fiber.NewError(404, "user_not_found")
    }
    return c.JSON(user)
})
```

### With Fiber's Built-in Errors

```go
app.Post("/users", func(c *fiber.Ctx) error {
    var user User
    if err := c.BodyParser(&user); err != nil {
        // Returns 400 Bad Request with i18n message
        return fiber.NewError(400, "invalid_request_body")
    }
    
    // Process user...
    return c.JSON(user)
})
```

### Different Error Codes

```go
app.Delete("/users/:id", func(c *fiber.Ctx) error {
    // Check permission
    if !hasPermission(c) {
        // Returns 403 Forbidden with i18n message
        return fiber.NewError(403, "permission_denied")
    }
    
    // Check if user exists
    user := getUserByID(c.Params("id"))
    if user == nil {
        // Returns 404 Not Found with i18n message
        return fiber.NewError(404, "user_not_found")
    }
    
    // Delete user
    if err := deleteUser(user.ID); err != nil {
        // Returns 500 Internal Server Error with i18n message
        return fiber.NewError(500, "delete_failed")
    }
    
    return c.JSON(fiber.Map{"message": "User deleted"})
})
```

### Generic Error Handling

```go
app.Get("/process", func(c *fiber.Ctx) error {
    // Any error will be caught and processed
    if err := processData(); err != nil {
        return err // FiberErrorHandler will handle this
    }
    return c.JSON(fiber.Map{"status": "success"})
})
```

## Response Examples

### 404 Not Found

**Request:**
```http
GET /users/999
Accept-Language: en
```

**Response:**
```json
{
  "meta": {
    "success": false,
    "message": "User not found"
  },
  "data": null
}
```

**With Indonesian:**
```http
GET /users/999
Accept-Language: id
```

**Response:**
```json
{
  "meta": {
    "success": false,
    "message": "Pengguna tidak ditemukan"
  },
  "data": null
}
```

### 400 Bad Request

**Request:**
```http
POST /users
Accept-Language: en
Content-Type: application/json

{ "invalid": "data" }
```

**Response:**
```json
{
  "meta": {
    "success": false,
    "message": "Invalid request body"
  },
  "data": null
}
```

### 500 Internal Server Error

**Request:**
```http
GET /process
Accept-Language: en
```

**Response:**
```json
{
  "meta": {
    "success": false,
    "message": "Processing failed"
  },
  "data": null
}
```

## Locale File Setup

Create translation files for error messages:

**locales/en.json:**
```json
{
  "user_not_found": "User not found",
  "invalid_request_body": "Invalid request body",
  "permission_denied": "Permission denied",
  "delete_failed": "Failed to delete user",
  "processing_failed": "Processing failed"
}
```

**locales/id.json:**
```json
{
  "user_not_found": "Pengguna tidak ditemukan",
  "invalid_request_body": "Body request tidak valid",
  "permission_denied": "Izin ditolak",
  "delete_failed": "Gagal menghapus pengguna",
  "processing_failed": "Pemrosesan gagal"
}
```

## Advanced Usage

### Custom Error Types

```go
type AppError struct {
    Code    int
    Message string
    Data    interface{}
}

func (e *AppError) Error() string {
    return e.Message
}

// In your handler
app.Get("/custom", func(c *fiber.Ctx) error {
    return &AppError{
        Code:    400,
        Message: "custom_error",
        Data: map[string]string{
            "Field": "value",
        },
    }
})
```

### Error Logging

```go
func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
    // Log the error
    log.Printf("Error: %v, Path: %s", err, ctx.Path())
    
    // Use the default error handler
    return response.FiberErrorHandler(ctx, err)
}

app := fiber.New(fiber.Config{
    ErrorHandler: CustomErrorHandler,
})
```

### Fallback Behavior

If i18nManager is not set, the error handler falls back to standard response functions:

```go
// Without i18n
app := fiber.New(fiber.Config{
    ErrorHandler: response.FiberErrorHandler,
})

// Will use error message as-is, without translation
app.Get("/test", func(c *fiber.Ctx) error {
    return fiber.NewError(404, "Resource not found")
})
```

## Error Code Mapping

| HTTP Status | Function Called | Description |
|-------------|----------------|-------------|
| 404 | `NotFoundI18n` | Resource not found |
| 400 | `BadRequestI18n` | Bad request/invalid input |
| 403 | `ErrorI18n(403)` | Forbidden/permission denied |
| 401 | `ErrorI18n(401)` | Unauthorized |
| 500 | `ErrorI18n(500)` | Internal server error |
| Other | `ErrorI18n(code)` | Custom status codes |

## Best Practices

1. **Always Use Message IDs** - Use translation keys instead of hardcoded messages
2. **Consistent Error Codes** - Use standard HTTP status codes appropriately
3. **Error Logging** - Implement logging wrapper for production environments
4. **Fallback Messages** - Provide default English messages in locale files
5. **Testing** - Test error responses with different languages

## Benefits

- **Automatic Translation** - No need to manually call i18n functions in error handling
- **Consistent Format** - All errors follow the same JSON structure
- **Clean Code** - Reduces boilerplate in route handlers
- **Centralized** - Single point for error response formatting
- **Multilingual** - Automatic language detection and translation

## Related Documentation

- [I18n Responses](i18n-responses.md) - Internationalized response functions
- [Standard Responses](standard-responses.md) - Basic response functions
- [I18n Package](../i18n.md) - Internationalization configuration
