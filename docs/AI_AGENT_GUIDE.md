# AI Agent Development Guide

This document is specifically designed for AI coding agents to understand and utilize the `go-pkg` library effectively in development projects.

## Package Overview

`go-pkg` is a comprehensive Go utility library providing:
- **Response handling** with i18n support
- **Validation** with multilingual error messages
- **Security** utilities (password hashing)
- **Internationalization (i18n)** with go-i18n
- **Database** management (GORM)
- **Helper** utilities (JSON, pointers, strings, IDs)
- **Custom types** (UTCTime)
- **Logger** utilities

## Import Path

```go
import "github.com/budimanlai/go-pkg/<package>"
```

---

## Quick Reference by Use Case

### Authentication & Security

```go
import "github.com/budimanlai/go-pkg/security"

// Hash password during registration
hashedPassword := security.HashPassword("userPassword123")
// Returns: bcrypt hash string or empty string on error

// Verify password during login
isValid, err := security.CheckPasswordHash("userPassword123", hashedPassword)
// Returns: (true, nil) if match, (false, nil) if no match, (false, error) on error
```

### HTTP Response Handling (Fiber)

```go
import (
    "github.com/budimanlai/go-pkg/response"
    "github.com/gofiber/fiber/v2"
)

// Setup (do once at app initialization)
response.SetI18nManager(i18nManager) // Optional, for i18n support

// Standard responses
response.Success(c, "Operation successful", data)
response.Error(c, "Something went wrong", fiber.StatusInternalServerError)
response.BadRequest(c, "Invalid input")
response.Unauthorized(c, "Authentication required")
response.NotFound(c, "Resource not found")

// With i18n (auto-translates message key)
response.SuccessWithI18n(c, "success.created", data)
response.ErrorWithI18n(c, "error.validation_failed", fiber.StatusBadRequest)
response.UnauthorizedWithI18n(c, "error.unauthorized")
response.NotFoundWithI18n(c, "error.not_found")
```

### Validation

```go
import "github.com/budimanlai/go-pkg/validator"

// Define struct with validation tags
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Age      int    `json:"age" validate:"required,min=18,max=100"`
}

// Setup validator (do once at app initialization)
v := validator.NewValidator()
validator.SetI18nManager(i18nManager) // Optional, for i18n error messages

// Validate
var req RegisterRequest
if err := c.BodyParser(&req); err != nil {
    return response.BadRequest(c, "Invalid JSON")
}

if err := v.ValidateStruct(req); err != nil {
    if validationErr, ok := err.(*validator.ValidationError); ok {
        return c.Status(fiber.StatusBadRequest).JSON(validationErr)
    }
    return response.Error(c, err.Error(), fiber.StatusBadRequest)
}

// Available validation tags:
// required, email, min, max, len, eq, ne, gt, gte, lt, lte,
// alpha, alphanum, numeric, url, uuid, datetime, oneof
```

### Internationalization (i18n)

```go
import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/gofiber/fiber/v2"
    "golang.org/x/text/language"
)

// Setup (do once at app initialization)
config := i18n.I18nConfig{
    DefaultLanguage: language.English,
    SupportedLangs:  []string{"en", "id", "zh"},
    LocalesPath:     "./locales", // JSON files: en.json, id.json, zh.json
}
i18nManager, err := i18n.NewI18nManager(config)

// Add middleware to Fiber app
app.Use(i18n.I18nMiddleware(config))

// Translate in handler
func handler(c *fiber.Ctx) error {
    lang := i18n.GetLanguage(c) // Gets language from context
    message := i18nManager.Translate(lang, "welcome.message")
    return c.JSON(fiber.Map{"message": message})
}

// Locale file format (locales/en.json):
{
    "welcome.message": "Welcome to our application",
    "error.not_found": "Resource not found",
    "success.created": "Created successfully"
}
```

### Database Management (GORM)

```go
import "github.com/budimanlai/go-pkg/databases"

// Setup database connection
config := databases.DbConfig{
    Host:     "localhost",
    Port:     "3306",
    User:     "root",
    Password: "password",
    Database: "mydb",
    Driver:   "mysql", // or "postgres"
    
    // Optional connection pool settings
    MaxIdleConns:    10,
    MaxOpenConns:    100,
    ConnMaxLifetime: 3600, // seconds
}

dbManager, err := databases.NewDbManager(config)
if err != nil {
    log.Fatal(err)
}

// Get GORM DB instance
db := dbManager.GetDb()

// Use GORM as usual
var users []User
db.Find(&users)
```

### Helper Functions

#### JSON Operations

```go
import "github.com/budimanlai/go-pkg/helpers"

// Unmarshal JSON string to struct (generic)
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

jsonStr := `{"name":"John","email":"john@example.com"}`
user, err := helpers.UnmarshalTo[User](jsonStr)

// Unmarshal map to struct (generic)
dataMap := map[string]interface{}{
    "name":  "Jane",
    "email": "jane@example.com",
}
user2, err := helpers.UnmarshalFromMap[User](dataMap)
```

#### Pointer Operations

```go
import "github.com/budimanlai/go-pkg/helpers"

// Create pointer (generic)
intPtr := helpers.Pointer(42)           // *int
strPtr := helpers.Pointer("hello")      // *string
boolPtr := helpers.Pointer(true)        // *bool

// Safely dereference with default (generic)
value := helpers.DerefPointer(intPtr, 0)      // Returns 42
value2 := helpers.DerefPointer(nil, 0)        // Returns 0 (default)
str := helpers.DerefPointer(strPtr, "default") // Returns "hello"
```

#### String & ID Generation

```go
import "github.com/budimanlai/go-pkg/helpers"

// Generate transaction ID (YYMMDDHHMMSS + 4 random digits)
trxID := helpers.GenerateTrxID()
// Output: "2411201534251234"

// With prefix/suffix
paymentID := helpers.GenerateTrxIDWithPrefix("PAY-")
// Output: "PAY-2411201534251234"

orderID := helpers.GenerateTrxIDWithSuffix("-ORD")
// Output: "2411201534251234-ORD"

// UUID v4
msgID := helpers.GenerateMessageID()
// Output: "550e8400-e29b-41d4-a716-446655440000"

// Short unique ID (first 8 chars of UUID)
shortID := helpers.GenerateUniqueID()
// Output: "550e8400"

// Random alphanumeric string
token := helpers.GenerateRandomString(32)
// Output: "aZ3bC9dE1FgH2iJ4kL5mN6oP7qR8sT9u"

// Normalize phone number
phone := helpers.NormalizePhoneNumber("+628123456789")
// Output: "628123456789"
phone2 := helpers.NormalizePhoneNumber("08123456789")
// Output: "628123456789" (adds country code)
```

#### Date Operations

```go
import "github.com/budimanlai/go-pkg/helpers"

// Parse date string to time.Time
date, err := helpers.StringToDate("2024-11-20")
// Returns: time.Time object or error
```

### Custom Types

#### UTCTime

```go
import "github.com/budimanlai/go-pkg/types"

// Use in structs for consistent UTC JSON serialization
type Event struct {
    ID        int            `json:"id"`
    Name      string         `json:"name"`
    CreatedAt types.UTCTime  `json:"created_at"`
    UpdatedAt types.UTCTime  `json:"updated_at"`
}

event := Event{
    ID:        1,
    Name:      "Conference",
    CreatedAt: types.UTCTime(time.Now()),
    UpdatedAt: types.UTCTime(time.Now()),
}

// JSON output automatically formats as UTC ISO8601:
// {"id":1,"name":"Conference","created_at":"2025-11-20T10:30:45Z","updated_at":"2025-11-20T10:30:45Z"}
```

---

## Common Patterns

### REST API Endpoint Pattern

```go
import (
    "github.com/budimanlai/go-pkg/response"
    "github.com/budimanlai/go-pkg/validator"
    "github.com/budimanlai/go-pkg/security"
    "github.com/gofiber/fiber/v2"
)

// Create user endpoint
func CreateUser(c *fiber.Ctx) error {
    // 1. Define request struct with validation
    type CreateUserRequest struct {
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required,min=8"`
        Name     string `json:"name" validate:"required,min=2"`
    }
    
    // 2. Parse request body
    var req CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    // 3. Validate
    v := validator.NewValidator()
    if err := v.ValidateStruct(req); err != nil {
        if validationErr, ok := err.(*validator.ValidationError); ok {
            return c.Status(fiber.StatusBadRequest).JSON(validationErr)
        }
        return response.Error(c, err.Error(), fiber.StatusBadRequest)
    }
    
    // 4. Hash password
    hashedPassword := security.HashPassword(req.Password)
    if hashedPassword == "" {
        return response.Error(c, "Failed to hash password", fiber.StatusInternalServerError)
    }
    
    // 5. Save to database
    user := User{
        Email:    req.Email,
        Password: hashedPassword,
        Name:     req.Name,
    }
    if err := db.Create(&user).Error; err != nil {
        return response.Error(c, "Failed to create user", fiber.StatusInternalServerError)
    }
    
    // 6. Return success response
    return response.SuccessWithI18n(c, "success.user_created", user)
}
```

### Login Endpoint Pattern

```go
func Login(c *fiber.Ctx) error {
    type LoginRequest struct {
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required"`
    }
    
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    // Validate
    v := validator.NewValidator()
    if err := v.ValidateStruct(req); err != nil {
        if validationErr, ok := err.(*validator.ValidationError); ok {
            return c.Status(fiber.StatusBadRequest).JSON(validationErr)
        }
        return response.Error(c, err.Error(), fiber.StatusBadRequest)
    }
    
    // Find user
    var user User
    if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
        return response.UnauthorizedWithI18n(c, "error.invalid_credentials")
    }
    
    // Verify password
    isValid, err := security.CheckPasswordHash(req.Password, user.Password)
    if err != nil {
        return response.Error(c, "Authentication error", fiber.StatusInternalServerError)
    }
    if !isValid {
        return response.UnauthorizedWithI18n(c, "error.invalid_credentials")
    }
    
    // Generate token (JWT or session)
    token := generateJWT(user.ID) // Your JWT implementation
    
    return response.Success(c, "Login successful", fiber.Map{
        "token": token,
        "user":  user,
    })
}
```

### Application Initialization Pattern

```go
func main() {
    // 1. Setup i18n
    i18nConfig := i18n.I18nConfig{
        DefaultLanguage: language.English,
        SupportedLangs:  []string{"en", "id"},
        LocalesPath:     "./locales",
    }
    i18nManager, err := i18n.NewI18nManager(i18nConfig)
    if err != nil {
        log.Fatal("Failed to initialize i18n:", err)
    }
    
    // 2. Setup database
    dbConfig := databases.DbConfig{
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        Database: os.Getenv("DB_NAME"),
        Driver:   "mysql",
    }
    dbManager, err := databases.NewDbManager(dbConfig)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    db := dbManager.GetDb()
    
    // 3. Auto-migrate models
    db.AutoMigrate(&User{}, &Product{}, &Order{})
    
    // 4. Setup Fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: response.CustomErrorHandler,
    })
    
    // 5. Add middleware
    app.Use(i18n.I18nMiddleware(i18nConfig))
    
    // 6. Setup response and validator with i18n
    response.SetI18nManager(i18nManager)
    validator.SetI18nManager(i18nManager)
    
    // 7. Register routes
    app.Post("/register", CreateUser)
    app.Post("/login", Login)
    app.Get("/users/:id", GetUser)
    
    // 8. Start server
    log.Fatal(app.Listen(":3000"))
}
```

---

## Critical Rules for AI Agents

### ✅ DO:

1. **Always validate user input** using `validator.ValidateStruct()`
2. **Hash passwords** with `security.HashPassword()` before storing
3. **Verify passwords** with `security.CheckPasswordHash()` during authentication
4. **Use i18n** for user-facing messages with `response.*WithI18n()` functions
5. **Check errors** from all functions, especially validation and database operations
6. **Use generics** for type-safe operations: `UnmarshalTo[T]`, `Pointer[T]`, `DerefPointer[T]`
7. **Normalize phone numbers** before storage: `helpers.NormalizePhoneNumber()`
8. **Use UTCTime** for consistent timezone handling in JSON APIs
9. **Setup i18n manager** before using response/validator i18n features
10. **Use appropriate response functions** based on HTTP status (Success, Error, NotFound, etc.)

### ❌ DON'T:

1. **Don't store plain text passwords** - always use `security.HashPassword()`
2. **Don't skip validation** - validate all user input
3. **Don't hardcode error messages** - use i18n keys for multilingual support
4. **Don't return raw validation errors** - wrap in `ValidationError` type
5. **Don't use deprecated functions** - only use functions documented here
6. **Don't compare password hashes manually** - use `security.CheckPasswordHash()`
7. **Don't forget error handling** - check all returns, especially `err != nil`
8. **Don't use type-specific pointer functions** - they don't exist (use generic `Pointer[T]`)
9. **Don't use ToJSON/FromJSON/IsJSON** - these don't exist (use `UnmarshalTo[T]`)
10. **Don't mix response types** - use consistent response format throughout API

---

## Function Reference Quick Lookup

### Security (`security`)
| Function | Purpose | Returns |
|----------|---------|---------|
| `HashPassword(password string)` | Hash password with bcrypt | `string` (hash or empty on error) |
| `CheckPasswordHash(password, hash string)` | Verify password | `(bool, error)` |

### Response (`response`)
| Function | Purpose | Status Code |
|----------|---------|-------------|
| `Success(c, message, data)` | Success response | 200 |
| `Error(c, message, status)` | Error response | Custom |
| `BadRequest(c, message)` | Bad request | 400 |
| `Unauthorized(c, message)` | Unauthorized | 401 |
| `NotFound(c, message)` | Not found | 404 |
| `SuccessWithI18n(c, key, data)` | Success with i18n | 200 |
| `ErrorWithI18n(c, key, status)` | Error with i18n | Custom |
| `UnauthorizedWithI18n(c, key)` | Unauthorized with i18n | 401 |
| `NotFoundWithI18n(c, key)` | Not found with i18n | 404 |

### Validator (`validator`)
| Function | Purpose | Returns |
|----------|---------|---------|
| `NewValidator()` | Create validator | `*Validator` |
| `ValidateStruct(s interface{})` | Validate struct | `error` (nil or ValidationError) |
| `SetI18nManager(manager)` | Setup i18n | `void` |

### Helpers (`helpers`)
| Category | Function | Purpose |
|----------|----------|---------|
| **JSON** | `UnmarshalTo[T](jsonStr)` | JSON string to type T |
| | `UnmarshalFromMap[T](map)` | Map to type T |
| **Pointer** | `Pointer[T](value)` | Create pointer to value |
| | `DerefPointer[T](ptr, default)` | Safe dereference |
| **ID Gen** | `GenerateTrxID()` | Transaction ID |
| | `GenerateTrxIDWithPrefix(prefix)` | Transaction ID with prefix |
| | `GenerateTrxIDWithSuffix(suffix)` | Transaction ID with suffix |
| | `GenerateMessageID()` | UUID v4 |
| | `GenerateUniqueID()` | Short unique ID (8 chars) |
| | `GenerateRandomString(length)` | Random alphanumeric |
| **String** | `NormalizePhoneNumber(phone)` | Normalize to intl format |
| **Date** | `StringToDate(dateStr)` | Parse YYYY-MM-DD to time.Time |

### Database (`databases`)
| Function | Purpose | Returns |
|----------|---------|---------|
| `NewDbManager(config)` | Connect to database | `(*DbManager, error)` |
| `GetDb()` | Get GORM instance | `*gorm.DB` |

### I18n (`i18n`)
| Function | Purpose | Returns |
|----------|---------|---------|
| `NewI18nManager(config)` | Create i18n manager | `(*I18nManager, error)` |
| `Translate(lang, key)` | Translate message | `string` |
| `GetLanguage(c)` | Get language from context | `language.Tag` |
| `I18nMiddleware(config)` | Fiber middleware | `fiber.Handler` |

---

## Validation Tags Reference

Common validation tags for struct fields:

```go
// Required fields
`validate:"required"`

// String validations
`validate:"email"`           // Valid email
`validate:"min=8"`           // Minimum length 8
`validate:"max=100"`         // Maximum length 100
`validate:"len=10"`          // Exact length 10
`validate:"alpha"`           // Only letters
`validate:"alphanum"`        // Letters and numbers
`validate:"url"`             // Valid URL

// Number validations
`validate:"min=18"`          // Minimum value 18
`validate:"max=100"`         // Maximum value 100
`validate:"gt=0"`            // Greater than 0
`validate:"gte=1"`           // Greater than or equal 1
`validate:"lt=100"`          // Less than 100
`validate:"lte=99"`          // Less than or equal 99

// Special validations
`validate:"oneof=male female"` // Must be one of values
`validate:"uuid"`            // Valid UUID
`validate:"datetime"`        // Valid datetime
`validate:"numeric"`         // Only numbers

// Multiple rules (comma-separated)
`validate:"required,email,min=5,max=100"`
```

---

## Error Handling Patterns

### Validation Errors

```go
if err := v.ValidateStruct(req); err != nil {
    // Check if it's a validation error
    if validationErr, ok := err.(*validator.ValidationError); ok {
        // Return structured validation error
        return c.Status(fiber.StatusBadRequest).JSON(validationErr)
    }
    // Other error
    return response.Error(c, err.Error(), fiber.StatusBadRequest)
}
```

### Database Errors

```go
if err := db.Create(&user).Error; err != nil {
    // Check for specific database errors
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return response.NotFound(c, "User not found")
    }
    if errors.Is(err, gorm.ErrDuplicatedKey) {
        return response.BadRequest(c, "Email already exists")
    }
    return response.Error(c, "Database error", fiber.StatusInternalServerError)
}
```

### Security Errors

```go
hashedPassword := security.HashPassword(password)
if hashedPassword == "" {
    return response.Error(c, "Failed to process password", fiber.StatusInternalServerError)
}

isValid, err := security.CheckPasswordHash(password, hash)
if err != nil {
    return response.Error(c, "Authentication error", fiber.StatusInternalServerError)
}
if !isValid {
    return response.Unauthorized(c, "Invalid credentials")
}
```

---

## Locale File Structure

Create JSON files in `./locales/` directory:

**locales/en.json:**
```json
{
    "success.created": "Created successfully",
    "success.updated": "Updated successfully",
    "success.deleted": "Deleted successfully",
    "success.user_created": "User registered successfully",
    "error.not_found": "Resource not found",
    "error.unauthorized": "You are not authorized",
    "error.invalid_credentials": "Invalid email or password",
    "error.validation_failed": "Validation failed",
    "validation.required": "{{.Field}} is required",
    "validation.email": "{{.Field}} must be a valid email",
    "validation.min": "{{.Field}} must be at least {{.Param}} characters"
}
```

**locales/id.json:**
```json
{
    "success.created": "Berhasil dibuat",
    "success.updated": "Berhasil diperbarui",
    "success.deleted": "Berhasil dihapus",
    "success.user_created": "Pengguna berhasil didaftarkan",
    "error.not_found": "Sumber daya tidak ditemukan",
    "error.unauthorized": "Anda tidak memiliki akses",
    "error.invalid_credentials": "Email atau password salah",
    "error.validation_failed": "Validasi gagal",
    "validation.required": "{{.Field}} wajib diisi",
    "validation.email": "{{.Field}} harus berupa email yang valid",
    "validation.min": "{{.Field}} minimal {{.Param}} karakter"
}
```

---

## Complete Example Project

See the [examples in response documentation](response/examples.md) and [validator examples](validator/examples.md) for complete working examples including:
- User registration and authentication
- CRUD operations with validation
- File upload handling
- Pagination
- Complex nested validation
- Multilingual APIs

---

## Dependencies

Required Go modules:
```bash
go get github.com/gofiber/fiber/v2
go get github.com/go-playground/validator/v10
go get github.com/nicksnyder/go-i18n/v2
go get golang.org/x/text/language
go get golang.org/x/crypto/bcrypt
go get github.com/google/uuid
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get gorm.io/driver/postgres
```

---

## Testing

Always run tests before deploying:
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Specific package tests
go test ./response/ -v
go test ./validator/ -v
go test ./security/ -v
```

---

## Support Documentation

For detailed documentation, refer to:
- [Response Package](response/README.md)
- [Validator Package](validator/README.md)
- [Security Package](security.md)
- [Helpers Package](helpers.md)
- [I18n Package](i18n.md)
- [Database Package](databases.md)
- [Types Package](types.md)

---

**Last Updated:** November 20, 2025  
**Package Version:** Compatible with Go 1.18+
