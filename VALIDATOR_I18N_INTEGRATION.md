# Validator I18n Integration - Summary

## Perubahan yang Dilakukan

### 1. Update Locale Files
Menambahkan message keys untuk validator dengan prefix `validator.` di semua file locale:

**Files Updated:**
- `locales/en.json` - Added English validator messages
- `locales/id.json` - Added Indonesian validator messages  
- `locales/zh.json` - Added Chinese validator messages

**Message Keys Added:**
- `validator.required` - Field wajib diisi
- `validator.email` - Validasi email
- `validator.min` - Minimal karakter/nilai
- `validator.max` - Maksimal karakter/nilai
- `validator.gte` - Greater than or equal
- `validator.lte` - Less than or equal
- `validator.len` - Panjang exact
- `validator.numeric` - Harus numeric
- `validator.alphanum` - Hanya alphanumeric
- `validator.default` - Fallback message

**Template Format:**
Menggunakan template placeholders:
- `{{.FieldName}}` - Nama field (auto title case)
- `{{.Param}}` - Parameter dari validation tag
- `{{.Tag}}` - Nama validation tag

### 2. Refactor validator.go

**Removed:**
- Hardcoded `Messages` map dengan multi-language
- Functions: `AddLanguage()`, `UpdateLanguage()`, `GetLanguages()`

**Added:**
- Import `github.com/budimanlai/go-pkg/i18n`
- Variable `i18nManager *i18n.I18nManager` - Global i18n manager
- Variable `DefaultMessages map[string]string` - Fallback English messages
- Function `SetI18nManager(manager *i18n.I18nManager)` - Set i18n instance

**Modified:**
- `getUserFriendlyMessage()` - Now uses i18n.Translate() with fallback to DefaultMessages
- Updated all function documentations

**Benefits:**
✅ Tidak perlu hardcode messages untuk setiap bahasa
✅ Menambah bahasa baru hanya perlu tambah file JSON di locales/
✅ Konsisten dengan sistem i18n di seluruh aplikasi
✅ Support template data yang lebih fleksibel
✅ Tidak perlu recompile untuk update/tambah bahasa

### 3. Created Examples

**examples/validator_without_i18n.go:**
- Demo penggunaan validator tanpa i18n setup
- Menggunakan DefaultMessages (English)

**examples/validator_with_i18n.go:**
- Demo penggunaan validator dengan i18n
- Menunjukkan validasi dalam 3 bahasa (en, id, zh)

### 4. Updated Documentation

**docs/validator.md:**
- Updated untuk mencerminkan integrasi dengan i18n
- Menambahkan section "Menambah Bahasa Baru"
- Menambahkan section "Template Placeholders"
- Menambahkan section "Penggunaan Tanpa I18n"
- Menambahkan section "Migration dari Versi Lama"
- Menambahkan section "Troubleshooting"

## Cara Penggunaan

### Setup (Dengan I18n - Recommended)

```go
// 1. Setup I18n
i18nConfig := i18n.I18nConfig{
    DefaultLanguage: language.English,
    SupportedLangs:  []string{"en", "id", "zh"},
    LocalesPath:     "locales",
}
i18nManager, _ := i18n.NewI18nManager(i18nConfig)

// 2. Set I18nManager ke validator
validator.SetI18nManager(i18nManager)

// 3. Validasi dengan bahasa tertentu
err := validator.ValidateStruct(user, "id") // Indonesian
err := validator.ValidateStruct(user, "en") // English
err := validator.ValidateStruct(user, "zh") // Chinese
```

### Tanpa I18n (Fallback)

```go
// Langsung pakai, akan menggunakan DefaultMessages (English)
err := validator.ValidateStruct(user, "en")
```

## Menambah Bahasa Baru

### Langkah 1: Buat File Locale

Create `locales/[lang_code].json`:

```json
{
    "validator.required": "{{.FieldName}} [required message in new language]",
    "validator.email": "{{.FieldName}} [email message in new language]",
    ...
}
```

### Langkah 2: Update I18n Config

```go
i18nConfig := i18n.I18nConfig{
    SupportedLangs: []string{"en", "id", "zh", "new_lang"},
    ...
}
```

### Langkah 3: Gunakan

```go
err := validator.ValidateStruct(user, "new_lang")
```

## Testing

```bash
# Test without i18n
go run examples/validator_without_i18n.go

# Test with i18n (multiple languages)
go run examples/validator_with_i18n.go
```

## Migration Guide

Untuk project yang sudah menggunakan validator versi lama:

1. **Hapus** semua pemanggilan `AddLanguage()` atau `UpdateLanguage()`
2. **Buat** file locale JSON untuk setiap bahasa
3. **Setup** i18n manager dan panggil `SetI18nManager()`
4. **Update** message format dari `%s` ke template `{{.FieldName}}`, `{{.Param}}`

## Notes

- Jika `SetI18nManager()` tidak dipanggil, validator tetap berfungsi dengan DefaultMessages (English)
- Template placeholder case-sensitive: gunakan `{{.FieldName}}` bukan `{{.fieldname}}`
- Message key harus menggunakan prefix `validator.` (contoh: `validator.required`)
- Untuk validation tag yang tidak ada message-nya, akan menggunakan `validator.default`

## Files Changed

- ✅ `locales/en.json` - Added validator messages
- ✅ `locales/id.json` - Added validator messages
- ✅ `locales/zh.json` - Added validator messages
- ✅ `validator/validator.go` - Integrated with i18n
- ✅ `examples/validator_without_i18n.go` - Example without i18n
- ✅ `examples/validator_with_i18n.go` - Example with i18n
- ✅ `docs/validator.md` - Updated documentation
