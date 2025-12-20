# Header Authentication Middleware

Middleware for authentication using `X-API-Key` header or custom header name.

## Features

- ✅ API key validation via HTTP headers
- ✅ Custom header name support
- ✅ Success handler for custom logic after successful validation
- ✅ Custom error handler
- ✅ Case-insensitive header names
- ✅ Support for all HTTP methods (GET, POST, PUT, DELETE, etc.)

## Installation

```go
import "github.com/budimanlai/go-pkg/middleware/auth"
```

## Usage

### Basic Usage with Default Header Name (X-API-Key)

```go
package main

import (
    "github.com/budimanlai/go-pkg/middleware/auth"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    // Setup key provider
    keyProvider := auth.NewBaseKeyProvider()
    keyProvider.Add("my-secret-api-key-123")
    keyProvider.Add("another-valid-key-456")

    // Create header auth middleware
    headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
        KeyProvider: keyProvider,
        // HeaderName default: "X-API-Key"
    })

    // Apply middleware
    app.Use(headerAuth.Middleware())

    app.Get("/api/data", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Access granted!",
            "data":    []string{"item1", "item2", "item3"},
        })
    })

    app.Listen(":3000")
}
```

**Request:**
```bash
curl -H "X-API-Key: my-secret-api-key-123" http://localhost:3000/api/data
```

### Custom Header Name

```go
headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
    KeyProvider: keyProvider,
    HeaderName:  "X-Custom-API-Key",
})
```

**Request:**
```bash
curl -H "X-Custom-API-Key: my-secret-api-key-123" http://localhost:3000/api/data
```

### With Success Handler

Success handler is called after API key is successfully validated. Useful for:
- Storing user information in context
- Logging
- Increment usage counter
- etc

```go
successHandler := func(c *fiber.Ctx, token string) error {
    // Get user info from database based on token
    user := getUserByAPIKey(token)
    
    // Store in context for use in next handlers
    c.Locals("user_id", user.ID)
    c.Locals("user_email", user.Email)
    c.Locals("api_key", token)
    
    // Log access
    log.Printf("API accessed by user: %s with key: %s", user.Email, token)
    
    return nil
}

headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
    KeyProvider:    keyProvider,
    HeaderName:     "X-API-Key",
    SuccessHandler: &successHandler,
})
```

### With Custom Error Handler

```go
customErrorHandler := func(c *fiber.Ctx, err error) error {
    return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
        "success": false,
        "error":   "Invalid or missing API key",
        "message": "Please provide a valid API key in the X-API-Key header",
    })
}

headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
    KeyProvider:  keyProvider,
    HeaderName:   "X-API-Key",
    ErrorHandler: customErrorHandler,
})
```

### With Database Key Provider

```go
import (
    "github.com/budimanlai/go-pkg/middleware/auth"
    "gorm.io/gorm"
)

func setupHeaderAuth(db *gorm.DB) *auth.HeaderAuth {
    // Gunakan database sebagai key provider
    keyProvider := auth.NewDbApiKey(db)
    
    headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
        KeyProvider: keyProvider,
        HeaderName:  "X-API-Key",
    })
    
    return headerAuth
}
```

### Apply ke Specific Routes

```go
app := fiber.New()

// Setup auth
keyProvider := auth.NewBaseKeyProvider()
keyProvider.Add("secret-key-123")

headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
    KeyProvider: keyProvider,
})

// Public routes (tidak butuh auth)
app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Welcome to public page")
})

app.Get("/about", func(c *fiber.Ctx) error {
    return c.SendString("About us page")
})

// Protected routes group
api := app.Group("/api")
api.Use(headerAuth.Middleware()) // Apply auth hanya ke /api/*

api.Get("/users", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"users": []string{"user1", "user2"}})
})

api.Get("/products", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"products": []string{"product1", "product2"}})
})
```

## Configuration Options

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `KeyProvider` | `BaseKey` | Interface for API key validation (required) | - |
| `HeaderName` | `string` | Header name for API key | `"X-API-Key"` |
| `SuccessHandler` | `*func(c *fiber.Ctx, token string) error` | Function called after successful validation | `nil` |
| `ErrorHandler` | `fiber.ErrorHandler` | Custom error handler for invalid/missing keys | `nil` |

### `NewHeaderAuth(config HeaderAuthConfig) *HeaderAuth`
Creates a new instance of HeaderAuth middleware.

### `GetHeaderName() string`
Gets the header name being used.

### `SetHeaderName(headerName string)`
Changes the header name being used.

### `Middleware() fiber.Handler`
Returns Fiber middleware handler.

## Response Codes

- `200 OK` - API key valid
- `401 Unauthorized` - Invalid or missing API key (default)
- Custom status code if using custom error handler

## Examples

### Complete Example with Database

```go
package main

import (
    "log"
    
    "github.com/budimanlai/go-pkg/middleware/auth"
    "github.com/gofiber/fiber/v2"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    // Setup database
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate
    db.AutoMigrate(&auth.ApiKey{})

    // Setup key provider with database
    keyProvider := auth.NewDbApiKey(db)
    
    // Add some keys
    keyProvider.AddKeyValue("client-1-key", "client-1-secret")
    keyProvider.AddKeyValue("client-2-key", "client-2-secret")

    // Success handler
    successHandler := func(c *fiber.Ctx, token string) error {
        // Get auth value from database
        authValue, _ := keyProvider.GetValue(token)
        c.Locals("client_secret", authValue)
        c.Locals("api_key", token)
        
        log.Printf("API accessed with key: %s", token)
        return nil
    }

    // Custom error handler
    errorHandler := func(c *fiber.Ctx, err error) error {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "success": false,
            "error":   "Unauthorized",
            "message": "Invalid or missing API key",
        })
    }

    // Create auth middleware
    headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
        KeyProvider:    keyProvider,
        HeaderName:     "X-API-Key",
        SuccessHandler: &successHandler,
        ErrorHandler:   errorHandler,
    })

    // Setup Fiber app
    app := fiber.New()

    // Public routes
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Welcome to the API",
            "version": "1.0.0",
        })
    })

    // Protected API routes
    api := app.Group("/api")
    api.Use(headerAuth.Middleware())

    api.Get("/profile", func(c *fiber.Ctx) error {
        apiKey := c.Locals("api_key").(string)
        return c.JSON(fiber.Map{
            "message": "Profile data",
            "api_key": apiKey,
        })
    })

    api.Get("/data", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "data": []string{"item1", "item2", "item3"},
        })
    })

    log.Fatal(app.Listen(":3000"))
}
```

## Testing

```bash
# Success request
curl -H "X-API-Key: client-1-key" http://localhost:3000/api/data

# Failed request (missing header)
curl http://localhost:3000/api/data

# Failed request (invalid key)
curl -H "X-API-Key: invalid-key" http://localhost:3000/api/data
```

## Key Provider Options

1. **BaseKeyProvider** (In-memory)
   ```go
   keyProvider := auth.NewBaseKeyProvider()
   keyProvider.Add("key-1")
   keyProvider.AddKeyValue("key-2", "value-2")
   ```

2. **DbApiKey** (Database)
   ```go
   keyProvider := auth.NewDbApiKey(db)
   keyProvider.AddKeyValue("api-key", "secret-value")
   ```

## Best Practices

1. **Use HTTPS** - Always use HTTPS in production to protect API keys
2. **Rotate Keys** - Periodically update API keys
3. **Different Keys per Client** - Give different API key for each client
4. **Log Access** - Use success handler for logging
5. **Rate Limiting** - Combine with rate limiting middleware
6. **Monitor Usage** - Track API key usage to detect abuse

## See Also

- [Basic Auth](./basic_auth.md)
- [QueryString Auth](./querystring_auth.md)
- [BaseKey Interface](./basekey_interface.md)
