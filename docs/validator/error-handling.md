# Error Handling

Comprehensive guide to handling validation errors using the `ValidationError` type.

## ValidationError Type

The `ValidationError` type provides structured error information with both backward-compatible and field-level error details.

### Definition

```go
type ValidationError struct {
    Messages []string            // All error messages (backward compatibility)
    Errors   map[string][]string // Field name -> error messages mapping
}
```

### Methods

#### Error()
Implements the `error` interface. Returns all validation error messages joined by semicolons.

```go
func (ve *ValidationError) Error() string
```

**Example:**
```go
verr := &ValidationError{
    Messages: []string{"Email is required", "Password is too short"},
}
fmt.Println(verr.Error())
// Output: Email is required; Password is too short
```

#### First()
Returns the first validation error message.

```go
func (ve *ValidationError) First() string
```

**Example:**
```go
verr := &ValidationError{
    Messages: []string{"Email is required", "Password is too short"},
}
fmt.Println(verr.First())
// Output: Email is required
```

#### All()
Returns all validation error messages as a slice.

```go
func (ve *ValidationError) All() []string
```

**Example:**
```go
verr := &ValidationError{
    Messages: []string{"Email is required", "Password is too short"},
}
for _, msg := range verr.All() {
    fmt.Println(msg)
}
// Output:
// Email is required
// Password is too short
```

#### GetFieldErrors()
Returns a map of field names to their error messages.

```go
func (ve *ValidationError) GetFieldErrors() map[string][]string
```

**Example:**
```go
verr := &ValidationError{
    Errors: map[string][]string{
        "email":    {"Email is required", "Email must be valid"},
        "password": {"Password is too short"},
    },
}

for field, errs := range verr.GetFieldErrors() {
    fmt.Printf("%s: %v\n", field, errs)
}
// Output:
// email: [Email is required Email must be valid]
// password: [Password is too short]
```

## Error Handling Patterns

### Basic Error Checking

```go
func validateUser(user User) error {
    if err := validator.ValidateStruct(user); err != nil {
        if verr, ok := err.(*validator.ValidationError); ok {
            // Handle validation error
            fmt.Println("Validation failed:", verr.First())
            return verr
        }
        // Handle other errors
        return err
    }
    return nil
}
```

### Type Assertion

Always use type assertion to access `ValidationError` methods:

```go
if err := validator.ValidateStruct(user); err != nil {
    // Type assert to ValidationError
    if verr, ok := err.(*validator.ValidationError); ok {
        // Now you can use ValidationError methods
        firstError := verr.First()
        allErrors := verr.All()
        fieldErrors := verr.GetFieldErrors()
    } else {
        // Not a validation error, handle differently
        return err
    }
}
```

### With Fiber Framework

```go
import (
    "github.com/budimanlai/go-pkg/validator"
    "github.com/budimanlai/go-pkg/response"
    "github.com/gofiber/fiber/v2"
)

app.Post("/users", func(c *fiber.Ctx) error {
    var user User
    
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    // Validate
    if err := validator.ValidateStructWithContext(c, user); err != nil {
        // ValidationErrorI18n automatically formats the error
        return response.ValidationErrorI18n(c, err)
    }
    
    // Process valid user
    return response.SuccessI18n(c, "user_created", user)
})
```

### Custom Error Response

```go
app.Post("/users", func(c *fiber.Ctx) error {
    var user User
    
    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    if err := validator.ValidateStruct(user); err != nil {
        if verr, ok := err.(*validator.ValidationError); ok {
            return c.Status(400).JSON(fiber.Map{
                "error":  verr.First(),
                "fields": verr.GetFieldErrors(),
            })
        }
        return c.Status(500).JSON(fiber.Map{
            "error": "Internal server error",
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "User created",
        "data":    user,
    })
})
```

## Error Response Formats

### Simple Error Response

```json
{
  "error": "Email is required"
}
```

**Code:**
```go
if err := validator.ValidateStruct(user); err != nil {
    if verr, ok := err.(*validator.ValidationError); ok {
        return c.JSON(fiber.Map{
            "error": verr.First(),
        })
    }
}
```

### Detailed Error Response

```json
{
  "error": "Validation failed",
  "errors": [
    "Email is required",
    "Password must be at least 8 characters",
    "Age must be greater than or equal to 18"
  ]
}
```

**Code:**
```go
if err := validator.ValidateStruct(user); err != nil {
    if verr, ok := err.(*validator.ValidationError); ok {
        return c.Status(400).JSON(fiber.Map{
            "error":  "Validation failed",
            "errors": verr.All(),
        })
    }
}
```

### Field-Level Error Response

```json
{
  "error": "Validation failed",
  "fields": {
    "email": [
      "Email is required",
      "Email must be valid"
    ],
    "password": [
      "Password must be at least 8 characters"
    ],
    "age": [
      "Age must be greater than or equal to 18"
    ]
  }
}
```

**Code:**
```go
if err := validator.ValidateStruct(user); err != nil {
    if verr, ok := err.(*validator.ValidationError); ok {
        return c.Status(400).JSON(fiber.Map{
            "error":  "Validation failed",
            "fields": verr.GetFieldErrors(),
        })
    }
}
```

### Standard Response Format (with response package)

```json
{
  "meta": {
    "success": false,
    "message": "Email is required",
    "errors": {
      "email": ["Email is required", "Email must be valid"],
      "password": ["Password must be at least 8 characters"],
      "age": ["Age must be greater than or equal to 18"]
    }
  },
  "data": null
}
```

**Code:**
```go
if err := validator.ValidateStructWithContext(c, user); err != nil {
    return response.ValidationErrorI18n(c, err)
}
```

## Advanced Error Handling

### Logging Validation Errors

```go
func validateAndLog(user User) error {
    if err := validator.ValidateStruct(user); err != nil {
        if verr, ok := err.(*validator.ValidationError); ok {
            // Log all validation errors
            log.Printf("Validation failed for user: %v", verr.All())
            
            // Log field-level errors
            for field, errs := range verr.GetFieldErrors() {
                log.Printf("Field %s: %v", field, errs)
            }
            
            return verr
        }
        return err
    }
    return nil
}
```

### Error Transformation

```go
func transformValidationError(err error) map[string]interface{} {
    if verr, ok := err.(*validator.ValidationError); ok {
        return map[string]interface{}{
            "type":    "validation_error",
            "message": verr.First(),
            "count":   len(verr.All()),
            "errors":  verr.GetFieldErrors(),
        }
    }
    return map[string]interface{}{
        "type":    "unknown_error",
        "message": err.Error(),
    }
}

// Usage
if err := validator.ValidateStruct(user); err != nil {
    errorData := transformValidationError(err)
    return c.Status(400).JSON(errorData)
}
```

### Multiple Struct Validation

```go
func validateMultiple(structs ...interface{}) error {
    var allMessages []string
    allFieldErrors := make(map[string][]string)
    
    for _, s := range structs {
        if err := validator.ValidateStruct(s); err != nil {
            if verr, ok := err.(*validator.ValidationError); ok {
                allMessages = append(allMessages, verr.All()...)
                
                for field, errs := range verr.GetFieldErrors() {
                    allFieldErrors[field] = append(allFieldErrors[field], errs...)
                }
            }
        }
    }
    
    if len(allMessages) > 0 {
        return &validator.ValidationError{
            Messages: allMessages,
            Errors:   allFieldErrors,
        }
    }
    
    return nil
}

// Usage
if err := validateMultiple(user, profile, settings); err != nil {
    return response.ValidationErrorI18n(c, err)
}
```

### Conditional Error Handling

```go
app.Post("/users", func(c *fiber.Ctx) error {
    var user User
    
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    if err := validator.ValidateStructWithContext(c, user); err != nil {
        if verr, ok := err.(*validator.ValidationError); ok {
            // Check for specific field errors
            fieldErrors := verr.GetFieldErrors()
            
            if errors, exists := fieldErrors["email"]; exists {
                log.Printf("Email validation failed: %v", errors)
            }
            
            if errors, exists := fieldErrors["password"]; exists {
                log.Printf("Password validation failed: %v", errors)
            }
            
            return response.ValidationErrorI18n(c, err)
        }
        return response.Error(c, 500, "Unexpected error")
    }
    
    return response.SuccessI18n(c, "user_created", user)
})
```

## Error Handling Best Practices

### 1. Always Type Assert

```go
// ✅ Good
if err := validator.ValidateStruct(user); err != nil {
    if verr, ok := err.(*validator.ValidationError); ok {
        return handleValidationError(verr)
    }
    return handleGenericError(err)
}

// ❌ Bad - assumes error is always ValidationError
if err := validator.ValidateStruct(user); err != nil {
    verr := err.(*validator.ValidationError) // Can panic!
    return verr.First()
}
```

### 2. Return First Error for User-Friendly Messages

```go
// ✅ Good - shows first error to user
if err := validator.ValidateStruct(user); err != nil {
    if verr, ok := err.(*validator.ValidationError); ok {
        return c.JSON(fiber.Map{
            "error": verr.First(), // User-friendly
            "details": verr.GetFieldErrors(), // For debugging
        })
    }
}
```

### 3. Log All Errors, Return First

```go
// ✅ Good - comprehensive logging
if err := validator.ValidateStruct(user); err != nil {
    if verr, ok := err.(*validator.ValidationError); ok {
        log.Printf("Validation errors: %v", verr.All()) // Log all
        return response.BadRequestI18n(c, verr.First(), nil) // Return first
    }
}
```

### 4. Use Field Errors for Forms

```go
// ✅ Good - field-level errors for forms
if err := validator.ValidateStruct(user); err != nil {
    if verr, ok := err.(*validator.ValidationError); ok {
        return c.Status(400).JSON(fiber.Map{
            "fields": verr.GetFieldErrors(), // Map to form fields
        })
    }
}
```

### 5. Use Standard Response Format

```go
// ✅ Good - consistent API response
if err := validator.ValidateStructWithContext(c, user); err != nil {
    return response.ValidationErrorI18n(c, err)
}
```

## Testing Validation Errors

```go
func TestValidateUser(t *testing.T) {
    user := User{
        Email:    "invalid-email",
        Password: "123",
        Age:      15,
    }
    
    err := validator.ValidateStruct(user)
    
    // Check error exists
    if err == nil {
        t.Fatal("Expected validation error")
    }
    
    // Type assert
    verr, ok := err.(*validator.ValidationError)
    if !ok {
        t.Fatal("Expected ValidationError type")
    }
    
    // Check error count
    if len(verr.All()) != 3 {
        t.Errorf("Expected 3 errors, got %d", len(verr.All()))
    }
    
    // Check specific field errors
    fieldErrors := verr.GetFieldErrors()
    
    if _, exists := fieldErrors["email"]; !exists {
        t.Error("Expected email error")
    }
    
    if _, exists := fieldErrors["password"]; !exists {
        t.Error("Expected password error")
    }
    
    if _, exists := fieldErrors["age"]; !exists {
        t.Error("Expected age error")
    }
}
```

## Related Documentation

- [README](README.md) - Package overview
- [Validation Tags](validation-tags.md) - Complete validation rules reference
- [I18n Integration](i18n-integration.md) - Multilingual error messages
- [Examples](examples.md) - Practical usage examples
