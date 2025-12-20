# JWT Authentication Middleware

Middleware for authentication using JSON Web Tokens (JWT). Supports token extraction from header, query string, or cookie.

## Features

- ✅ JWT token validation with various signing methods (HS256, HS384, HS512)
- ✅ Multiple token sources (Header, Query String, Cookie)
- ✅ Custom auth scheme (default: Bearer)
- ✅ Automatic token expiration checking
- ✅ Success handler for custom logic after validation
- ✅ Custom error handler
- ✅ Claims storage in context
- ✅ Support for all HTTP methods

## Installation

```go
import "github.com/budimanlai/go-pkg/middleware/auth"
```

```bash
go get github.com/golang-jwt/jwt/v5
```

## Usage

### Basic Usage with Bearer Token

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

**Generate JWT Token (for testing):**
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

### Token from Query String

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

### Token from Cookie

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

### With Success Handler

Success handler is called after JWT is successfully validated:

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

### With Custom Error Handler

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
| `SecretKey` | `string` | Secret key for signing/validating JWT (required) | - |
| `SigningMethod` | `string` | Signing method: "HS256", "HS384", "HS512" | `"HS256"` |
| `TokenLookup` | `string` | Token location: "header:Name", "query:name", "cookie:name" | `"header:Authorization"` |
| `AuthScheme` | `string` | Authorization scheme (e.g., "Bearer") | `"Bearer"` |
| `ContextKey` | `string` | Key for storing claims in context | `"user"` |
| `SuccessHandler` | `func` | Handler called after successful validation | `nil` |
| `ErrorHandler` | `fiber.ErrorHandler` | Custom error handler | `nil` |
| `Claims` | `jwt.Claims` | Custom claims struct | `jwt.MapClaims{}` |

## Complete Example with Login

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
Creates a new instance of JWTAuth middleware.

### `Middleware() fiber.Handler`
Returns Fiber middleware handler.

### `GetSecretKey() string`
Gets the secret key being used.

### `GetSigningMethod() string`
Gets the signing method being used.

### `GetContextKey() string`
Gets the context key for storing claims.

## Security Best Practices

1. **Secret Key**
   - Use strong and random secret key
   - Minimum 32 characters
   - Store in environment variables
   - Don't commit to repository

2. **Token Expiration**
   - Set reasonable expiration time (24 hours - 7 days)
   - Implement refresh token mechanism
   - Force re-login for sensitive operations

3. **HTTPS**
   - Always use HTTPS in production
   - Tokens can be stolen if using HTTP

4. **Token Storage (Client-side)**
   - Store in httpOnly cookie (recommended)
   - Or localStorage with XSS protection
   - Avoid sessionStorage for persistent login

5. **Signing Method**
   - Use minimum HS256
   - HS512 for higher security
   - Don't use "none" algorithm

6. **Claims Validation**
   - Validate exp (expiration)
   - Validate iss (issuer) if used
   - Validate aud (audience) if used

7. **Token Revocation**
   - Implement blacklist mechanism
   - Store revoked tokens in Redis/Database
   - Check blacklist in success handler

## Performance Tips

1. **Token Size**
   - Don't store large data in JWT claims
   - JWT is sent in every request
   - Store only identifiers (user_id), fetch details from database

2. **Caching**
   - Cache user data from database
   - Use Redis for token blacklist
   - Cache permissions in memory

3. **Middleware Order**
   - Place JWT middleware after rate limiting
   - Before handlers that need authentication

## See Also

- [Header Auth](./header-auth.md)
- [Basic Auth](./basic_auth.md)
- [QueryString Auth](./querystring_auth.md)
- [golang-jwt/jwt Documentation](https://github.com/golang-jwt/jwt)
