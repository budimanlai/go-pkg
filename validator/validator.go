package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// Validator is the global validator instance from go-playground/validator.
	// It is initialized automatically and can be used throughout the application.
	Validator *validator.Validate

	// Messages contains validation error messages organized by language and validation tag.
	// Structure: map[language]map[tag]messageTemplate
	//
	// Supported languages: "id" (Indonesian), "en" (English)
	// Message templates use %s placeholders for field name and parameter values.
	//
	// Example:
	//   Messages["id"]["required"] = "%s wajib diisi"
	//   Messages["en"]["required"] = "%s is required"
	Messages = map[string]map[string]string{
		"id": {
			"required": "%s wajib diisi",
			"email":    "%s harus berupa alamat email yang valid",
			"min":      "%s minimal %s karakter",
			"max":      "%s maksimal %s karakter",
			"gte":      "%s harus lebih besar atau sama dengan %s",
			"lte":      "%s harus lebih kecil atau sama dengan %s",
			"len":      "%s harus memiliki panjang %s",
			"numeric":  "%s harus berupa angka",
			"alphanum": "%s hanya boleh berisi huruf dan angka",
			"default":  "%s tidak valid (%s)",
		},
		"en": {
			"required": "%s is required",
			"email":    "%s must be a valid email address",
			"min":      "%s must be at least %s characters",
			"max":      "%s must be at most %s characters",
			"gte":      "%s must be greater than or equal to %s",
			"lte":      "%s must be less than or equal to %s",
			"len":      "%s must be exactly %s characters",
			"numeric":  "%s must be numeric",
			"alphanum": "%s must contain only letters and numbers",
			"default":  "%s is invalid (%s)",
		},
	}
)

func init() {
	Validator = validator.New()
}

// ValidationError is a custom error type for validation failures.
// It collects multiple validation error messages and provides convenient methods
// to access them individually or as a group.
//
// Fields:
//   - Messages: Slice of validation error messages
//
// Methods:
//   - Error(): Returns all messages joined by semicolon (implements error interface)
//   - First(): Returns the first error message
//   - All(): Returns all error messages as a slice
//
// Example:
//
//	if err := ValidateStruct(user, "en"); err != nil {
//	    if verr, ok := err.(*ValidationError); ok {
//	        fmt.Println(verr.First())  // Get first error
//	        fmt.Println(verr.All())    // Get all errors
//	    }
//	}
type ValidationError struct {
	Messages []string
}

// Error implements the error interface for ValidationError.
// It returns all validation error messages joined by semicolons.
//
// Returns:
//   - string: All error messages concatenated with "; " separator
//
// Example:
//
//	verr := &ValidationError{Messages: []string{"Email is required", "Password is too short"}}
//	fmt.Println(verr.Error())
//	// Output: Email is required; Password is too short
func (ve *ValidationError) Error() string {
	return strings.Join(ve.Messages, "; ")
}

// First returns the first validation error message.
// If there are no error messages, it returns an empty string.
//
// Returns:
//   - string: First error message, or empty string if no messages exist
//
// Example:
//
//	verr := &ValidationError{Messages: []string{"Email is required", "Password is too short"}}
//	fmt.Println(verr.First())
//	// Output: Email is required
func (ve *ValidationError) First() string {
	if len(ve.Messages) > 0 {
		return ve.Messages[0]
	}
	return ""
}

// All returns all validation error messages as a slice.
//
// Returns:
//   - []string: Slice containing all error messages
//
// Example:
//
//	verr := &ValidationError{Messages: []string{"Email is required", "Password is too short"}}
//	for _, msg := range verr.All() {
//	    fmt.Println(msg)
//	}
//	// Output:
//	// Email is required
//	// Password is too short
func (ve *ValidationError) All() []string {
	return ve.Messages
}

// ValidateStruct validates a struct using validation tags and returns localized error messages.
// It uses the global Validator instance and translates errors based on the specified language.
//
// The function automatically converts field names to title case for better readability.
// If validation succeeds, it returns nil. If validation fails, it returns a ValidationError
// containing all validation error messages in the specified language.
//
// Parameters:
//   - s: The struct to validate (must have validation tags)
//   - lang: Language code for error messages ("en" or "id")
//
// Returns:
//   - error: nil if validation succeeds, *ValidationError if validation fails
//
// Example:
//
//	type User struct {
//	    Email    string `validate:"required,email"`
//	    Password string `validate:"required,min=8"`
//	    Age      int    `validate:"gte=18"`
//	}
//
//	user := User{Email: "invalid", Password: "123", Age: 15}
//	if err := ValidateStruct(user, "en"); err != nil {
//	    if verr, ok := err.(*ValidationError); ok {
//	        fmt.Println(verr.First())
//	        // Output: Email must be a valid email address
//	    }
//	}
func ValidateStruct(s interface{}, lang string) error {
	err := Validator.Struct(s)
	if err == nil {
		return nil
	}

	var messages []string
	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		for _, e := range validateErrs {
			message := getUserFriendlyMessage(e.Field(), e.Tag(), e.Param(), lang)
			messages = append(messages, message)
		}
	} else {
		// Jika bukan validation error, kembalikan error asli
		messages = append(messages, err.Error())
	}
	return &ValidationError{Messages: messages}
}

// getUserFriendlyMessage generates a user-friendly error message based on validation failure details.
// It formats the message using the appropriate language template with field name and parameter values.
//
// The function:
//   - Converts field names to title case for readability
//   - Selects the appropriate message template based on language and validation tag
//   - Falls back to Indonesian ("id") if the specified language is not found
//   - Uses "default" template if the validation tag is not found
//
// Parameters:
//   - field: Name of the field that failed validation
//   - tag: Validation tag that failed (e.g., "required", "email", "min")
//   - param: Parameter value for the validation tag (e.g., "8" for min=8)
//   - lang: Language code for the error message
//
// Returns:
//   - string: Formatted user-friendly error message
//
// Example:
//
//	msg := getUserFriendlyMessage("email", "required", "", "en")
//	// Returns: "Email is required"
//
//	msg := getUserFriendlyMessage("password", "min", "8", "en")
//	// Returns: "Password must be at least 8 characters"
func getUserFriendlyMessage(field, tag, param, lang string) string {
	// Gunakan unicode-aware title caser
	caser := cases.Title(language.Und)
	fieldName := caser.String(field)

	// Ambil pesan berdasarkan bahasa, default ke "id" jika tidak ada
	langMessages, exists := Messages[lang]
	if !exists {
		langMessages = Messages["id"]
	}

	template, exists := langMessages[tag]
	if !exists {
		template = langMessages["default"]
	}

	if param == "" {
		return fmt.Sprintf(template, fieldName)
	}
	return fmt.Sprintf(template, fieldName, param)
}

// AddLanguage adds a new language with custom validation messages to the Messages map.
// If the language already exists, this function does nothing (messages are not overwritten).
// To update an existing language, use UpdateLanguage instead.
//
// Parameters:
//   - lang: Language code (e.g., "fr", "es", "zh")
//   - messages: Map of validation tags to message templates
//
// Example:
//
//	AddLanguage("fr", map[string]string{
//	    "required": "%s est requis",
//	    "email":    "%s doit être une adresse email valide",
//	    "min":      "%s doit comporter au moins %s caractères",
//	    "default":  "%s n'est pas valide (%s)",
//	})
func AddLanguage(lang string, messages map[string]string) {
	if _, exists := Messages[lang]; !exists {
		Messages[lang] = messages
	}
}

// UpdateLanguage updates or adds validation messages for a language.
// If the language already exists, it completely replaces the existing messages.
// If the language doesn't exist, it creates a new entry.
//
// Parameters:
//   - lang: Language code (e.g., "en", "id", "zh")
//   - messages: Map of validation tags to message templates
//
// Example:
//
//	// Update existing language
//	UpdateLanguage("en", map[string]string{
//	    "required": "%s is mandatory",
//	    "email":    "%s must be valid email",
//	    "default":  "%s is not valid (%s)",
//	})
//
//	// Add new language
//	UpdateLanguage("es", map[string]string{
//	    "required": "%s es requerido",
//	    "email":    "%s debe ser un correo electrónico válido",
//	    "default":  "%s no es válido (%s)",
//	})
func UpdateLanguage(lang string, messages map[string]string) {
	Messages[lang] = messages
}

// GetLanguages returns a list of all available language codes in the Messages map.
// The order of languages in the returned slice is not guaranteed.
//
// Returns:
//   - []string: Slice of language codes
//
// Example:
//
//	langs := GetLanguages()
//	fmt.Println(langs)
//	// Output (order may vary): [id en]
//
//	for _, lang := range GetLanguages() {
//	    fmt.Printf("Language: %s\n", lang)
//	}
func GetLanguages() []string {
	var langs []string
	for lang := range Messages {
		langs = append(langs, lang)
	}
	return langs
}
