# Header Authentication Middleware

Middleware untuk autentikasi menggunakan header `X-API-Key` atau custom header name.

## Features

- ✅ Validasi API key melalui HTTP header
- ✅ Custom header name support
- ✅ Success handler untuk custom logic setelah validasi berhasil
- ✅ Custom error handler
- ✅ Case-insensitive header names
- ✅ Support untuk semua HTTP methods (GET, POST, PUT, DELETE, dll)

## Installation

```go
import "github.com/budimanlai/go-pkg/middleware/auth"
```

## Usage

### Basic Usage dengan Default Header Name (X-API-Key)

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

### Dengan Success Handler

Success handler dipanggil setelah API key berhasil divalidasi. Berguna untuk:
- Menyimpan informasi user ke context
- Logging
- Increment usage counter
- dll

```go
successHandler := func(c *fiber.Ctx, token string) error {
    // Ambil user info dari database berdasarkan token
    user := getUserByAPIKey(token)
    
    // Simpan ke context untuk digunakan di handler selanjutnya
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

### Dengan Custom Error Handler

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

### Dengan Database Key Provider

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
| `KeyProvider` | `BaseKey` | Interface untuk validasi API keys (required) | - |
| `HeaderName` | `string` | Nama header untuk API key | `"X-API-Key"` |
| `SuccessHandler` | `*func(c *fiber.Ctx, token string) error` | Function yang dipanggil setelah validasi berhasil | `nil` |
| `ErrorHandler` | `fiber.ErrorHandler` | Custom error handler untuk invalid/missing keys | `nil` |

## Methods

### `NewHeaderAuth(config HeaderAuthConfig) *HeaderAuth`
Membuat instance baru dari HeaderAuth middleware.

### `GetHeaderName() string`
Mendapatkan nama header yang digunakan.

### `SetHeaderName(name string)`
Mengubah nama header yang digunakan.

### `Middleware() fiber.Handler`
Mengembalikan Fiber middleware handler.

## Response Codes

- `200 OK` - API key valid
- `401 Unauthorized` - API key invalid atau missing (default)
- Custom status code jika menggunakan custom error handler

## Examples

### Complete Example dengan Database

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

    // Setup key provider dengan database
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

1. **Gunakan HTTPS** - Selalu gunakan HTTPS di production untuk melindungi API keys
2. **Rotate Keys** - Rutin update API keys secara berkala
3. **Different Keys per Client** - Berikan API key yang berbeda untuk setiap client
4. **Log Access** - Gunakan success handler untuk logging
5. **Rate Limiting** - Kombinasikan dengan rate limiting middleware
6. **Monitor Usage** - Track penggunaan API key untuk mendeteksi abuse

## See Also

- [Basic Auth](./basic_auth.md)
- [QueryString Auth](./querystring_auth.md)
- [BaseKey Interface](./basekey_interface.md)
