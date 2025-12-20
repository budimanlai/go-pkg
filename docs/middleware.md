# Middleware Package

Package middleware menyediakan berbagai authentication middleware untuk aplikasi Fiber. Semua middleware dirancang dengan fokus pada fleksibilitas, keamanan, dan performa.

## Overview

Middleware package mendukung berbagai metode autentikasi:

- **JWT Authentication** - Token-based authentication dengan JSON Web Tokens
- **Header Authentication** - API key authentication via headers (X-API-Key)
- **Basic Authentication** - HTTP Basic Authentication (username/password)
- **Query String Authentication** - API key via query parameters
- **Database API Key** - Database-backed API key management

## Features

- ✅ Multiple authentication methods
- ✅ Flexible key storage (in-memory, database)
- ✅ Thread-safe operations
- ✅ Custom success/error handlers
- ✅ Context integration
- ✅ Full test coverage (87 tests)
- ✅ Production-ready

## Installation

```bash
go get github.com/budimanlai/go-pkg/middleware/auth
go get github.com/gofiber/fiber/v2
go get github.com/golang-jwt/jwt/v5
go get gorm.io/gorm
```

## Quick Start

### JWT Authentication

```go
import "github.com/budimanlai/go-pkg/middleware/auth"

jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey: "your-secret-key",
})

app.Use(jwtAuth.Middleware())
```

### Header Authentication (X-API-Key)

```go
keyProvider := auth.NewBaseKeyProvider()
keyProvider.Add("secret-api-key-123")

headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
    KeyProvider: keyProvider,
    HeaderName:  "X-API-Key",
})

app.Use(headerAuth.Middleware())
```

### Basic Authentication

```go
keyProvider := auth.NewBaseKeyProvider()
keyProvider.AddKeyValue("admin", "password123")

basicAuth := auth.NewBasicAuth(auth.BasicAuthConfig{
    KeyProvider: keyProvider,
})

app.Use(basicAuth.Middleware())
```

## Authentication Methods

### 1. JWT Authentication

JSON Web Token authentication mendukung multiple token sources (header, query, cookie).

**Documentation:** [JWT Auth Guide](./jwt-auth.md)

**Example:**
```go
jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey:     "my-secret-key",
    SigningMethod: "HS256",
    TokenLookup:   "header:Authorization",
    AuthScheme:    "Bearer",
    ContextKey:    "user",
})

app.Use(jwtAuth.Middleware())

app.Get("/profile", func(c *fiber.Ctx) error {
    claims := c.Locals("user").(jwt.MapClaims)
    return c.JSON(claims)
})
```

**Token Sources:**
- Header: `Authorization: Bearer <token>`
- Query: `?token=<token>`
- Cookie: `Cookie: jwt=<token>`

### 2. Header Authentication

API key authentication melalui HTTP headers (default: X-API-Key).

**Documentation:** [Header Auth Guide](./header-auth.md)

**Example:**
```go
keyProvider := auth.NewBaseKeyProvider()
keyProvider.Add("api-key-123")
keyProvider.Add("api-key-456")

headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
    KeyProvider: keyProvider,
    HeaderName:  "X-API-Key",
})

app.Use(headerAuth.Middleware())
```

**Request:**
```bash
curl -H "X-API-Key: api-key-123" http://localhost:3000/api/data
```

### 3. Basic Authentication

HTTP Basic Authentication dengan username dan password.

**Example:**
```go
keyProvider := auth.NewBaseKeyProvider()
keyProvider.AddKeyValue("admin", "secure-password")
keyProvider.AddKeyValue("user", "user-password")

basicAuth := auth.NewBasicAuth(auth.BasicAuthConfig{
    KeyProvider: keyProvider,
})

app.Use(basicAuth.Middleware())
```

**Request:**
```bash
curl -u admin:secure-password http://localhost:3000/api/data
```

### 4. Query String Authentication

API key authentication via query parameters.

**Example:**
```go
keyProvider := auth.NewBaseKeyProvider()
keyProvider.Add("secret-token-123")

queryAuth := auth.NewQueryStringAuth(auth.QueryStringAuthConfig{
    KeyProvider: keyProvider,
    ParamName:   "access-token",
})

app.Use(queryAuth.Middleware())
```

**Request:**
```bash
curl http://localhost:3000/api/data?access-token=secret-token-123
```

### 5. Database API Key

Database-backed API key storage menggunakan GORM.

**Example:**
```go
import "gorm.io/gorm"

db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
db.AutoMigrate(&auth.ApiKey{})

dbKeyProvider := auth.NewDbApiKey(db)
dbKeyProvider.Add("api-key-123")

headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
    KeyProvider: dbKeyProvider,
})

app.Use(headerAuth.Middleware())
```

## Key Providers

### In-Memory Provider (BaseKeyProvider)

Thread-safe in-memory key storage.

```go
keyProvider := auth.NewBaseKeyProvider()

// Add API key only
keyProvider.Add("api-key-123")

// Add key-value pair (for Basic Auth)
keyProvider.AddKeyValue("username", "password")

// Check if key exists
exists := keyProvider.IsExists("api-key-123")

// Get value for key
value, err := keyProvider.GetValue("username")

// Remove key
keyProvider.Remove("api-key-123")

// Remove all keys
keyProvider.RemoveAll()

// Replace key
keyProvider.Replace("old-key", "new-key")
```

### Database Provider (DbApiKey)

Database-backed key storage dengan status management.

```go
dbProvider := auth.NewDbApiKey(db)

// Add new API key
dbProvider.Add("new-api-key")

// Add with auth key (for key-value)
dbProvider.AddKeyValue("username", "password")

// Check if exists and active
exists := dbProvider.IsExists("new-api-key")

// Get auth key (password)
authKey, err := dbProvider.GetValue("username")

// Deactivate key (soft delete)
dbProvider.Remove("new-api-key")

// Delete all inactive keys
dbProvider.RemoveAll()

// Replace key
dbProvider.Replace("old-key", "new-key")
```

**Database Schema:**
```sql
CREATE TABLE api_key (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    api_key VARCHAR(255) NOT NULL UNIQUE,
    auth_key VARCHAR(255),
    status ENUM('active', 'inactive') DEFAULT 'active',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

## Custom Handlers

### Success Handler

Dipanggil setelah autentikasi berhasil:

```go
successHandler := func(c *fiber.Ctx, key string) error {
    // Log access
    log.Printf("User with key %s accessed %s", key, c.Path())
    
    // Store additional data
    c.Locals("api_key", key)
    c.Locals("authenticated", true)
    
    // Check permissions
    if !hasPermission(key, c.Path()) {
        return fiber.NewError(fiber.StatusForbidden, "No permission")
    }
    
    return nil
}

headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
    KeyProvider:    keyProvider,
    SuccessHandler: successHandler,
})
```

### Error Handler

Dipanggil ketika autentikasi gagal:

```go
errorHandler := func(c *fiber.Ctx, err error) error {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "success": false,
        "error":   "Unauthorized",
        "message": "Invalid or missing API key",
        "code":    "AUTH_FAILED",
    })
}

headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
    KeyProvider:  keyProvider,
    ErrorHandler: errorHandler,
})
```

## Multiple Authentication Methods

Kombinasi berbagai authentication methods:

```go
app := fiber.New()

// Public routes
app.Post("/login", loginHandler)
app.Post("/register", registerHandler)

// API Key authentication
apiGroup := app.Group("/api")
apiGroup.Use(headerAuth.Middleware())
apiGroup.Get("/public-data", publicDataHandler)

// JWT authentication for user endpoints
userGroup := app.Group("/user")
userGroup.Use(jwtAuth.Middleware())
userGroup.Get("/profile", profileHandler)
userGroup.Put("/profile", updateProfileHandler)

// Admin routes with Basic Auth
adminGroup := app.Group("/admin")
adminGroup.Use(basicAuth.Middleware())
adminGroup.Get("/users", listUsersHandler)
adminGroup.Delete("/users/:id", deleteUserHandler)
```

## Rate Limiting Integration

Kombinasi dengan rate limiting:

```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

// Rate limiter
limiter := limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
})

// Apply rate limiter first, then auth
app.Use(limiter)
app.Use(jwtAuth.Middleware())
```

## Dynamic Key Management

### Runtime Key Updates

```go
// Add new key at runtime
keyProvider.Add("new-api-key-789")

// Remove compromised key
keyProvider.Remove("compromised-key-123")

// JWT secret key rotation
jwtAuth.SetSecretKey("new-secret-key")
```

### Database Sync

```go
// Periodic sync from database
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        keys := loadKeysFromDatabase()
        for _, key := range keys {
            keyProvider.Add(key)
        }
    }
}()
```

## Complete Example

```go
package main

import (
    "log"
    "time"
    
    "github.com/budimanlai/go-pkg/middleware/auth"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    app := fiber.New()
    
    // Database setup
    db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/dbname"), &gorm.Config{})
    db.AutoMigrate(&auth.ApiKey{})
    
    // Key providers
    memoryProvider := auth.NewBaseKeyProvider()
    memoryProvider.AddKeyValue("admin", "admin123")
    
    dbProvider := auth.NewDbApiKey(db)
    dbProvider.Add("db-api-key-123")
    
    // JWT Auth
    jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
        SecretKey: "my-secret-key",
        SuccessHandler: func(c *fiber.Ctx, claims jwt.MapClaims) error {
            log.Printf("JWT user %s accessed %s", claims["user_id"], c.Path())
            return nil
        },
    })
    
    // Header Auth
    headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
        KeyProvider: dbProvider,
        HeaderName:  "X-API-Key",
    })
    
    // Basic Auth
    basicAuth := auth.NewBasicAuth(auth.BasicAuthConfig{
        KeyProvider: memoryProvider,
    })
    
    // Public routes
    app.Post("/login", loginHandler)
    
    // JWT protected routes
    user := app.Group("/user")
    user.Use(jwtAuth.Middleware())
    user.Get("/profile", profileHandler)
    
    // API key protected routes
    api := app.Group("/api")
    api.Use(headerAuth.Middleware())
    api.Get("/data", dataHandler)
    
    // Admin routes with Basic Auth
    admin := app.Group("/admin")
    admin.Use(basicAuth.Middleware())
    admin.Get("/users", listUsersHandler)
    
    log.Fatal(app.Listen(":3000"))
}

func loginHandler(c *fiber.Ctx) error {
    // Generate JWT token
    claims := jwt.MapClaims{
        "user_id": "123",
        "email":   "user@example.com",
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte("my-secret-key"))
    
    return c.JSON(fiber.Map{
        "token": tokenString,
    })
}

func profileHandler(c *fiber.Ctx) error {
    claims := c.Locals("user").(jwt.MapClaims)
    return c.JSON(claims)
}

func dataHandler(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"data": "sensitive information"})
}

func listUsersHandler(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"users": []string{"user1", "user2"}})
}
```

## Testing

Package ini memiliki comprehensive test coverage:

```bash
# Run all tests
go test ./middleware/auth/...

# Run with coverage
go test -cover ./middleware/auth/...

# Run specific test
go test -v ./middleware/auth/... -run TestJWTAuth

# Run with race detection
go test -race ./middleware/auth/...
```

**Test Statistics:**
- Total tests: 87
- Coverage: High
- Thread-safety tested: Yes
- Edge cases covered: Yes

## Security Best Practices

1. **Store Secrets Securely**
   - Use environment variables
   - Never commit secrets to repository
   - Rotate keys regularly

2. **Use HTTPS in Production**
   - Always use HTTPS in production
   - Tokens can be intercepted over HTTP

3. **Implement Rate Limiting**
   - Prevent brute force attacks
   - Use Fiber's rate limiter

4. **Validate Token Expiration**
   - Set reasonable expiration times
   - Implement token refresh mechanism

5. **Log Authentication Events**
   - Log successful and failed attempts
   - Monitor for suspicious activity

6. **Use Strong Secrets**
   - Minimum 32 characters for JWT secrets
   - Use cryptographically random keys

7. **Implement Token Revocation**
   - Blacklist compromised tokens
   - Use Redis for distributed blacklist

## Performance Tips

1. **Use In-Memory Provider** for small key sets
2. **Cache Database Queries** for DbApiKey
3. **Use Connection Pooling** for database
4. **Minimize Success Handler** logic
5. **Use Context Timeout** for long operations

## Documentation

- [JWT Authentication](./jwt-auth.md)
- [Header Authentication](./header-auth.md)
- [Response Package](./response/README.md)

## See Also

- [Security Package](./security.md)
- [Response Package](./response/README.md)
- [Databases Package](./databases.md)
