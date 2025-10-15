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
	Validator *validator.Validate
	// Messages berisi pesan error per bahasa per tag
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

// ValidationError adalah tipe error custom untuk validasi
type ValidationError struct {
	Messages []string
}

// Error mengimplementasikan interface error
func (ve *ValidationError) Error() string {
	return strings.Join(ve.Messages, "; ")
}

// First mengembalikan pesan error pertama
func (ve *ValidationError) First() string {
	if len(ve.Messages) > 0 {
		return ve.Messages[0]
	}
	return ""
}

// All mengembalikan semua pesan error
func (ve *ValidationError) All() []string {
	return ve.Messages
}

// ValidateStruct memvalidasi struct dan mengembalikan ValidationError jika ada error
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

// getUserFriendlyMessage mengembalikan pesan error yang user-friendly berdasarkan field, tag, param, dan bahasa
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

// AddLanguage menambahkan bahasa baru dengan pesan custom
func AddLanguage(lang string, messages map[string]string) {
	if _, exists := Messages[lang]; !exists {
		Messages[lang] = messages
	}
}

// UpdateLanguage mengupdate pesan bahasa yang ada atau menambah jika belum ada
func UpdateLanguage(lang string, messages map[string]string) {
	Messages[lang] = messages
}

// GetLanguages mengembalikan list bahasa yang tersedia
func GetLanguages() []string {
	var langs []string
	for lang := range Messages {
		langs = append(langs, lang)
	}
	return langs
}
