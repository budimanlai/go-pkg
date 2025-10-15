# Response Package Documentation

Paket response ini menyediakan helper functions untuk membuat HTTP response standar dalam format JSON untuk aplikasi Fiber. Mendukung response biasa dan response dengan internasionalisasi (i18n).

## Instalasi

Pastikan Anda memiliki dependensi yang diperlukan:

```bash
go get github.com/gofiber/fiber/v2
```

Untuk fitur i18n, pastikan paket i18n sudah diinstall:

```bash
go get github.com/budimanlai/go-pkg/i18n
```

## Struktur Response

Semua response mengembalikan JSON dengan struktur:

```json
{
    "meta": {
        "success": true|false,
        "message": "response message"
    },
    "data": null|object
}
```

- `meta.success`: Boolean indicating success status
- `meta.message`: Response message string
- `data`: Response data (null for error responses, object for success)

## Fungsi Response Standar

### Success

Mengembalikan response sukses dengan status 200.

```go
func Success(c *fiber.Ctx, message string, data interface{}) error
```

**Contoh:**

```go
return response.Success(c, "Operation successful", map[string]string{"user_id": "123"})
```

**Output:**

```json
{
    "meta": {
        "success": true,
        "message": "Operation successful"
    },
    "data": {
        "user_id": "123"
    }
}
```

### Error

Mengembalikan response error dengan status code custom.

```go
func Error(c *fiber.Ctx, code int, message string) error
```

**Contoh:**

```go
return response.Error(c, 500, "Internal server error")
```

**Output:**

```json
{
    "meta": {
        "success": false,
        "message": "Internal server error"
    },
    "data": null
}
```

### BadRequest

Mengembalikan response bad request dengan status 400.

```go
func BadRequest(c *fiber.Ctx, message string) error
```

**Contoh:**

```go
return response.BadRequest(c, "Invalid input data")
```

### NotFound

Mengembalikan response not found dengan status 404.

```go
func NotFound(c *fiber.Ctx, message string) error
```

**Contoh:**

```go
return response.NotFound(c, "Resource not found")
```

## Fungsi Response dengan I18n

Fungsi ini menggunakan paket i18n untuk menerjemahkan pesan berdasarkan bahasa yang ditentukan dalam context Fiber.

### Setup I18n

Sebelum menggunakan fungsi i18n, Anda perlu setup I18nManager:

```go
import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/response"
    "golang.org/x/text/language"
)

i18nConfig := i18n.I18nConfig{
    DefaultLanguage: language.Indonesian,
    SupportedLangs:  []string{"en", "id", "zh"},
    LocalesPath:     "locales",
}

i18nManager, err := i18n.NewI18nManager(i18nConfig)
if err != nil {
    panic(err)
}

response.SetI18nManager(i18nManager)
```

### SuccessI18n

Mengembalikan response sukses dengan pesan yang diterjemahkan.

```go
func SuccessI18n(c *fiber.Ctx, messageID string, data interface{}) error
```

**Contoh:**

```go
return response.SuccessI18n(c, "welcome", map[string]string{"user": "John"})
```

Jika bahasa context adalah "id", akan menggunakan pesan dari `locales/id.json`.

### ErrorI18n

Mengembalikan response error dengan pesan yang diterjemahkan.

```go
func ErrorI18n(c *fiber.Ctx, code int, messageID string, data interface{}) error
```

**Contoh:**

```go
return response.ErrorI18n(c, 400, "invalid_input", map[string]string{"field": "email"})
```

### BadRequestI18n

Mengembalikan bad request dengan pesan yang diterjemahkan.

```go
func BadRequestI18n(c *fiber.Ctx, messageID string, data interface{}) error
```

### NotFoundI18n

Mengembalikan not found dengan pesan yang diterjemahkan.

```go
func NotFoundI18n(c *fiber.Ctx, messageID string) error
```

## Bahasa Context

Fungsi i18n menggunakan bahasa dari `c.Locals("language")`. Jika tidak ada, akan fallback ke bahasa default dari I18nManager.

Untuk set bahasa, Anda bisa menggunakan middleware i18n atau set manual:

```go
c.Locals("language", "id") // Set bahasa Indonesia
```

## Contoh Lengkap

Lihat `tests/main.go` untuk contoh penggunaan lengkap:

```go
package main

import (
    pkg_i18n "github.com/budimanlai/go-pkg/i18n"
    pkg_response "github.com/budimanlai/go-pkg/response"
    "github.com/gofiber/fiber/v2"
    "golang.org/x/text/language"
)

func main() {
    i18nConfig := pkg_i18n.I18nConfig{
        DefaultLanguage: language.Indonesian,
        SupportedLangs:  []string{"en", "id", "zh"},
        LocalesPath:     "locales",
    }

    app := fiber.New()
    i18nManager, err := pkg_i18n.NewI18nManagerWithFiber(app, i18nConfig)
    if err != nil {
        panic(err)
    }
    pkg_response.SetI18nManager(i18nManager)

    app.Get("/", func(c *fiber.Ctx) error {
        return pkg_response.Success(c, "Success", map[string]string{"foo": "bar"})
    })

    app.Get("/i18n", func(c *fiber.Ctx) error {
        // Menggunakan query ?lang=zh atau header Accept-Language
        return pkg_response.SuccessI18n(c, "sukses", map[string]string{"Name": "Budiman"})
    })

    app.Listen(":3000")
}
```

## Testing

Jalankan unit tests dengan:

```bash
go test ./response
```

Tests mencakup:
- Response standar (Success, Error, BadRequest, NotFound)
- Response dengan i18n
- Verifikasi struktur JSON response
- Status code yang benar

## Catatan

- Jika I18nManager tidak diset, fungsi i18n akan fallback ke fungsi standar dengan messageID sebagai message langsung.
- Pastikan file locales tersedia dan berisi messageID yang digunakan.