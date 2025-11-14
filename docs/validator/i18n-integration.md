# I18n Integration

Guide to configuring and using internationalization (i18n) for validation error messages in multiple languages.

## Overview

The validator package integrates with the i18n package to provide automatic translation of validation error messages based on user language preferences.

## Setup

### 1. Initialize I18n Manager

```go
import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/validator"
)

func main() {
    // Initialize i18n manager
    i18nMgr, err := i18n.NewI18nManager(i18n.Config{
        LocalesPath:     "./locales",
        DefaultLanguage: "en",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Set i18n manager for validator
    validator.SetI18nManager(i18nMgr)
}
```

### 2. Create Locale Files

Create JSON files in your `locales` directory with validation message translations.

**Directory structure:**
```
locales/
├── en.json
├── id.json
└── zh.json
```

## Locale File Format

All validator message keys must be prefixed with `validator.`:

### English (locales/en.json)

```json
{
  "validator.required": "{{.FieldName}} is required",
  "validator.email": "{{.FieldName}} must be a valid email address",
  "validator.min": "{{.FieldName}} must be at least {{.Param}} characters",
  "validator.max": "{{.FieldName}} must be at most {{.Param}} characters",
  "validator.gte": "{{.FieldName}} must be greater than or equal to {{.Param}}",
  "validator.lte": "{{.FieldName}} must be less than or equal to {{.Param}}",
  "validator.len": "{{.FieldName}} must be exactly {{.Param}} characters",
  "validator.numeric": "{{.FieldName}} must be numeric",
  "validator.alphanum": "{{.FieldName}} must contain only letters and numbers",
  "validator.alpha": "{{.FieldName}} must contain only letters",
  "validator.url": "{{.FieldName}} must be a valid URL",
  "validator.uri": "{{.FieldName}} must be a valid URI",
  "validator.uuid": "{{.FieldName}} must be a valid UUID",
  "validator.ipv4": "{{.FieldName}} must be a valid IPv4 address",
  "validator.ipv6": "{{.FieldName}} must be a valid IPv6 address",
  "validator.eqfield": "{{.FieldName}} must equal {{.Param}}",
  "validator.nefield": "{{.FieldName}} must not equal {{.Param}}",
  "validator.gtfield": "{{.FieldName}} must be greater than {{.Param}}",
  "validator.ltefield": "{{.FieldName}} must be less than or equal to {{.Param}}",
  "validator.oneof": "{{.FieldName}} must be one of: {{.Param}}",
  "validator.default": "{{.FieldName}} is invalid ({{.Tag}})"
}
```

### Indonesian (locales/id.json)

```json
{
  "validator.required": "{{.FieldName}} harus diisi",
  "validator.email": "{{.FieldName}} harus berupa alamat email yang valid",
  "validator.min": "{{.FieldName}} minimal {{.Param}} karakter",
  "validator.max": "{{.FieldName}} maksimal {{.Param}} karakter",
  "validator.gte": "{{.FieldName}} minimal {{.Param}}",
  "validator.lte": "{{.FieldName}} maksimal {{.Param}}",
  "validator.len": "{{.FieldName}} harus tepat {{.Param}} karakter",
  "validator.numeric": "{{.FieldName}} harus berupa angka",
  "validator.alphanum": "{{.FieldName}} hanya boleh berisi huruf dan angka",
  "validator.alpha": "{{.FieldName}} hanya boleh berisi huruf",
  "validator.url": "{{.FieldName}} harus berupa URL yang valid",
  "validator.uri": "{{.FieldName}} harus berupa URI yang valid",
  "validator.uuid": "{{.FieldName}} harus berupa UUID yang valid",
  "validator.ipv4": "{{.FieldName}} harus berupa alamat IPv4 yang valid",
  "validator.ipv6": "{{.FieldName}} harus berupa alamat IPv6 yang valid",
  "validator.eqfield": "{{.FieldName}} harus sama dengan {{.Param}}",
  "validator.nefield": "{{.FieldName}} tidak boleh sama dengan {{.Param}}",
  "validator.gtfield": "{{.FieldName}} harus lebih besar dari {{.Param}}",
  "validator.ltefield": "{{.FieldName}} harus kurang dari atau sama dengan {{.Param}}",
  "validator.oneof": "{{.FieldName}} harus salah satu dari: {{.Param}}",
  "validator.default": "{{.FieldName}} tidak valid ({{.Tag}})"
}
```

### Chinese (locales/zh.json)

```json
{
  "validator.required": "{{.FieldName}}是必填项",
  "validator.email": "{{.FieldName}}必须是有效的电子邮件地址",
  "validator.min": "{{.FieldName}}至少需要{{.Param}}个字符",
  "validator.max": "{{.FieldName}}最多{{.Param}}个字符",
  "validator.gte": "{{.FieldName}}必须大于或等于{{.Param}}",
  "validator.lte": "{{.FieldName}}必须小于或等于{{.Param}}",
  "validator.len": "{{.FieldName}}必须恰好是{{.Param}}个字符",
  "validator.numeric": "{{.FieldName}}必须是数字",
  "validator.alphanum": "{{.FieldName}}只能包含字母和数字",
  "validator.alpha": "{{.FieldName}}只能包含字母",
  "validator.url": "{{.FieldName}}必须是有效的URL",
  "validator.uri": "{{.FieldName}}必须是有效的URI",
  "validator.uuid": "{{.FieldName}}必须是有效的UUID",
  "validator.ipv4": "{{.FieldName}}必须是有效的IPv4地址",
  "validator.ipv6": "{{.FieldName}}必须是有效的IPv6地址",
  "validator.eqfield": "{{.FieldName}}必须等于{{.Param}}",
  "validator.nefield": "{{.FieldName}}不能等于{{.Param}}",
  "validator.gtfield": "{{.FieldName}}必须大于{{.Param}}",
  "validator.ltefield": "{{.FieldName}}必须小于或等于{{.Param}}",
  "validator.oneof": "{{.FieldName}}必须是以下之一：{{.Param}}",
  "validator.default": "{{.FieldName}}无效 ({{.Tag}})"
}
```

## Template Variables

Validation messages support template variables for dynamic content:

| Variable | Description | Example Value |
|----------|-------------|---------------|
| `{{.FieldName}}` | Field name from JSON tag or struct field | `email`, `password`, `age` |
| `{{.Param}}` | Parameter value from validation tag | `8` (from `min=8`), `100` (from `max=100`) |
| `{{.Tag}}` | Validation tag name | `required`, `email`, `min` |

### Example Usage

**Validation tag:**
```go
type User struct {
    Password string `json:"password" validate:"required,min=8"`
}
```

**English message:**
```json
{
  "validator.min": "{{.FieldName}} must be at least {{.Param}} characters"
}
```

**Result:** `"password must be at least 8 characters"`

**Indonesian message:**
```json
{
  "validator.min": "{{.FieldName}} minimal {{.Param}} karakter"
}
```

**Result:** `"password minimal 8 karakter"`

## Validation Methods with I18n

### ValidateStruct (Default Language)

Uses the default language from i18n configuration:

```go
user := User{
    Email:    "invalid",
    Password: "123",
}

err := validator.ValidateStruct(user)
// Returns errors in default language (e.g., "en")
```

### ValidateStructWithLang (Specific Language)

Validates with a specific language code:

```go
user := User{
    Email:    "invalid",
    Password: "123",
}

// English errors
err := validator.ValidateStructWithLang(user, "en")

// Indonesian errors
err := validator.ValidateStructWithLang(user, "id")

// Chinese errors
err := validator.ValidateStructWithLang(user, "zh")
```

### ValidateStructWithContext (Auto Language Detection)

Automatically detects language from Fiber context:

```go
app.Post("/users", func(c *fiber.Ctx) error {
    var user User
    
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    // Language automatically detected from:
    // 1. Accept-Language header
    // 2. ?lang query parameter
    // 3. Default language
    if err := validator.ValidateStructWithContext(c, user); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    return response.SuccessI18n(c, "user_created", user)
})
```

## Complete Example

### Application Setup

```go
package main

import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/response"
    "github.com/budimanlai/go-pkg/validator"
    "github.com/gofiber/fiber/v2"
    "log"
)

type User struct {
    Name     string `json:"name" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Age      int    `json:"age" validate:"required,gte=18"`
}

func main() {
    // Initialize i18n
    i18nMgr, err := i18n.NewI18nManager(i18n.Config{
        LocalesPath:     "./locales",
        DefaultLanguage: "en",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Configure packages
    validator.SetI18nManager(i18nMgr)
    response.SetI18nManager(i18nMgr)
    
    // Create Fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: response.FiberErrorHandler,
    })
    
    // Add i18n middleware
    app.Use(i18n.I18nMiddleware(i18nMgr))
    
    // Routes
    app.Post("/users", createUser)
    
    log.Fatal(app.Listen(":3000"))
}

func createUser(c *fiber.Ctx) error {
    var user User
    
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    // Validate with automatic language detection
    if err := validator.ValidateStructWithContext(c, user); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Create user...
    return response.SuccessI18n(c, "user_created", user)
}
```

### Testing with Different Languages

**English Request:**
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en" \
  -d '{"email":"invalid","password":"123","age":15}'
```

**Response:**
```json
{
  "meta": {
    "success": false,
    "message": "name is required",
    "errors": {
      "name": ["name is required"],
      "email": ["email must be a valid email address"],
      "password": ["password must be at least 8 characters"],
      "age": ["age must be greater than or equal to 18"]
    }
  },
  "data": null
}
```

**Indonesian Request:**
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -H "Accept-Language: id" \
  -d '{"email":"invalid","password":"123","age":15}'
```

**Response:**
```json
{
  "meta": {
    "success": false,
    "message": "name harus diisi",
    "errors": {
      "name": ["name harus diisi"],
      "email": ["email harus berupa alamat email yang valid"],
      "password": ["password minimal 8 karakter"],
      "age": ["age minimal 18"]
    }
  },
  "data": null
}
```

**Chinese Request:**
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -H "Accept-Language: zh" \
  -d '{"email":"invalid","password":"123","age":15}'
```

**Response:**
```json
{
  "meta": {
    "success": false,
    "message": "name是必填项",
    "errors": {
      "name": ["name是必填项"],
      "email": ["email必须是有效的电子邮件地址"],
      "password": ["password至少需要8个字符"],
      "age": ["age必须大于或等于18"]
    }
  },
  "data": null
}
```

## Fallback Behavior

### Without I18n Manager

If `SetI18nManager` is not called, validator uses default English messages:

```go
// No i18n setup
err := validator.ValidateStruct(user)
// Returns: "Email must be a valid email address"
```

### Missing Translation

If a translation key is not found, it falls back to the default message:

```json
{
  "validator.email": "{{.FieldName}} must be a valid email address"
}
```

## Field Name Extraction

The validator automatically extracts field names from JSON tags for consistent error messages:

```go
type User struct {
    Email string `json:"email" validate:"required"`
    // Error message will use "email" not "Email"
}

type Product struct {
    ProductName string `json:"product_name" validate:"required"`
    // Error message will use "product_name"
}

type Item struct {
    Name string `validate:"required"`
    // No JSON tag - will use "Name"
}
```

## Best Practices

1. **Always prefix with validator.** - All validation message keys must start with `validator.`
2. **Provide all common tags** - Include translations for commonly used tags
3. **Use template variables** - Make messages dynamic with `{{.FieldName}}`, `{{.Param}}`, etc.
4. **Consistent field names** - Use JSON tags for consistent naming
5. **Default fallback** - Always include `validator.default` for unknown tags
6. **Test all languages** - Verify translations with different Accept-Language headers
7. **Update together** - Keep all locale files in sync when adding new validation rules

## Complete Tag Reference

Common validation tags that should be translated:

```json
{
  "validator.required": "...",
  "validator.email": "...",
  "validator.min": "...",
  "validator.max": "...",
  "validator.gte": "...",
  "validator.lte": "...",
  "validator.gt": "...",
  "validator.lt": "...",
  "validator.eq": "...",
  "validator.ne": "...",
  "validator.len": "...",
  "validator.alpha": "...",
  "validator.alphanum": "...",
  "validator.numeric": "...",
  "validator.url": "...",
  "validator.uri": "...",
  "validator.uuid": "...",
  "validator.uuid4": "...",
  "validator.ipv4": "...",
  "validator.ipv6": "...",
  "validator.mac": "...",
  "validator.isbn": "...",
  "validator.datetime": "...",
  "validator.eqfield": "...",
  "validator.nefield": "...",
  "validator.gtfield": "...",
  "validator.ltefield": "...",
  "validator.oneof": "...",
  "validator.contains": "...",
  "validator.excludes": "...",
  "validator.startswith": "...",
  "validator.endswith": "...",
  "validator.unique": "...",
  "validator.default": "..."
}
```

## Related Documentation

- [README](README.md) - Package overview
- [Validation Tags](validation-tags.md) - Complete validation rules
- [Error Handling](error-handling.md) - ValidationError type
- [I18n Package](../i18n.md) - I18n configuration and setup
