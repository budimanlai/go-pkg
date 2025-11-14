# Validator Package Examples

Practical examples and common patterns for using the validator package in real-world applications.

## Table of Contents

- [Basic Validation](#basic-validation)
- [User Registration](#user-registration)
- [Complex Structs](#complex-structs)
- [Password Validation](#password-validation)
- [Custom Error Messages](#custom-error-messages)
- [Nested Structures](#nested-structures)
- [Conditional Validation](#conditional-validation)
- [File Upload Validation](#file-upload-validation)

## Basic Validation

Simple struct validation with common rules:

```go
package main

import (
    "fmt"
    "github.com/budimanlai/go-pkg/validator"
)

type User struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Age      int    `json:"age" validate:"required,gte=18"`
}

func main() {
    // Valid user
    user := User{
        Email:    "john@example.com",
        Password: "secure123",
        Age:      25,
    }
    
    if err := validator.ValidateStruct(user); err != nil {
        fmt.Println("Validation failed:", err)
    } else {
        fmt.Println("Validation passed")
    }
    
    // Invalid user
    invalidUser := User{
        Email:    "invalid-email",
        Password: "123",
        Age:      15,
    }
    
    if err := validator.ValidateStruct(invalidUser); err != nil {
        if verr, ok := err.(*validator.ValidationError); ok {
            fmt.Println("First error:", verr.First())
            fmt.Println("All errors:", verr.All())
            fmt.Println("Field errors:", verr.GetFieldErrors())
        }
    }
}
```

## User Registration

Complete user registration with validation:

```go
package main

import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/response"
    "github.com/budimanlai/go-pkg/validator"
    "github.com/gofiber/fiber/v2"
    "log"
)

type RegisterRequest struct {
    Name            string `json:"name" validate:"required,min=3,max=50"`
    Email           string `json:"email" validate:"required,email"`
    Password        string `json:"password" validate:"required,min=8,max=64"`
    PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
    Age             int    `json:"age" validate:"required,gte=18,lte=100"`
    Phone           string `json:"phone" validate:"omitempty,min=10,max=15"`
    Website         string `json:"website" validate:"omitempty,url"`
}

func main() {
    // Setup i18n
    i18nMgr, _ := i18n.NewI18nManager(i18n.Config{
        LocalesPath:     "./locales",
        DefaultLanguage: "en",
    })
    validator.SetI18nManager(i18nMgr)
    response.SetI18nManager(i18nMgr)
    
    app := fiber.New()
    app.Use(i18n.I18nMiddleware(i18nMgr))
    
    app.Post("/register", register)
    
    log.Fatal(app.Listen(":3000"))
}

func register(c *fiber.Ctx) error {
    var req RegisterRequest
    
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    // Validate
    if err := validator.ValidateStructWithContext(c, req); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Check if email already exists
    if emailExists(req.Email) {
        return response.BadRequestI18n(c, "email_already_registered", nil)
    }
    
    // Create user
    user := createUser(req)
    
    return response.SuccessI18n(c, "registration_success", fiber.Map{
        "id":    user.ID,
        "email": user.Email,
    })
}
```

## Complex Structs

Validation for complex business objects:

```go
type Address struct {
    Street  string `json:"street" validate:"required,min=5"`
    City    string `json:"city" validate:"required,min=2"`
    State   string `json:"state" validate:"required,len=2"`
    ZipCode string `json:"zip_code" validate:"required,numeric,len=5"`
    Country string `json:"country" validate:"required,len=2"`
}

type PaymentMethod struct {
    Type       string `json:"type" validate:"required,oneof=card paypal bank"`
    CardNumber string `json:"card_number" validate:"required_if=Type card,numeric,len=16"`
    CardCVV    string `json:"card_cvv" validate:"required_if=Type card,numeric,len=3"`
    CardExpiry string `json:"card_expiry" validate:"required_if=Type card,len=5"`
    Email      string `json:"email" validate:"required_if=Type paypal,email"`
    BankCode   string `json:"bank_code" validate:"required_if=Type bank,numeric"`
}

type Order struct {
    OrderID       string        `json:"order_id" validate:"required,uuid4"`
    CustomerEmail string        `json:"customer_email" validate:"required,email"`
    Items         []OrderItem   `json:"items" validate:"required,min=1,dive"`
    Total         float64       `json:"total" validate:"required,gt=0"`
    Address       Address       `json:"address" validate:"required"`
    Payment       PaymentMethod `json:"payment" validate:"required"`
}

type OrderItem struct {
    ProductID string  `json:"product_id" validate:"required,uuid4"`
    Quantity  int     `json:"quantity" validate:"required,gte=1,lte=100"`
    Price     float64 `json:"price" validate:"required,gt=0"`
}

func createOrder(c *fiber.Ctx) error {
    var order Order
    
    if err := c.BodyParser(&order); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    if err := validator.ValidateStructWithContext(c, order); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Validate items total matches order total
    calculatedTotal := 0.0
    for _, item := range order.Items {
        calculatedTotal += item.Price * float64(item.Quantity)
    }
    
    if calculatedTotal != order.Total {
        return response.BadRequestI18n(c, "total_mismatch", nil)
    }
    
    // Process order
    processOrder(&order)
    
    return response.SuccessI18n(c, "order_created", order)
}
```

## Password Validation

Advanced password validation with confirmation:

```go
type PasswordChange struct {
    CurrentPassword string `json:"current_password" validate:"required"`
    NewPassword     string `json:"new_password" validate:"required,min=8,max=64,containsany=!@#$%^&*"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type ResetPassword struct {
    Token           string `json:"token" validate:"required,uuid4"`
    NewPassword     string `json:"new_password" validate:"required,min=8,max=64"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

func changePassword(c *fiber.Ctx) error {
    userID := c.Locals("userID").(int)
    
    var req PasswordChange
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    if err := validator.ValidateStructWithContext(c, req); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Verify current password
    if !verifyPassword(userID, req.CurrentPassword) {
        return response.BadRequestI18n(c, "invalid_current_password", nil)
    }
    
    // Check if new password is same as current
    if req.CurrentPassword == req.NewPassword {
        return response.BadRequestI18n(c, "new_password_same_as_current", nil)
    }
    
    // Update password
    updatePassword(userID, req.NewPassword)
    
    return response.SuccessI18n(c, "password_changed", nil)
}

func resetPassword(c *fiber.Ctx) error {
    var req ResetPassword
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    if err := validator.ValidateStructWithContext(c, req); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Verify token
    userID, err := verifyResetToken(req.Token)
    if err != nil {
        return response.BadRequestI18n(c, "invalid_reset_token", nil)
    }
    
    // Update password
    updatePassword(userID, req.NewPassword)
    
    // Invalidate token
    invalidateResetToken(req.Token)
    
    return response.SuccessI18n(c, "password_reset_success", nil)
}
```

## Custom Error Messages

Handling validation errors with custom responses:

```go
type Product struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    SKU         string  `json:"sku" validate:"required,alphanum,len=10"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Stock       int     `json:"stock" validate:"required,gte=0"`
    Category    string  `json:"category" validate:"required,oneof=electronics clothing food"`
    Description string  `json:"description" validate:"required,min=10,max=500"`
}

func createProduct(c *fiber.Ctx) error {
    var product Product
    
    if err := c.BodyParser(&product); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    if err := validator.ValidateStructWithContext(c, product); err != nil {
        if verr, ok := err.(*validator.ValidationError); ok {
            // Log all errors for debugging
            log.Printf("Product validation failed: %v", verr.All())
            
            // Check for specific field errors
            fieldErrors := verr.GetFieldErrors()
            
            // Custom handling for SKU
            if errors, exists := fieldErrors["sku"]; exists {
                log.Printf("SKU validation failed: %v", errors)
                // Could check SKU uniqueness here
            }
            
            // Custom handling for price
            if errors, exists := fieldErrors["price"]; exists {
                log.Printf("Price validation failed: %v", errors)
            }
            
            return response.ValidationErrorI18n(c, err)
        }
        return response.Error(c, 500, "Unexpected error")
    }
    
    // Check SKU uniqueness
    if skuExists(product.SKU) {
        return response.BadRequestI18n(c, "sku_already_exists", fiber.Map{
            "SKU": product.SKU,
        })
    }
    
    // Create product
    saveProduct(&product)
    
    return response.SuccessI18n(c, "product_created", product)
}
```

## Nested Structures

Validation for nested and complex structures:

```go
type Company struct {
    Name    string   `json:"name" validate:"required,min=3"`
    Email   string   `json:"email" validate:"required,email"`
    Address Address  `json:"address" validate:"required"`
    Contact Contact  `json:"contact" validate:"required"`
    Employees []Employee `json:"employees" validate:"omitempty,dive"`
}

type Contact struct {
    Name  string `json:"name" validate:"required,min=3"`
    Phone string `json:"phone" validate:"required,min=10,max=15"`
    Email string `json:"email" validate:"required,email"`
}

type Employee struct {
    Name       string `json:"name" validate:"required,min=3"`
    Email      string `json:"email" validate:"required,email"`
    Position   string `json:"position" validate:"required"`
    Department string `json:"department" validate:"required,oneof=IT HR Sales Marketing"`
    Salary     float64 `json:"salary" validate:"required,gt=0"`
}

func createCompany(c *fiber.Ctx) error {
    var company Company
    
    if err := c.BodyParser(&company); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    // Validates company and all nested structures
    if err := validator.ValidateStructWithContext(c, company); err != nil {
        if verr, ok := err.(*validator.ValidationError); ok {
            // Show which nested field failed
            for field, errs := range verr.GetFieldErrors() {
                log.Printf("Field %s errors: %v", field, errs)
            }
        }
        return response.ValidationErrorI18n(c, err)
    }
    
    // Additional business logic validation
    if len(company.Employees) > 1000 {
        return response.BadRequestI18n(c, "too_many_employees", fiber.Map{
            "Max":   "1000",
            "Count": len(company.Employees),
        })
    }
    
    saveCompany(&company)
    
    return response.SuccessI18n(c, "company_created", company)
}
```

## Conditional Validation

Validation that depends on other field values:

```go
type ShippingInfo struct {
    Method          string  `json:"method" validate:"required,oneof=standard express overnight pickup"`
    Address         *Address `json:"address" validate:"required_unless=Method pickup"`
    PickupLocation  string  `json:"pickup_location" validate:"required_if=Method pickup"`
    Insurance       bool    `json:"insurance"`
    InsuranceAmount float64 `json:"insurance_amount" validate:"required_if=Insurance true,gt=0"`
}

type ContactInfo struct {
    Email string `json:"email" validate:"required_without=Phone,omitempty,email"`
    Phone string `json:"phone" validate:"required_without=Email,omitempty,min=10"`
}

type Subscription struct {
    Plan        string `json:"plan" validate:"required,oneof=free basic premium enterprise"`
    BillingCycle string `json:"billing_cycle" validate:"required_unless=Plan free,oneof=monthly yearly"`
    PaymentMethod string `json:"payment_method" validate:"required_unless=Plan free,oneof=card paypal"`
    PromoCode   string `json:"promo_code" validate:"omitempty,alphanum,len=8"`
}

func processShipping(c *fiber.Ctx) error {
    var shipping ShippingInfo
    
    if err := c.BodyParser(&shipping); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    if err := validator.ValidateStructWithContext(c, shipping); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Calculate shipping cost
    cost := calculateShippingCost(shipping)
    
    return response.SuccessI18n(c, "shipping_calculated", fiber.Map{
        "cost": cost,
        "method": shipping.Method,
    })
}

func updateContact(c *fiber.Ctx) error {
    var contact ContactInfo
    
    if err := c.BodyParser(&contact); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    // At least one of Email or Phone must be provided
    if err := validator.ValidateStructWithContext(c, contact); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    saveContact(&contact)
    
    return response.SuccessI18n(c, "contact_updated", contact)
}

func subscribe(c *fiber.Ctx) error {
    var sub Subscription
    
    if err := c.BodyParser(&sub); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    if err := validator.ValidateStructWithContext(c, sub); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Validate promo code if provided
    if sub.PromoCode != "" {
        if !isValidPromoCode(sub.PromoCode) {
            return response.BadRequestI18n(c, "invalid_promo_code", nil)
        }
    }
    
    createSubscription(&sub)
    
    return response.SuccessI18n(c, "subscription_created", sub)
}
```

## File Upload Validation

Validating file uploads with metadata:

```go
type FileUpload struct {
    Filename    string `json:"filename" validate:"required,min=1,max=255"`
    ContentType string `json:"content_type" validate:"required,oneof=image/jpeg image/png image/gif application/pdf"`
    Size        int64  `json:"size" validate:"required,gt=0,lte=10485760"` // Max 10MB
    Description string `json:"description" validate:"omitempty,max=500"`
}

type BulkUpload struct {
    Files []FileUpload `json:"files" validate:"required,min=1,max=10,dive"`
}

func uploadFile(c *fiber.Ctx) error {
    // Get file
    file, err := c.FormFile("file")
    if err != nil {
        return response.BadRequestI18n(c, "file_required", nil)
    }
    
    // Create validation struct
    upload := FileUpload{
        Filename:    file.Filename,
        ContentType: file.Header.Get("Content-Type"),
        Size:        file.Size,
        Description: c.FormValue("description"),
    }
    
    // Validate
    if err := validator.ValidateStructWithContext(c, upload); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Save file
    path := saveFile(file)
    
    return response.SuccessI18n(c, "file_uploaded", fiber.Map{
        "filename": upload.Filename,
        "path":     path,
        "size":     upload.Size,
    })
}

func bulkUpload(c *fiber.Ctx) error {
    var bulk BulkUpload
    
    if err := c.BodyParser(&bulk); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    // Validate all files
    if err := validator.ValidateStructWithContext(c, bulk); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Process uploads
    results := processUploads(bulk.Files)
    
    return response.SuccessI18n(c, "files_uploaded", fiber.Map{
        "count":   len(results),
        "results": results,
    })
}
```

## Testing Validation

Unit tests for validation:

```go
package main

import (
    "testing"
    "github.com/budimanlai/go-pkg/validator"
)

func TestUserValidation(t *testing.T) {
    tests := []struct {
        name      string
        user      User
        wantError bool
        errorCount int
    }{
        {
            name: "Valid user",
            user: User{
                Email:    "john@example.com",
                Password: "secure123",
                Age:      25,
            },
            wantError: false,
        },
        {
            name: "Invalid email",
            user: User{
                Email:    "invalid-email",
                Password: "secure123",
                Age:      25,
            },
            wantError:  true,
            errorCount: 1,
        },
        {
            name: "All fields invalid",
            user: User{
                Email:    "invalid",
                Password: "123",
                Age:      15,
            },
            wantError:  true,
            errorCount: 3,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validator.ValidateStruct(tt.user)
            
            if tt.wantError {
                if err == nil {
                    t.Error("Expected error, got nil")
                    return
                }
                
                verr, ok := err.(*validator.ValidationError)
                if !ok {
                    t.Error("Expected ValidationError type")
                    return
                }
                
                if len(verr.All()) != tt.errorCount {
                    t.Errorf("Expected %d errors, got %d", tt.errorCount, len(verr.All()))
                }
            } else {
                if err != nil {
                    t.Errorf("Expected no error, got %v", err)
                }
            }
        })
    }
}
```

## Related Documentation

- [README](README.md) - Package overview
- [Validation Tags](validation-tags.md) - Complete validation rules
- [Error Handling](error-handling.md) - ValidationError type
- [I18n Integration](i18n-integration.md) - Multilingual validation
