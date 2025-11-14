# I18n Response Functions

I18n response functions provide automatic translation of messages based on the user's language preference. They integrate seamlessly with the i18n package to deliver multilingual API responses.

## Setup

Before using i18n responses, you must configure the i18n manager:

```go
import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/response"
)

// Initialize i18n manager
i18nMgr, err := i18n.NewI18nManager(i18n.Config{
    LocalesPath:     "./locales",
    DefaultLanguage: "en",
})

// Set i18n manager for response package
response.SetI18nManager(i18nMgr)

// Use i18n middleware in your Fiber app
app.Use(i18n.I18nMiddleware(i18nMgr))
```

## SetI18nManager

Configures the global i18n manager for response translations.

### Signature

```go
func SetI18nManager(manager *i18n.I18nManager)
```

### Parameters

- `manager` (*i18n.I18nManager) - The initialized i18n manager instance

### Example

```go
i18nMgr, _ := i18n.NewI18nManager(i18n.Config{
    LocalesPath:     "./locales",
    DefaultLanguage: "en",
})
response.SetI18nManager(i18nMgr)
```

## SuccessI18n

Returns a 200 OK response with a translated success message and optional data.

### Signature

```go
func SuccessI18n(c *fiber.Ctx, messageID string, data interface{}) error
```

### Parameters

- `c` (*fiber.Ctx) - The Fiber context
- `messageID` (string) - Message identifier to translate
- `data` (interface{}) - Response data to include (can be nil)

### Response Format

```json
{
  "meta": {
    "success": true,
    "message": "Translated success message"
  },
  "data": {
    // your data here
  }
}
```

### Examples

**Basic success response:**

```go
app.Post("/users", func(c *fiber.Ctx) error {
    user := createUser()
    return response.SuccessI18n(c, "user_created", user)
})
```

**Locale files (locales/en.json):**

```json
{
  "user_created": "User created successfully"
}
```

**Locale files (locales/id.json):**

```json
{
  "user_created": "Pengguna berhasil dibuat"
}
```

## ErrorI18n

Returns an error response with a translated message and custom status code.

### Signature

```go
func ErrorI18n(c *fiber.Ctx, code int, messageID string, data interface{}) error
```

### Parameters

- `c` (*fiber.Ctx) - The Fiber context
- `code` (int) - HTTP status code
- `messageID` (string) - Message identifier to translate
- `data` (interface{}) - Template data for message interpolation (can be nil)

### Response Format

```json
{
  "meta": {
    "success": false,
    "message": "Translated error message"
  },
  "data": null
}
```

### Examples

**Simple error:**

```go
app.Get("/process", func(c *fiber.Ctx) error {
    if err := process(); err != nil {
        return response.ErrorI18n(c, 500, "processing_failed", nil)
    }
    return response.SuccessI18n(c, "processing_success", nil)
})
```

**Error with template data:**

```go
app.Get("/files/:id", func(c *fiber.Ctx) error {
    file := getFile(c.Params("id"))
    if file == nil {
        return response.ErrorI18n(c, 404, "file_not_found", map[string]string{
            "Filename": c.Params("id"),
        })
    }
    return response.SuccessI18n(c, "file_retrieved", file)
})
```

**Locale files with template:**

```json
{
  "file_not_found": "File '{{.Filename}}' not found"
}
```

## BadRequestI18n

Returns a 400 Bad Request response with a translated message.

### Signature

```go
func BadRequestI18n(c *fiber.Ctx, messageID string, data interface{}) error
```

### Parameters

- `c` (*fiber.Ctx) - The Fiber context
- `messageID` (string) - Message identifier to translate
- `data` (interface{}) - Template data for message interpolation (can be nil)

### Response Format

```json
{
  "meta": {
    "success": false,
    "message": "Translated bad request message"
  },
  "data": null
}
```

### Examples

**Simple bad request:**

```go
app.Post("/users", func(c *fiber.Ctx) error {
    var user User
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    return response.SuccessI18n(c, "user_created", user)
})
```

**Bad request with data:**

```go
app.Post("/upload", func(c *fiber.Ctx) error {
    file, err := c.FormFile("document")
    if err != nil {
        return response.BadRequestI18n(c, "file_required", nil)
    }
    
    maxSize := 10 * 1024 * 1024 // 10MB
    if file.Size > int64(maxSize) {
        return response.BadRequestI18n(c, "file_too_large", map[string]string{
            "MaxSize": "10MB",
            "Size": fmt.Sprintf("%.2fMB", float64(file.Size)/1024/1024),
        })
    }
    
    return response.SuccessI18n(c, "file_uploaded", nil)
})
```

## NotFoundI18n

Returns a 404 Not Found response with a translated message.

### Signature

```go
func NotFoundI18n(c *fiber.Ctx, messageID string) error
```

### Parameters

- `c` (*fiber.Ctx) - The Fiber context
- `messageID` (string) - Message identifier to translate

### Response Format

```json
{
  "meta": {
    "success": false,
    "message": "Translated not found message"
  },
  "data": null
}
```

### Examples

**Resource not found:**

```go
app.Get("/users/:id", func(c *fiber.Ctx) error {
    user := getUserByID(c.Params("id"))
    if user == nil {
        return response.NotFoundI18n(c, "user_not_found")
    }
    return response.SuccessI18n(c, "user_retrieved", user)
})
```

**Locale files:**

```json
{
  "user_not_found": "User not found"
}
```

## ValidationErrorI18n

Returns a 400 Bad Request response with detailed validation errors. Automatically extracts field-level errors from `validator.ValidationError`.

### Signature

```go
func ValidationErrorI18n(c *fiber.Ctx, err error) error
```

### Parameters

- `c` (*fiber.Ctx) - The Fiber context
- `err` (error) - The validation error (should be *validator.ValidationError)

### Response Format

```json
{
  "meta": {
    "success": false,
    "message": "First validation error message",
    "errors": {
      "Email": ["Email is required", "Email must be valid"],
      "Password": ["Password must be at least 8 characters"]
    }
  },
  "data": null
}
```

### Examples

**With validator package:**

```go
import (
    "github.com/budimanlai/go-pkg/validator"
    "github.com/budimanlai/go-pkg/response"
)

type User struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Age      int    `json:"age" validate:"required,gte=18"`
}

app.Post("/register", func(c *fiber.Ctx) error {
    var user User
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    // Validate with context (uses language from request)
    if err := validator.ValidateStructWithContext(c, user); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Create user
    createUser(&user)
    return response.SuccessI18n(c, "user_registered", user)
})
```

**Response example (English):**

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

**Response example (Indonesian):**

```json
{
  "meta": {
    "success": false,
    "message": "Email harus diisi",
    "errors": {
      "email": ["Email harus diisi"],
      "password": ["Password minimal 8 karakter"],
      "age": ["Age minimal 18"]
    }
  },
  "data": null
}
```

## Language Detection

The i18n response functions automatically detect the user's language from the request context. The language is set by the `I18nMiddleware` based on:

1. `Accept-Language` header
2. `lang` query parameter
3. Default language from configuration

### Example Flow

```go
// Setup
i18nMgr, _ := i18n.NewI18nManager(i18n.Config{
    LocalesPath:     "./locales",
    DefaultLanguage: "en",
})
response.SetI18nManager(i18nMgr)

app.Use(i18n.I18nMiddleware(i18nMgr))

// Handler
app.Post("/users", func(c *fiber.Ctx) error {
    // Language is automatically detected from:
    // - Accept-Language: id
    // - or ?lang=id
    // - or default: en
    
    return response.SuccessI18n(c, "user_created", user)
    // Returns Indonesian if Accept-Language: id
    // Returns English if Accept-Language: en
})
```

## Locale File Structure

Organize your translation files in the `locales` directory:

```
locales/
├── en.json
├── id.json
└── zh.json
```

**locales/en.json:**

```json
{
  "user_created": "User created successfully",
  "user_updated": "User updated successfully",
  "user_deleted": "User deleted successfully",
  "user_not_found": "User not found",
  "invalid_request_body": "Invalid request body",
  "processing_failed": "Processing failed",
  "file_not_found": "File '{{.Filename}}' not found",
  "file_too_large": "File size {{.Size}} exceeds maximum allowed size {{.MaxSize}}"
}
```

**locales/id.json:**

```json
{
  "user_created": "Pengguna berhasil dibuat",
  "user_updated": "Pengguna berhasil diperbarui",
  "user_deleted": "Pengguna berhasil dihapus",
  "user_not_found": "Pengguna tidak ditemukan",
  "invalid_request_body": "Body request tidak valid",
  "processing_failed": "Pemrosesan gagal",
  "file_not_found": "File '{{.Filename}}' tidak ditemukan",
  "file_too_large": "Ukuran file {{.Size}} melebihi batas maksimum {{.MaxSize}}"
}
```

## Best Practices

1. **Consistent Message IDs** - Use snake_case for message IDs (e.g., `user_created`, `invalid_email`)
2. **Template Data** - Use meaningful template variable names (e.g., `{{.Filename}}`, `{{.MaxSize}}`)
3. **Fallback Handling** - Always set a default language in case translation is not found
4. **Middleware Order** - Place `I18nMiddleware` before your route handlers
5. **Testing** - Test your API with different `Accept-Language` headers to ensure translations work

## Related Documentation

- [I18n Package](../i18n.md) - Internationalization setup and configuration
- [Standard Responses](standard-responses.md) - Non-i18n response functions
- [Validator Package](../validator/README.md) - Struct validation with i18n support
