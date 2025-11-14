package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/budimanlai/go-pkg/i18n"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// Validator is the global validator instance from go-playground/validator.
	// It is initialized automatically and can be used throughout the application.
	Validator *validator.Validate

	// i18nManager holds the I18nManager instance for translation support.
	// If set, validator will use i18n for error messages. If nil, it falls back to default English messages.
	i18nManager *i18n.I18nManager

	// DefaultMessages contains fallback validation error messages in English.
	// These are used when i18nManager is not set or translation is not found.
	// Message templates use placeholders: {{.FieldName}}, {{.Param}}, {{.Tag}}
	DefaultMessages = map[string]string{
		"required": "{{.FieldName}} is required",
		"email":    "{{.FieldName}} must be a valid email address",
		"min":      "{{.FieldName}} must be at least {{.Param}} characters",
		"max":      "{{.FieldName}} must be at most {{.Param}} characters",
		"gte":      "{{.FieldName}} must be greater than or equal to {{.Param}}",
		"lte":      "{{.FieldName}} must be less than or equal to {{.Param}}",
		"len":      "{{.FieldName}} must be exactly {{.Param}} characters",
		"numeric":  "{{.FieldName}} must be numeric",
		"alphanum": "{{.FieldName}} must contain only letters and numbers",
		"default":  "{{.FieldName}} is invalid ({{.Tag}})",
	}
)

func init() {
	Validator = validator.New()
}

// SetI18nManager sets the global I18nManager instance for validator translations.
// Once set, all validation error messages will use i18n translations with the "validator." prefix.
// If not set, validator will use DefaultMessages (English).
//
// Parameters:
//   - manager: *i18n.I18nManager - The initialized i18n manager instance
//
// Example:
//
//	i18nMgr, _ := i18n.NewI18nManager(config)
//	validator.SetI18nManager(i18nMgr)
func SetI18nManager(manager *i18n.I18nManager) {
	i18nManager = manager
}

// ValidationError is a custom error type for validation failures.
// It collects multiple validation error messages and provides convenient methods
// to access them individually or as a group, including per-field error details.
//
// Fields:
//   - Messages: Slice of all validation error messages (for backward compatibility)
//   - Errors: Map of field names to their error messages (for detailed error reporting)
//
// Methods:
//   - Error(): Returns all messages joined by semicolon (implements error interface)
//   - First(): Returns the first error message
//   - All(): Returns all error messages as a slice
//   - GetFieldErrors(): Returns map of field names to their error messages
//
// Example:
//
//	if err := ValidateStruct(user, "en"); err != nil {
//	    if verr, ok := err.(*ValidationError); ok {
//	        fmt.Println(verr.First())  // Get first error
//	        fmt.Println(verr.All())    // Get all errors
//	        fmt.Println(verr.GetFieldErrors())  // Get errors per field
//	    }
//	}
type ValidationError struct {
	Messages []string            // All error messages (backward compatibility)
	Errors   map[string][]string // Field name -> error messages mapping
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

// GetFieldErrors returns a map of field names to their error messages.
// This is useful for displaying field-specific errors in UI forms.
//
// Returns:
//   - map[string][]string: Map where keys are field names and values are slices of error messages for that field
//
// Example:
//
//	verr := &ValidationError{
//	    Errors: map[string][]string{
//	        "Email": {"Email is required", "Email must be valid"},
//	        "Password": {"Password is too short"},
//	    },
//	}
//	for field, errs := range verr.GetFieldErrors() {
//	    fmt.Printf("%s: %v\n", field, errs)
//	}
func (ve *ValidationError) GetFieldErrors() map[string][]string {
	return ve.Errors
}

// getLanguageFromContext retrieves the language code from the Fiber context.
// It attempts to get the language set by I18nMiddleware from context locals.
// If not found, it falls back to the default language from i18nManager, or "en" if i18nManager is not set.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context
//
// Returns:
//   - string: Language code (e.g., "en", "id", "zh")
func getLanguageFromContext(c *fiber.Ctx) string {
	if lang, ok := c.Locals("language").(string); ok {
		return lang
	}

	if i18nManager != nil {
		return i18nManager.DefaultLanguage
	}

	return "en" // fallback to English
}

// getFieldName retrieves the field name from json tag if available, otherwise returns the struct field name.
// This ensures consistency between request/response JSON field names and validation error messages.
//
// Parameters:
//   - s: The struct being validated
//   - fieldName: The struct field name from validator
//
// Returns:
//   - string: JSON tag name if exists, otherwise original field name in title case
//
// Example:
//
//	type User struct {
//	    Email string `json:"email" validate:"required"`
//	}
//	// getFieldName will return "email" instead of "Email"
func getFieldName(s interface{}, fieldName string) string {
	// Get the type of the struct
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// If not a struct, return title case of field name
	if val.Kind() != reflect.Struct {
		caser := cases.Title(language.Und)
		return caser.String(fieldName)
	}

	// Try to find the field in the struct
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == fieldName {
			// Check if json tag exists
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" {
				// Handle json tag with options (e.g., "email,omitempty")
				parts := strings.Split(jsonTag, ",")
				if parts[0] != "" && parts[0] != "-" {
					return parts[0]
				}
			}
			// If no json tag or json:"-", return title case of field name
			break
		}
	}

	// Fallback to title case of original field name
	caser := cases.Title(language.Und)
	return caser.String(fieldName)
}

// ValidateStruct validates a struct using validation tags with the default language.
// If i18nManager is set, it uses the default language from i18nManager.
// Otherwise, it uses "en" (English) as the default language.
//
// Parameters:
//   - s: The struct to validate (must have validation tags)
//
// Returns:
//   - error: nil if validation succeeds, *ValidationError if validation fails
//
// Example:
//
//	type User struct {
//	    Email    string `validate:"required,email"`
//	    Password string `validate:"required,min=8"`
//	}
//
//	user := User{Email: "invalid", Password: "123"}
//	if err := ValidateStruct(user); err != nil {
//	    if verr, ok := err.(*ValidationError); ok {
//	        fmt.Println(verr.First())
//	        // Output: Email must be a valid email address
//	    }
//	}
func ValidateStruct(s interface{}) error {
	defaultLang := "en"
	if i18nManager != nil {
		defaultLang = i18nManager.DefaultLanguage
	}
	return ValidateStructWithLang(s, defaultLang)
}

// ValidateStructWithLang validates a struct using validation tags with a specified language.
// It uses the global Validator instance and translates errors based on the specified language.
//
// The function automatically converts field names to title case for better readability.
// If i18nManager is set, it will use i18n translations from locale files with "validator." prefix.
// Otherwise, it falls back to DefaultMessages (English).
//
// If validation succeeds, it returns nil. If validation fails, it returns a ValidationError
// containing all validation error messages in the specified language.
//
// Parameters:
//   - s: The struct to validate (must have validation tags)
//   - lang: Language code for error messages (e.g., "en", "id", "zh")
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
//	if err := ValidateStructWithLang(user, "id"); err != nil {
//	    if verr, ok := err.(*ValidationError); ok {
//	        fmt.Println(verr.First())
//	        // Output: Email harus berupa alamat email yang valid
//	    }
//	}
func ValidateStructWithLang(s interface{}, lang string) error {
	err := Validator.Struct(s)
	if err == nil {
		return nil
	}

	var messages []string
	fieldErrors := make(map[string][]string)

	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		for _, e := range validateErrs {
			// Get field name from json tag if available
			fieldName := getFieldName(s, e.Field())
			message := getUserFriendlyMessage(fieldName, e.Tag(), e.Param(), lang)
			messages = append(messages, message)

			// Add to field errors map using json tag name
			fieldErrors[fieldName] = append(fieldErrors[fieldName], message)
		}
	} else {
		// Jika bukan validation error, kembalikan error asli
		messages = append(messages, err.Error())
	}

	return &ValidationError{
		Messages: messages,
		Errors:   fieldErrors,
	}
}

// ValidateStructWithContext validates a struct using validation tags with language from Fiber context.
// It extracts the language from the Fiber context (set by I18nMiddleware) and uses it for error messages.
// If language is not found in context, it falls back to the default language.
//
// The function automatically converts field names to title case for better readability.
// If i18nManager is set, it will use i18n translations from locale files with "validator." prefix.
// Otherwise, it falls back to DefaultMessages (English).
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context containing language information
//   - s: The struct to validate (must have validation tags)
//
// Returns:
//   - error: nil if validation succeeds, *ValidationError if validation fails
//
// Example:
//
//	app.Post("/users", func(c *fiber.Ctx) error {
//	    var user User
//	    if err := c.BodyParser(&user); err != nil {
//	        return err
//	    }
//
//	    if err := ValidateStructWithContext(c, user); err != nil {
//	        if verr, ok := err.(*ValidationError); ok {
//	            return c.Status(400).JSON(fiber.Map{
//	                "error": verr.First(),
//	            })
//	        }
//	    }
//
//	    return c.JSON(user)
//	})
func ValidateStructWithContext(c *fiber.Ctx, s interface{}) error {
	lang := getLanguageFromContext(c)
	return ValidateStructWithLang(s, lang)
}

// getUserFriendlyMessage generates a user-friendly error message based on validation failure details.
// It uses i18n for translations if i18nManager is set, otherwise falls back to DefaultMessages.
//
// The function:
//   - Uses field name from json tag for consistency with request/response
//   - Uses i18n with "validator." prefix for message keys (e.g., "validator.required")
//   - Falls back to DefaultMessages if i18n is not available or translation not found
//   - Supports template data with FieldName, Param, and Tag placeholders
//
// Parameters:
//   - fieldName: Name of the field (from json tag or struct field name)
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
//	// Returns: "email is required"
//
//	msg := getUserFriendlyMessage("password", "min", "8", "id")
//	// Returns: "password minimal 8 karakter"
func getUserFriendlyMessage(fieldName, tag, param, lang string) string {
	// Prepare template data (fieldName already processed from json tag or title case)
	templateData := map[string]string{
		"FieldName": fieldName,
		"Param":     param,
		"Tag":       tag,
	}

	// Try to get message from i18n if available
	if i18nManager != nil {
		fmt.Println("Tag:", tag)
		messageKey := "validator." + tag
		message := i18nManager.Translate(lang, messageKey, templateData)

		// Check if translation was found (i18n returns the key if not found)
		if !strings.Contains(message, "Missing translation") {
			return message
		}

		// Try default key if specific tag not found
		messageKey = "validator.default"
		message = i18nManager.Translate(lang, messageKey, templateData)
		if !strings.Contains(message, "Missing translation") {
			return message
		}
	}

	fmt.Println("Falling back to default English messages")

	// Fallback to default English messages
	template, exists := DefaultMessages[tag]
	if !exists {
		template = DefaultMessages["default"]
	}

	// Simple template replacement for default messages
	message := strings.ReplaceAll(template, "{{.FieldName}}", fieldName)
	message = strings.ReplaceAll(message, "{{.Param}}", param)
	message = strings.ReplaceAll(message, "{{.Tag}}", tag)

	return message
}
