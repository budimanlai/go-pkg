package validator

import (
	"testing"
)

type TestUser struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=130"`
}

func TestValidateStruct_Valid(t *testing.T) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	err := ValidateStruct(user, "id")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateStruct_Invalid_ID(t *testing.T) {
	user := &TestUser{
		Name:  "",
		Email: "invalid-email",
		Age:   -5,
	}

	err := ValidateStruct(user, "id")
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
}

func TestValidateStruct_Invalid_EN(t *testing.T) {
	user := &TestUser{
		Name:  "",
		Email: "invalid-email",
		Age:   -5,
	}

	err := ValidateStruct(user, "en")
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

func TestValidateStruct_UnknownLang(t *testing.T) {
	user := &TestUser{
		Name: "",
	}

	err := ValidateStruct(user, "unknown")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	// Should default to "id"
	expectedFirst := "Name wajib diisi"
	if valErr.First() != expectedFirst {
		t.Errorf("Expected first error '%s', got '%s'", expectedFirst, valErr.First())
	}
}

func TestAddLanguage(t *testing.T) {
	// Add new language
	jvMessages := map[string]string{
		"required": "%s kudu diisi",
		"default":  "%s ora valid",
	}
	AddLanguage("jv", jvMessages)

	// Check if added
	langs := GetLanguages()
	found := false
	for _, lang := range langs {
		if lang == "jv" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Language 'jv' not added")
	}

	// Test validation with new language
	user := &TestUser{Name: ""}
	err := ValidateStruct(user, "jv")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	expectedFirst := "Name kudu diisi"
	if valErr.First() != expectedFirst {
		t.Errorf("Expected first error '%s', got '%s'", expectedFirst, valErr.First())
	}
}

func TestUpdateLanguage(t *testing.T) {
	// Update existing language
	customID := map[string]string{
		"required": "%s harus ada",
		"default":  "%s salah",
	}
	UpdateLanguage("id", customID)

	// Test validation
	user := &TestUser{Name: ""}
	err := ValidateStruct(user, "id")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	expectedFirst := "Name harus ada"
	if valErr.First() != expectedFirst {
		t.Errorf("Expected first error '%s', got '%s'", expectedFirst, valErr.First())
	}
}

func TestGetLanguages(t *testing.T) {
	langs := GetLanguages()
	if len(langs) < 2 {
		t.Errorf("Expected at least 2 languages, got %d", len(langs))
	}

	// Check if "id" and "en" exist
	hasID := false
	hasEN := false
	for _, lang := range langs {
		if lang == "id" {
			hasID = true
		}
		if lang == "en" {
			hasEN = true
		}
	}
	if !hasID {
		t.Error("Language 'id' not found")
	}
	if !hasEN {
		t.Error("Language 'en' not found")
	}
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
