package validator

import (
	"testing"

	"github.com/budimanlai/go-pkg/i18n"
	"golang.org/x/text/language"
)

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

// Setup i18n for tests
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

func TestValidateStruct_Valid(t *testing.T) {
	setupI18n()

	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	err := ValidateStruct(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateStructWithLang_Invalid_ID(t *testing.T) {
	setupI18n()

	user := &TestUser{
		Name:  "",
		Email: "invalid-email",
		Age:   -5,
	}

	err := ValidateStructWithLang(user, "id")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	expectedFirst := "Name wajib diisi"
	if valErr.First() != expectedFirst {
		t.Errorf("Expected first error '%s', got '%s'", expectedFirst, valErr.First())
	}

	all := valErr.All()
	if len(all) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(all))
	}

	expectedMessages := []string{
		"Name wajib diisi",
		"Email harus berupa alamat email yang valid",
		"Age harus lebih besar atau sama dengan 0",
	}
	for i, msg := range expectedMessages {
		if all[i] != msg {
			t.Errorf("Expected error %d '%s', got '%s'", i, msg, all[i])
		}
	}

	// Test field errors
	fieldErrors := valErr.GetFieldErrors()
	if len(fieldErrors) != 3 {
		t.Errorf("Expected 3 field errors, got %d", len(fieldErrors))
	}

	// Check Name field errors
	if nameErrs, ok := fieldErrors["Name"]; ok {
		if len(nameErrs) != 1 {
			t.Errorf("Expected 1 error for Name, got %d", len(nameErrs))
		}
		if nameErrs[0] != "Name wajib diisi" {
			t.Errorf("Expected Name error 'Name wajib diisi', got '%s'", nameErrs[0])
		}
	} else {
		t.Error("Expected Name in field errors")
	}

	// Check Email field errors
	if emailErrs, ok := fieldErrors["Email"]; ok {
		if len(emailErrs) != 1 {
			t.Errorf("Expected 1 error for Email, got %d", len(emailErrs))
		}
	} else {
		t.Error("Expected Email in field errors")
	}

	// Check Age field errors
	if ageErrs, ok := fieldErrors["Age"]; ok {
		if len(ageErrs) != 1 {
			t.Errorf("Expected 1 error for Age, got %d", len(ageErrs))
		}
	} else {
		t.Error("Expected Age in field errors")
	}
}

func TestValidateStructWithLang_Invalid_EN(t *testing.T) {
	setupI18n()

	user := &TestUser{
		Name:  "",
		Email: "invalid-email",
		Age:   -5,
	}

	err := ValidateStructWithLang(user, "en")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	expectedFirst := "Name is required"
	if valErr.First() != expectedFirst {
		t.Errorf("Expected first error '%s', got '%s'", expectedFirst, valErr.First())
	}
}

func TestValidateStruct_DefaultLanguage(t *testing.T) {
	setupI18n()

	user := &TestUser{
		Name: "",
	}

	err := ValidateStruct(user)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	// Should use default language (English)
	expectedFirst := "Name is required"
	if valErr.First() != expectedFirst {
		t.Errorf("Expected first error '%s', got '%s'", expectedFirst, valErr.First())
	}
}

func TestValidateStruct_WithoutI18n(t *testing.T) {
	// Reset i18nManager to nil
	SetI18nManager(nil)

	user := &TestUser{
		Name: "",
	}

	err := ValidateStruct(user)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	// Should use DefaultMessages (English)
	expectedFirst := "Name is required"
	if valErr.First() != expectedFirst {
		t.Errorf("Expected first error '%s', got '%s'", expectedFirst, valErr.First())
	}
}

func TestValidateStructWithLang_Chinese(t *testing.T) {
	setupI18n()

	user := &TestUser{
		Name:  "",
		Email: "invalid-email",
	}

	err := ValidateStructWithLang(user, "zh")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	// Check Chinese translation exists
	first := valErr.First()
	if first == "" {
		t.Error("Expected error message, got empty string")
	}
}

func TestValidateStructWithContext(t *testing.T) {
	setupI18n()

	t.Skip("Skipping Fiber context test - requires full HTTP request setup")

	// Note: Testing ValidateStructWithContext requires a full Fiber app
	// with HTTP request. For integration tests, see examples/validator_with_fiber.go
}

func TestValidationError_Methods(t *testing.T) {
	messages := []string{"Error 1", "Error 2"}
	valErr := &ValidationError{Messages: messages}

	// Test Error()
	expectedError := "Error 1; Error 2"
	if valErr.Error() != expectedError {
		t.Errorf("Expected error string '%s', got '%s'", expectedError, valErr.Error())
	}

	// Test First()
	if valErr.First() != "Error 1" {
		t.Errorf("Expected first 'Error 1', got '%s'", valErr.First())
	}

	// Test All()
	all := valErr.All()
	if len(all) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(all))
	}
	if all[0] != "Error 1" || all[1] != "Error 2" {
		t.Errorf("All messages incorrect: %v", all)
	}
}

func TestValidateStruct_WithJSONTag(t *testing.T) {
	setupI18n()

	user := &TestUserWithJSON{
		Name:  "",
		Email: "invalid-email",
		Age:   15,
	}

	err := ValidateStructWithLang(user, "en")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	// Check field errors use json tag names
	fieldErrors := valErr.GetFieldErrors()

	// Check that field names match json tags (lowercase)
	if _, exists := fieldErrors["name"]; !exists {
		t.Error("Expected field 'name' (from json tag), not found in errors")
	}

	if _, exists := fieldErrors["email"]; !exists {
		t.Error("Expected field 'email' (from json tag), not found in errors")
	}

	if _, exists := fieldErrors["age"]; !exists {
		t.Error("Expected field 'age' (from json tag), not found in errors")
	}

	// Check that messages use json tag names
	if !contains(valErr.All(), "name is required") {
		t.Errorf("Expected message with 'name' (lowercase from json tag), got: %v", valErr.All())
	}

	if !contains(valErr.All(), "email must be a valid email address") {
		t.Errorf("Expected message with 'email' (lowercase from json tag), got: %v", valErr.All())
	}
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
