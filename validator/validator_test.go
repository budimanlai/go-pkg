package validator

import (
	"net/http/httptest"
	"testing"

	"github.com/budimanlai/go-pkg/i18n"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

// Test Structs
type TestUser struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=130"`
}

type TestUserWithJSON struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"gte=18,lte=130"`
}

type TestProduct struct {
	Name        string `json:"product_name" validate:"required,min=3,max=100"`
	Price       int    `json:"price" validate:"required,gte=0"`
	SKU         string `json:"sku" validate:"required,len=8,alphanum"`
	Description string `json:"description" validate:"max=500"`
}

type TestAddress struct {
	Street  string `validate:"required"`
	City    string `validate:"required,min=2"`
	ZipCode string `validate:"required,numeric,len=5"`
}

// Setup Helper
func setupI18n() *i18n.I18nManager {
	i18nConfig := i18n.I18nConfig{
		DefaultLanguage: language.English,
		SupportedLangs:  []string{"en", "id", "zh"},
		LocalesPath:     "../locales",
	}
	manager, _ := i18n.NewI18nManager(i18nConfig)
	SetI18nManager(manager)
	return manager
}

// Helper function
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// ValidateStruct Valid Tests
func TestValidateStruct_Valid(t *testing.T) {
	t.Run("valid_struct", func(t *testing.T) {
		setupI18n()
		user := &TestUser{Name: "John Doe", Email: "john@example.com", Age: 25}
		if err := ValidateStruct(user); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("valid_with_json_tags", func(t *testing.T) {
		setupI18n()
		user := &TestUserWithJSON{Name: "Jane", Email: "jane@example.com", Age: 30}
		if err := ValidateStruct(user); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("valid_product", func(t *testing.T) {
		setupI18n()
		product := &TestProduct{Name: "Product", Price: 1000, SKU: "ABC12345", Description: "Test"}
		if err := ValidateStruct(product); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

// ValidateStructWithLang Indonesian Tests
func TestValidateStructWithLang_Invalid_ID(t *testing.T) {
	t.Run("multiple_errors_indonesian", func(t *testing.T) {
		setupI18n()
		user := &TestUser{Name: "", Email: "invalid-email", Age: -5}
		err := ValidateStructWithLang(user, "id")
		if err == nil {
			t.Error("Expected error, got nil")
		}
		valErr := err.(*ValidationError)
		if valErr.First() != "Name wajib diisi" {
			t.Errorf("Expected 'Name wajib diisi', got '%s'", valErr.First())
		}
		if len(valErr.All()) != 3 {
			t.Errorf("Expected 3 errors, got %d", len(valErr.All()))
		}
	})
}

// ValidateStructWithLang English Tests
func TestValidateStructWithLang_Invalid_EN(t *testing.T) {
	t.Run("multiple_errors_english", func(t *testing.T) {
		setupI18n()
		user := &TestUser{Name: "", Email: "invalid", Age: -5}
		err := ValidateStructWithLang(user, "en")
		valErr := err.(*ValidationError)
		if valErr.First() != "Name is required" {
			t.Errorf("Expected 'Name is required', got '%s'", valErr.First())
		}
	})

	t.Run("email_validation", func(t *testing.T) {
		setupI18n()
		user := &TestUser{Name: "John", Email: "not-email", Age: 25}
		err := ValidateStructWithLang(user, "en")
		valErr := err.(*ValidationError)
		if !contains(valErr.All(), "Email must be a valid email address") {
			t.Errorf("Expected email error, got: %v", valErr.All())
		}
	})

	t.Run("range_validation", func(t *testing.T) {
		setupI18n()
		user := &TestUser{Name: "John", Email: "john@example.com", Age: 150}
		err := ValidateStructWithLang(user, "en")
		valErr := err.(*ValidationError)
		if !contains(valErr.All(), "Age must be less than or equal to 130") {
			t.Errorf("Expected age error, got: %v", valErr.All())
		}
	})
}

// Chinese Translation Tests
func TestValidateStructWithLang_Chinese(t *testing.T) {
	t.Run("chinese_translation", func(t *testing.T) {
		setupI18n()
		user := &TestUser{Name: "", Email: "invalid"}
		err := ValidateStructWithLang(user, "zh")
		valErr := err.(*ValidationError)
		if valErr.First() == "" {
			t.Error("Expected error message")
		}
	})
}

// Default Language Tests
func TestValidateStruct_DefaultLanguage(t *testing.T) {
	t.Run("uses_default_english", func(t *testing.T) {
		setupI18n()
		user := &TestUser{Name: ""}
		err := ValidateStruct(user)
		valErr := err.(*ValidationError)
		if valErr.First() != "Name is required" {
			t.Errorf("Expected English default, got '%s'", valErr.First())
		}
	})

	t.Run("default_language_from_manager", func(t *testing.T) {
		cfg := i18n.I18nConfig{DefaultLanguage: language.Indonesian, SupportedLangs: []string{"id"}, LocalesPath: "../locales"}
		mgr, _ := i18n.NewI18nManager(cfg)
		SetI18nManager(mgr)
		user := &TestUser{Name: ""}
		err := ValidateStruct(user)
		valErr := err.(*ValidationError)
		if valErr.First() != "Name wajib diisi" {
			t.Errorf("Expected Indonesian, got '%s'", valErr.First())
		}
		setupI18n()
	})
}

// Without I18n Tests
func TestValidateStruct_WithoutI18n(t *testing.T) {
	t.Run("fallback_to_defaults", func(t *testing.T) {
		SetI18nManager(nil)
		user := &TestUser{Name: ""}
		err := ValidateStruct(user)
		valErr := err.(*ValidationError)
		if valErr.First() != "Name is required" {
			t.Errorf("Expected default message, got '%s'", valErr.First())
		}
		setupI18n()
	})
}

// ValidateStructWithContext Tests
func TestValidateStructWithContext(t *testing.T) {
	t.Run("context_with_english", func(t *testing.T) {
		setupI18n()
		app := fiber.New()
		app.Post("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "en")
			user := &TestUser{Name: ""}
			err := ValidateStructWithContext(c, user)
			valErr := err.(*ValidationError)
			if valErr.First() != "Name is required" {
				t.Errorf("Expected English, got '%s'", valErr.First())
			}
			return c.SendString("OK")
		})
		req := httptest.NewRequest("POST", "/test", nil)
		app.Test(req)
	})

	t.Run("context_with_indonesian", func(t *testing.T) {
		setupI18n()
		app := fiber.New()
		app.Post("/test", func(c *fiber.Ctx) error {
			c.Locals("language", "id")
			user := &TestUser{Name: ""}
			err := ValidateStructWithContext(c, user)
			valErr := err.(*ValidationError)
			if valErr.First() != "Name wajib diisi" {
				t.Errorf("Expected Indonesian, got '%s'", valErr.First())
			}
			return c.SendString("OK")
		})
		req := httptest.NewRequest("POST", "/test", nil)
		app.Test(req)
	})

	t.Run("context_without_language", func(t *testing.T) {
		setupI18n()
		app := fiber.New()
		app.Post("/test", func(c *fiber.Ctx) error {
			user := &TestUser{Name: ""}
			err := ValidateStructWithContext(c, user)
			valErr := err.(*ValidationError)
			if valErr.First() != "Name is required" {
				t.Errorf("Expected default English, got '%s'", valErr.First())
			}
			return c.SendString("OK")
		})
		req := httptest.NewRequest("POST", "/test", nil)
		app.Test(req)
	})
}

// ValidationError Methods Tests
func TestValidationError_Methods(t *testing.T) {
	t.Run("error_method", func(t *testing.T) {
		valErr := &ValidationError{Messages: []string{"Error 1", "Error 2"}}
		if valErr.Error() != "Error 1; Error 2" {
			t.Errorf("Expected 'Error 1; Error 2', got '%s'", valErr.Error())
		}
	})

	t.Run("first_method", func(t *testing.T) {
		valErr := &ValidationError{Messages: []string{"First", "Second"}}
		if valErr.First() != "First" {
			t.Errorf("Expected 'First', got '%s'", valErr.First())
		}
	})

	t.Run("first_empty", func(t *testing.T) {
		valErr := &ValidationError{Messages: []string{}}
		if valErr.First() != "" {
			t.Error("Expected empty string")
		}
	})

	t.Run("all_method", func(t *testing.T) {
		valErr := &ValidationError{Messages: []string{"E1", "E2"}}
		if len(valErr.All()) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(valErr.All()))
		}
	})

	t.Run("get_field_errors", func(t *testing.T) {
		valErr := &ValidationError{
			Errors: map[string][]string{"email": {"Required", "Invalid"}, "name": {"Required"}},
		}
		errors := valErr.GetFieldErrors()
		if len(errors) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(errors))
		}
		if len(errors["email"]) != 2 {
			t.Errorf("Expected 2 email errors, got %d", len(errors["email"]))
		}
	})
}

// JSON Tag Tests
func TestValidateStruct_WithJSONTag(t *testing.T) {
	t.Run("json_tag_field_names", func(t *testing.T) {
		setupI18n()
		user := &TestUserWithJSON{Name: "", Email: "invalid", Age: 15}
		err := ValidateStructWithLang(user, "en")
		valErr := err.(*ValidationError)
		fieldErrors := valErr.GetFieldErrors()
		
		if _, exists := fieldErrors["name"]; !exists {
			t.Error("Expected 'name' from json tag")
		}
		if _, exists := fieldErrors["email"]; !exists {
			t.Error("Expected 'email' from json tag")
		}
		if _, exists := fieldErrors["age"]; !exists {
			t.Error("Expected 'age' from json tag")
		}
	})

	t.Run("complex_json_tags", func(t *testing.T) {
		setupI18n()
		product := &TestProduct{Name: "", Price: -1, SKU: "123"}
		err := ValidateStructWithLang(product, "en")
		valErr := err.(*ValidationError)
		fieldErrors := valErr.GetFieldErrors()
		
		if _, exists := fieldErrors["product_name"]; !exists {
			t.Error("Expected 'product_name' from json tag")
		}
	})
}

// Edge Cases Tests
func TestValidation_EdgeCases(t *testing.T) {
	t.Run("nil_struct", func(t *testing.T) {
		setupI18n()
		var user *TestUser
		err := ValidateStruct(user)
		if err == nil {
			t.Error("Expected error for nil struct")
		}
	})

	t.Run("empty_struct", func(t *testing.T) {
		setupI18n()
		type Empty struct{}
		if err := ValidateStruct(&Empty{}); err != nil {
			t.Errorf("Expected no error for empty struct, got %v", err)
		}
	})

	t.Run("struct_without_tags", func(t *testing.T) {
		setupI18n()
		type NoTags struct{ Name string }
		if err := ValidateStruct(&NoTags{Name: ""}); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("boundary_values", func(t *testing.T) {
		setupI18n()
		user := &TestUser{Name: "John", Email: "john@test.com", Age: 130}
		if err := ValidateStruct(user); err != nil {
			t.Errorf("Expected no error at boundary, got %v", err)
		}
		user.Age = 0
		if err := ValidateStruct(user); err != nil {
			t.Errorf("Expected no error at lower boundary, got %v", err)
		}
	})
}

// Validation Rules Tests
func TestValidationRules(t *testing.T) {
	t.Run("numeric_validation", func(t *testing.T) {
		setupI18n()
		addr := &TestAddress{Street: "Main", City: "NY", ZipCode: "abc"}
		if err := ValidateStruct(addr); err == nil {
			t.Error("Expected numeric validation error")
		}
	})

	t.Run("len_validation", func(t *testing.T) {
		setupI18n()
		addr := &TestAddress{Street: "Main", City: "NY", ZipCode: "123"}
		if err := ValidateStruct(addr); err == nil {
			t.Error("Expected length validation error")
		}
	})

	t.Run("alphanum_validation", func(t *testing.T) {
		setupI18n()
		product := &TestProduct{Name: "Product", Price: 100, SKU: "ABC-123"}
		if err := ValidateStruct(product); err == nil {
			t.Error("Expected alphanum validation error")
		}
	})

	t.Run("min_length", func(t *testing.T) {
		setupI18n()
		addr := &TestAddress{Street: "Main", City: "A", ZipCode: "12345"}
		err := ValidateStruct(addr)
		valErr := err.(*ValidationError)
		if !contains(valErr.All(), "City must be at least 2 characters") {
			t.Errorf("Expected min length error, got: %v", valErr.All())
		}
	})

	t.Run("max_length", func(t *testing.T) {
		setupI18n()
		product := &TestProduct{Name: "Valid", Price: 100, SKU: "ABC12345", Description: string(make([]byte, 600))}
		if err := ValidateStruct(product); err == nil {
			t.Error("Expected max length error")
		}
	})
}

// SetI18nManager Tests
func TestSetI18nManager(t *testing.T) {
	t.Run("set_manager", func(t *testing.T) {
		manager := setupI18n()
		if i18nManager == nil {
			t.Error("Expected i18nManager to be set")
		}
		if i18nManager != manager {
			t.Error("Expected managers to match")
		}
	})

	t.Run("nil_manager", func(t *testing.T) {
		SetI18nManager(nil)
		if i18nManager != nil {
			t.Error("Expected nil manager")
		}
		user := &TestUser{Name: ""}
		if err := ValidateStruct(user); err == nil {
			t.Error("Expected error even without manager")
		}
		setupI18n()
	})
}
