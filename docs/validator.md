# Validator Package Documentation

Package `validator` menyediakan validasi struct dengan pesan error yang user-friendly dan mendukung multi bahasa melalui integrasi dengan package `i18n`.

## Fitur
- Validasi struct menggunakan tag `validate` dari library `go-playground/validator`.
- Pesan error yang mudah dipahami oleh user.
- **Integrasi dengan package i18n** untuk multi bahasa yang fleksibel.
- Bahasa dikelola melalui file JSON di folder `locales` dengan prefix `validator.`.
- Fallback ke bahasa Inggris default jika i18n tidak diset.
- Tipe error custom dengan method `First()` dan `All()`.

## Instalasi
Package ini adalah bagian dari `github.com/budimanlai/go-pkg`. Import sebagai:

```go
import "github.com/budimanlai/go-pkg/validator"
```

## Cara Pakai

### 1. Setup I18n (Opsional tapi Disarankan)

Jika ingin menggunakan multi bahasa, setup i18n terlebih dahulu:

```go
import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/validator"
    "golang.org/x/text/language"
)

// Setup I18n
i18nConfig := i18n.I18nConfig{
    DefaultLanguage: language.English,
    SupportedLangs:  []string{"en", "id", "zh"},
    LocalesPath:     "locales",
}

i18nManager, err := i18n.NewI18nManager(i18nConfig)
if err != nil {
    log.Fatal(err)
}

// Set I18nManager ke validator
validator.SetI18nManager(i18nManager)
```

### 2. Definisikan Struct dengan Tag Validate

**PENTING:** Gunakan tag `json` untuk konsistensi field name antara request/response dan error messages.

```go
type User struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Age      int    `json:"age" validate:"gte=18,lte=130"`
    Username string `json:"username" validate:"required,alphanum"`
}
```

**Keuntungan menggunakan json tag:**
- Field name di error message sama dengan field name di JSON request/response
- Consistency untuk frontend developer
- Lebih mudah mapping error ke form field di UI

**Contoh tanpa json tag:**
```json
{
  "errors": {
    "Email": ["Email is required"]  // Title case
  }
}
```

**Contoh dengan json tag:**
```json
{
  "errors": {
    "email": ["email is required"]  // Lowercase, sama dengan request
  }
}
```

### 3. Validasi Struct dengan Bahasa Tertentu

Ada 3 cara untuk melakukan validasi:

#### a. ValidateStruct() - Menggunakan Bahasa Default

```go
user := &User{
    Email:    "invalid-email",
    Password: "123",
}

// Validasi dengan bahasa default (dari i18nManager.DefaultLanguage)
// Jika i18nManager tidak diset, menggunakan "en" (English)
err := validator.ValidateStruct(user)
if err != nil {
    if valErr, ok := err.(*validator.ValidationError); ok {
        fmt.Println(valErr.First())
    }
}
```

#### b. ValidateStructWithLang() - Menggunakan Bahasa Spesifik

```go
type User struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

user := &User{
    Email:    "invalid-email",
    Password: "123",
}

// Validasi dengan bahasa Indonesia
err := validator.ValidateStructWithLang(user, "id")
if err != nil {
    if valErr, ok := err.(*validator.ValidationError); ok {
        fmt.Println(valErr.First())
        // Output: email harus berupa alamat email yang valid
        
        // Get field errors
        for field, errs := range valErr.GetFieldErrors() {
            fmt.Printf("%s: %v\n", field, errs)
        }
        // Output:
        // email: [email harus berupa alamat email yang valid]
        // password: [password minimal 8 karakter]
    }
}

// Validasi dengan bahasa Inggris
err = validator.ValidateStructWithLang(user, "en")
if err != nil {
    if valErr, ok := err.(*validator.ValidationError); ok {
        fmt.Println(valErr.First())
        // Output: email must be a valid email address
    }
}
```

#### c. ValidateStructWithContext() - Menggunakan Bahasa dari Fiber Context

Cocok untuk API endpoint yang menggunakan Fiber framework dengan I18nMiddleware:

```go
import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/validator"
    "github.com/gofiber/fiber/v2"
)

app := fiber.New()

// Setup I18n Middleware
app.Use(i18n.I18nMiddleware(i18nConfig))

// Endpoint with automatic language detection
app.Post("/users", func(c *fiber.Ctx) error {
    var user User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Bahasa otomatis diambil dari:
    // 1. Query parameter ?lang=id
    // 2. Header Accept-Language
    // 3. Default language
    if err := validator.ValidateStructWithContext(c, user); err != nil {
        if valErr, ok := err.(*validator.ValidationError); ok {
            return c.Status(400).JSON(fiber.Map{
                "success": false,
                "message": valErr.First(),
                "errors":  valErr.All(),
            })
        }
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    user,
    })
})
```

### 4. Handle Error
- Jika validasi berhasil, `err` adalah `nil`.
- Jika gagal, `err` adalah `*ValidationError` dengan pesan error dalam bahasa yang dipilih.
- Gunakan `err.First()` untuk error pertama, `err.All()` untuk semua error.
- Gunakan `err.GetFieldErrors()` untuk error per field (berguna untuk UI).

```go
if err := validator.ValidateStruct(user); err != nil {
    if valErr, ok := err.(*validator.ValidationError); ok {
        // Get first error only
        fmt.Println(valErr.First())
        
        // Get all errors
        for _, msg := range valErr.All() {
            fmt.Println(msg)
        }
        
        // Get error as string (all errors joined with semicolon)
        fmt.Println(valErr.Error())
        
        // Get errors by field (untuk UI form validation)
        fieldErrors := valErr.GetFieldErrors()
        for field, messages := range fieldErrors {
            fmt.Printf("Field %s errors:\n", field)
            for _, msg := range messages {
                fmt.Printf("  - %s\n", msg)
            }
        }
        // Output:
        // Field Email errors:
        //   - Email is required
        //   - Email must be a valid email address
        // Field Password errors:
        //   - Password must be at least 8 characters
    }
}
```

### 5. Integrasi dengan Response Package

Untuk API, gunakan `response.ValidationErrorI18n()` yang secara otomatis menggunakan bahasa dari context dan mengembalikan JSON dengan detail field errors:

```go
import (
    "github.com/budimanlai/go-pkg/response"
    "github.com/budimanlai/go-pkg/validator"
)

app.Post("/users", func(c *fiber.Ctx) error {
    var user User
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }

    // Validasi dengan response otomatis
    if err := validator.ValidateStructWithContext(c, user); err != nil {
        return response.ValidationErrorI18n(c, err)
    }

    // Proses user...
    return response.Success(c, user)
})
```

Response JSON:
```json
{
    "meta": {
        "success": false,
        "message": "email is required",
        "errors": {
            "email": [
                "email is required",
                "email must be a valid email address"
            ],
            "password": [
                "password must be at least 8 characters"
            ]
        }
    },
    "data": null
}
```

**Catatan:** Field names menggunakan nama dari json tag (lowercase), bukan nama struct field (Title Case).

## Perbandingan 3 Fungsi Validasi

| Fungsi | Bahasa | Use Case |
|--------|--------|----------|
| `ValidateStruct(s)` | Default dari i18nManager atau "en" | Validasi tanpa perlu specify bahasa |
| `ValidateStructWithLang(s, lang)` | Explicit dari parameter | Validasi dengan bahasa spesifik |
| `ValidateStructWithContext(c, s)` | Dari Fiber context | API endpoint dengan I18nMiddleware |

### Contoh Penggunaan Lengkap dengan Fiber

```go
package main

import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/validator"
    "github.com/gofiber/fiber/v2"
    "golang.org/x/text/language"
)

func main() {
    // Setup I18n
    i18nConfig := i18n.I18nConfig{
        DefaultLanguage: language.English,
        SupportedLangs:  []string{"en", "id", "zh"},
        LocalesPath:     "locales",
    }
    i18nManager, _ := i18n.NewI18nManager(i18nConfig)
    validator.SetI18nManager(i18nManager)

    app := fiber.New()
    app.Use(i18n.I18nMiddleware(i18nConfig))

    // Endpoint 1: Language from context (recommended for APIs)
    app.Post("/users", func(c *fiber.Ctx) error {
        var user User
        c.BodyParser(&user)
        
        if err := validator.ValidateStructWithContext(c, user); err != nil {
            return handleValidationError(c, err)
        }
        return c.JSON(user)
    })

    // Endpoint 2: Always English
    app.Post("/users/en", func(c *fiber.Ctx) error {
        var user User
        c.BodyParser(&user)
        
        if err := validator.ValidateStructWithLang(user, "en"); err != nil {
            return handleValidationError(c, err)
        }
        return c.JSON(user)
    })

    // Endpoint 3: Default language
    app.Post("/users/default", func(c *fiber.Ctx) error {
        var user User
        c.BodyParser(&user)
        
        if err := validator.ValidateStruct(user); err != nil {
            return handleValidationError(c, err)
        }
        return c.JSON(user)
    })

    app.Listen(":3000")
}

func handleValidationError(c *fiber.Ctx, err error) error {
    if valErr, ok := err.(*validator.ValidationError); ok {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": valErr.First(),
            "errors":  valErr.All(),
        })
    }
    return err
}
```

## Menambah Bahasa Baru

Untuk menambah bahasa baru, cukup tambahkan file JSON di folder `locales` dengan message key yang menggunakan prefix `validator.`:

### 1. Buat File Locale Baru

Contoh: `locales/fr.json` untuk bahasa Perancis:

```json
{
    "validator.required": "{{.FieldName}} est requis",
    "validator.email": "{{.FieldName}} doit être une adresse email valide",
    "validator.min": "{{.FieldName}} doit comporter au moins {{.Param}} caractères",
    "validator.max": "{{.FieldName}} doit comporter au plus {{.Param}} caractères",
    "validator.gte": "{{.FieldName}} doit être supérieur ou égal à {{.Param}}",
    "validator.lte": "{{.FieldName}} doit être inférieur ou égal à {{.Param}}",
    "validator.len": "{{.FieldName}} doit comporter exactement {{.Param}} caractères",
    "validator.numeric": "{{.FieldName}} doit être numérique",
    "validator.alphanum": "{{.FieldName}} ne doit contenir que des lettres et des chiffres",
    "validator.default": "{{.FieldName}} n'est pas valide ({{.Tag}})"
}
```

### 2. Update I18n Config

```go
i18nConfig := i18n.I18nConfig{
    DefaultLanguage: language.English,
    SupportedLangs:  []string{"en", "id", "zh", "fr"}, // Tambahkan "fr"
    LocalesPath:     "locales",
}
```

### 3. Gunakan Bahasa Baru

```go
err := validator.ValidateStruct(user, "fr")
// Pesan error akan dalam bahasa Perancis
```

## Template Placeholders

Message template mendukung placeholder berikut:
- `{{.FieldName}}`: Nama field yang divalidasi (otomatis dikonversi ke title case)
- `{{.Param}}`: Parameter dari validation tag (contoh: "8" untuk `min=8`)
- `{{.Tag}}`: Nama validation tag (contoh: "required", "email")

## Tag Validasi yang Didukung

Message default sudah tersedia untuk tag berikut:
- `required`: Wajib diisi.
- `email`: Harus email valid.
- `min`: Minimal panjang (untuk string) atau nilai (untuk number).
- `max`: Maksimal panjang atau nilai.
- `gte`: Lebih besar atau sama dengan.
- `lte`: Lebih kecil atau sama dengan.
- `len`: Panjang tepat.
- `numeric`: Harus angka.
- `alphanum`: Hanya huruf dan angka.
- `default`: Pesan fallback untuk tag yang tidak memiliki message khusus.

Untuk tag validasi lainnya yang didukung oleh `go-playground/validator`, lihat [dokumentasi resmi](https://github.com/go-playground/validator).

## Penggunaan Tanpa I18n

Jika tidak memanggil `SetI18nManager()`, validator akan menggunakan pesan default dalam bahasa Inggris:

```go
import "github.com/budimanlai/go-pkg/validator"

type User struct {
    Email string `validate:"required,email"`
}

user := User{Email: "invalid"}
err := validator.ValidateStruct(user, "en") // Bahasa akan diabaikan, selalu English
if err != nil {
    fmt.Println(err)
    // Output: Email must be a valid email address
}
```

## Contoh Lengkap

### Example 1: Tanpa I18n (Simple)
Lihat file `examples/validator_without_i18n.go`

```bash
go run examples/validator_without_i18n.go
```

### Example 2: Dengan I18n (Multi-Language)
Lihat file `examples/validator_with_i18n.go`

```bash
go run examples/validator_with_i18n.go
```

### Example 3: Dengan Fiber (REST API)
Lihat file `examples/validator_with_fiber.go`

```bash
go run examples/validator_with_fiber.go

# Test dengan curl
curl -X POST http://localhost:3000/users?lang=id \
  -H "Content-Type: application/json" \
  -d '{"email":"invalid","password":"123","age":15,"username":""}'
```

## API Reference

### Functions

#### `SetI18nManager(manager *i18n.I18nManager)`
Set global I18nManager instance untuk validator translations.

#### `ValidateStruct(s interface{}) error`
Validasi struct dengan bahasa default. Returns `*ValidationError` jika validasi gagal.

**Default Language:**
- Jika i18nManager diset: Menggunakan `i18nManager.DefaultLanguage`
- Jika i18nManager tidak diset: Menggunakan "en" (English)

#### `ValidateStructWithLang(s interface{}, lang string) error`
Validasi struct dengan bahasa spesifik. Returns `*ValidationError` jika validasi gagal.

**Parameters:**
- `s`: Struct to validate
- `lang`: Language code (e.g., "en", "id", "zh")

#### `ValidateStructWithContext(c *fiber.Ctx, s interface{}) error`
Validasi struct dengan bahasa dari Fiber context. Returns `*ValidationError` jika validasi gagal.

Bahasa diambil dari (berurutan):
1. Query parameter `?lang=id`
2. Header `Accept-Language`
3. Default language

**Parameters:**
- `c`: Fiber context containing language information
- `s`: Struct to validate

### Types

#### `ValidationError`
Custom error type untuk validation failures.

**Methods:**
- `Error() string`: Returns semua error messages joined dengan semicolon
- `First() string`: Returns error message pertama
- `All() []string`: Returns semua error messages sebagai slice

### Variables

#### `Validator *validator.Validate`
Global validator instance yang dapat digunakan langsung.

#### `DefaultMessages map[string]string`
Fallback validation messages dalam bahasa Inggris.

## Migration dari Versi Lama

Jika sebelumnya menggunakan hardcoded `Messages` map:

### Sebelum (Old)
```go
validator.AddLanguage("fr", map[string]string{
    "required": "%s est requis",
    "email":    "%s doit être valide",
})

err := validator.ValidateStruct(user, "fr")
```

### Sekarang (New)
```go
// 1. Buat file locales/fr.json
{
    "validator.required": "{{.FieldName}} est requis",
    "validator.email": "{{.FieldName}} doit être valide"
}

// 2. Setup i18n dengan bahasa baru
i18nConfig := i18n.I18nConfig{
    DefaultLanguage: language.English,
    SupportedLangs:  []string{"en", "id", "fr"},
    LocalesPath:     "locales",
}
i18nManager, _ := i18n.NewI18nManager(i18nConfig)
validator.SetI18nManager(i18nManager)

// 3. Validasi seperti biasa
err := validator.ValidateStruct(user, "fr")
```

**Keuntungan:**
- Lebih mudah maintain (messages di file JSON)
- Tidak perlu recompile untuk update/tambah bahasa
- Konsisten dengan sistem i18n di aplikasi
- Support template data yang lebih fleksibel

## Troubleshooting

### Pesan Error Tidak Muncul dalam Bahasa yang Diinginkan

**Solusi:**
1. Pastikan `SetI18nManager()` sudah dipanggil
2. Pastikan file locale untuk bahasa tersebut ada di folder `locales`
3. Pastikan message key menggunakan prefix `validator.` (contoh: `validator.required`)
4. Pastikan bahasa sudah ditambahkan di `SupportedLangs` config i18n

### Pesan Error Masih dalam Bahasa Inggris

Jika i18n tidak diset atau translation tidak ditemukan, validator akan fallback ke `DefaultMessages` yang dalam bahasa Inggris.

## Lihat Juga

- [I18n Documentation](./i18n.md) - Dokumentasi package i18n
- [go-playground/validator](https://github.com/go-playground/validator) - Library validator yang digunakan

PASS
ok      github.com/budimanlai/go-pkg/validator  0.166s
```

## Dependencies
- `github.com/go-playground/validator/v10`
- `golang.org/x/text/cases`
- `golang.org/x/text/language`

## Lisensi
Sesuai dengan project utama.