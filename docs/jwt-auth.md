# JWT Authentication Middleware

Middleware untuk autentikasi menggunakan JSON Web Token (JWT). Mendukung ekstraksi token dari header, query string, atau cookie.

## Features

- ✅ Validasi JWT token dengan berbagai signing methods (HS256, HS384, HS512)
- ✅ Multiple token sources (Header, Query String, Cookie)
- ✅ Custom auth scheme (default: Bearer)
- ✅ Automatic token expiration checking
- ✅ Success handler untuk custom logic setelah validasi
- ✅ Custom error handler
- ✅ Claims storage di context
- ✅ Support untuk semua HTTP methods

## Installation

```go
import "github.com/budimanlai/go-pkg/middleware/auth"
```

```bash
go get github.com/golang-jwt/jwt/v5
```

## Usage

### Basic Usage dengan Bearer Token

```go
package main

import (
    "time"
    
    "github.com/budimanlai/go-pkg/middleware/auth"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

func main() {
    app := fiber.New()

    // Setup JWT middleware
    jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
        SecretKey: "your-secret-key-here",
        // Defaults:
        // SigningMethod: "HS256"
        // TokenLookup: "header:Authorization"
        // AuthScheme: "Bearer"
        // ContextKey: "user"
    })

    // Apply middleware
    app.Use(jwtAuth.Middleware())

    app.Get("/api/profile", func(c *fiber.Ctx) error {
        // Get claims from context
        claims := c.Locals("user").(jwt.MapClaims)
        userID := claims["user_id"].(string)
        
        return c.JSON(fiber.Map{
            "message": "Profile data",
            "user_id": userID,
        })
    })

    app.Listen(":3000")
}
```

**Generate JWT Token (untuk testing):**
```go
func generateToken(userID string, secretKey string) string {
    claims := jwt.MapClaims{
        "user_id": userID,
        "email":   "user@example.com",
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 jam
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte(secretKey))
    return tokenString
}
```

**Request:**
```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:3000/api/profile
```

### Token dari Query String

```go
jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey:   "your-secret-key",
    TokenLookup: "query:token",
    AuthScheme:  "", // No auth scheme for query param
})
```

**Request:**
```bash
curl http://localhost:3000/api/profile?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token dari Cookie

```go
jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey:   "your-secret-key",
    TokenLookup: "cookie:jwt",
    AuthScheme:  "", // No auth scheme for cookie
})
```

**Request:**
```bash
curl --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:3000/api/profile
```

### Dengan Success Handler

Success handler dipanggil setelah JWT berhasil divalidasi:

```go
successHandler := func(c *fiber.Ctx, claims jwt.MapClaims) error {
    // Extract user info from claims
    userID := claims["user_id"].(string)
    email := claims["email"].(string)
    role := claims["role"].(string)
    
    // Store additional data in context
    c.Locals("user_id", userID)
    c.Locals("email", email)
    c.Locals("role", role)
    
    // Check user permissions from database
    hasPermission := checkUserPermission(userID, c.Path())
    if !hasPermission {
        return fiber.NewError(fiber.StatusForbidden, "No permission")
    }
    
    // Log access
    log.Printf("User %s accessed %s", email, c.Path())
    
    return nil
}

jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey:      "your-secret-key",
    SuccessHandler: successHandler,
})
```

### Dengan Custom Error Handler

```go
customErrorHandler := func(c *fiber.Ctx, err error) error {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "success": false,
        "error":   "Unauthorized",
        "message": "Invalid or expired JWT token",
        "code":    "JWT_INVALID",
    })
}

jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey:    "your-secret-key",
    ErrorHandler: customErrorHandler,
})
```

### Custom Signing Method

```go
jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey:     "your-secret-key",
    SigningMethod: "HS512", // HS256, HS384, or HS512
})
```

### Custom Context Key

```go
jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey:  "your-secret-key",
    ContextKey: "jwt_claims", // Default: "user"
})

app.Get("/test", func(c *fiber.Ctx) error {
    claims := c.Locals("jwt_claims").(jwt.MapClaims)
    return c.JSON(claims)
})
```

### Custom Auth Scheme

```go
// Default: "Bearer"
jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey:  "your-secret-key",
    AuthScheme: "JWT", // Custom scheme
})
```

**Request:**
```bash
curl -H "Authorization: JWT eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:3000/api/profile
```

### Apply ke Specific Routes

```go
app := fiber.New()

// Setup JWT auth
jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey: "your-secret-key",
})

// Public routes
app.Post("/login", loginHandler)
app.Post("/register", registerHandler)

// Protected routes group
api := app.Group("/api")
api.Use(jwtAuth.Middleware())

api.Get("/profile", getProfile)
api.Get("/users", getUsers)
api.Post("/posts", createPost)
```

## Configuration Options

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `SecretKey` | `string` | Secret key untuk sign/validate JWT (required) | - |
| `SigningMethod` | `string` | Signing method: "HS256", "HS384", "HS512" | `"HS256"` |
| `TokenLookup` | `string` | Token location: "header:Name", "query:name", "cookie:name" | `"header:Authorization"` |
| `AuthScheme` | `string` | Authorization scheme (e.g., "Bearer") | `"Bearer"` |
| `ContextKey` | `string` | Key untuk menyimpan claims di context | `"user"` |
| `SuccessHandler` | `func` | Handler dipanggil setelah validasi berhasil | `nil` |
| `ErrorHandler` | `fiber.ErrorHandler` | Custom error handler | `nil` |
| `Claims` | `jwt.Claims` | Custom claims struct | `jwt.MapClaims{}` |

## Complete Example dengan Login

```go
package main

import (
    "log"
    "time"
    
    "github.com/budimanlai/go-pkg/middleware/auth"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

var secretKey = "my-super-secret-key-change-in-production"

// Generate JWT token
func generateToken(userID string, email string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "email":   email,
        "role":    "user",
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

// Login handler
func loginHandler(c *fiber.Ctx) error {
    type LoginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    // Validate credentials (simplified)
    if req.Email != "user@example.com" || req.Password != "password" {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
    }
    
    // Generate token
    token, err := generateToken("123", req.Email)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "token":   token,
        "type":    "Bearer",
    })
}

// Protected handler
func profileHandler(c *fiber.Ctx) error {
    claims := c.Locals("user").(jwt.MapClaims)
    
    return c.JSON(fiber.Map{
        "success": true,
        "data": fiber.Map{
            "user_id": claims["user_id"],
            "email":   claims["email"],
            "role":    claims["role"],
        },
    })
}

func main() {
    app := fiber.New()
    
    // Success handler
    successHandler := func(c *fiber.Ctx, claims jwt.MapClaims) error {
        log.Printf("User %s accessed %s", claims["email"], c.Path())
        return nil
    }
    
    // Error handler
    errorHandler := func(c *fiber.Ctx, err error) error {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "success": false,
            "error":   "Unauthorized",
            "message": err.Error(),
        })
    }
    
    // Setup JWT middleware
    jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
        SecretKey:      secretKey,
        SuccessHandler: successHandler,
        ErrorHandler:   errorHandler,
    })
    
    // Public routes
    app.Post("/login", loginHandler)
    
    // Protected routes
    api := app.Group("/api")
    api.Use(jwtAuth.Middleware())
    api.Get("/profile", profileHandler)
    
    log.Fatal(app.Listen(":3000"))
}
```

**Testing:**

1. **Login:**
```bash
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Response:
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "type": "Bearer"
}
```

2. **Access Protected Route:**
```bash
curl http://localhost:3000/api/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Response:
{
  "success": true,
  "data": {
    "user_id": "123",
    "email": "user@example.com",
    "role": "user"
  }
}
```

3. **Without Token (Error):**
```bash
curl http://localhost:3000/api/profile

# Response:
{
  "success": false,
  "error": "Unauthorized",
  "message": "missing or malformed JWT"
}
```

## JWT Claims Structure

Recommended JWT claims structure:

```go
claims := jwt.MapClaims{
    // Standard claims
    "exp": time.Now().Add(time.Hour * 24).Unix(),  // Expiration time
    "iat": time.Now().Unix(),                       // Issued at
    "nbf": time.Now().Unix(),                       // Not before
    
    // Custom claims
    "user_id": "123",
    "email":   "user@example.com",
    "role":    "admin",
    "permissions": []string{"read", "write"},
}
```

## Error Responses

Default error response (401 Unauthorized):
```json
{
  "error": "Unauthorized",
  "message": "missing or malformed JWT"
}
```

With custom error handler, you can customize the response format.

## Methods

### `NewJWTAuth(config JWTConfig) *JWTAuth`
Membuat instance baru dari JWTAuth middleware.

### `Middleware() fiber.Handler`
Mengembalikan Fiber middleware handler.

### `GetSecretKey() string`
Mendapatkan secret key yang digunakan.

### `GetSigningMethod() string`
Mendapatkan signing method yang digunakan.

### `GetContextKey() string`
Mendapatkan context key untuk menyimpan claims.

## Security Best Practices

1. **Secret Key**
   - Gunakan secret key yang kuat dan random
   - Minimal 32 karakter
   - Simpan di environment variables
   - Jangan commit ke repository

2. **Token Expiration**
   - Set expiration time yang reasonable (24 jam - 7 hari)
   - Implement refresh token mechanism
   - Force re-login untuk operasi sensitive

3. **HTTPS**
   - Selalu gunakan HTTPS di production
   - Token dapat dicuri jika menggunakan HTTP

4. **Token Storage (Client-side)**
   - Simpan di httpOnly cookie (recommended)
   - Atau localStorage dengan XSS protection
   - Hindari sessionStorage untuk persistent login

5. **Signing Method**
   - Gunakan minimal HS256
   - HS512 untuk keamanan lebih tinggi
   - Jangan gunakan "none" algorithm

6. **Claims Validation**
   - Validasi exp (expiration)
   - Validasi iss (issuer) jika menggunakan
   - Validasi aud (audience) jika menggunakan

7. **Token Revocation**
   - Implement blacklist mechanism
   - Store revoked tokens di Redis/Database
   - Check blacklist di success handler

## Performance Tips

1. **Token Size**
   - Jangan simpan data besar di JWT claims
   - JWT dikirim di setiap request
   - Simpan hanya identifier (user_id), fetch detail dari database

2. **Caching**
   - Cache user data dari database
   - Gunakan Redis untuk blacklist tokens
   - Cache permissions di memory

3. **Middleware Order**
   - Letakkan JWT middleware setelah rate limiting
   - Sebelum handler yang butuh authentication

## See Also

- [Header Auth](./header-auth.md)
- [Basic Auth](./basic_auth.md)
- [QueryString Auth](./querystring_auth.md)
- [golang-jwt/jwt Documentation](https://github.com/golang-jwt/jwt)
