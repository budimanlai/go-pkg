# Validator Package Documentation

Package `validator` menyediakan validasi struct dengan pesan error yang user-friendly dan mendukung multi bahasa.

## Fitur
- Validasi struct menggunakan tag `validate` dari library `go-playground/validator`.
- Pesan error yang mudah dipahami oleh user.
- Mendukung multi bahasa (default: Indonesia dan Inggris).
- Kemampuan untuk menambah atau mengupdate bahasa custom.
- Tipe error custom dengan method `First()` dan `All()`.

## Instalasi
Package ini adalah bagian dari `github.com/budimanlai/go-pkg`. Import sebagai:

```go
import "github.com/budimanlai/go-pkg/validator"
```

## Cara Pakai

### 1. Definisikan Struct dengan Tag Validate
```go
type User struct {
    Name  string `validate:"required"`
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=0,lte=130"`
}
```

### 2. Validasi Struct
```go
user := &User{
    Name:  "John Doe",
    Email: "john@example.com",
    Age:   25,
}

err := validator.ValidateStruct(user, "id") // "id" untuk bahasa Indonesia, "en" untuk Inggris
if err != nil {
    fmt.Println("Validasi error:", err)
    if valErr, ok := err.(*validator.ValidationError); ok {
        fmt.Println("Error pertama:", valErr.First())
        fmt.Println("Semua Error:", valErr.All())
    }
} else {
    fmt.Println("Validasi berhasil!")
}
```

### 3. Handle Error
- Jika validasi berhasil, `err` adalah `nil`.
- Jika gagal, `err` adalah `*ValidationError` dengan pesan error dalam bahasa yang dipilih.
- Gunakan `err.First()` untuk error pertama, `err.All()` untuk semua error.

### 4. Multi Bahasa
- Bahasa default: "id" (Indonesia), "en" (Inggris).
- Untuk bahasa lain, gunakan `AddLanguage` atau `UpdateLanguage`.

```go
// Menambah bahasa baru
javaneseMessages := map[string]string{
    "required": "%s kudu diisi",
    "email":    "%s kudu alamat email sing bener",
    "default":  "%s ora valid (%s)",
}
validator.AddLanguage("jv", javaneseMessages)

// Menggunakan bahasa baru
err := validator.ValidateStruct(user, "jv")
```

### 5. Custom Pesan
- Akses langsung `validator.Messages` untuk modifikasi.
- Atau gunakan `UpdateLanguage` untuk bahasa tertentu.

```go
validator.UpdateLanguage("id", map[string]string{
    "required": "%s harus diisi",
    // ...
})
```

## Tag Validasi yang Didukung
- `required`: Wajib diisi.
- `email`: Harus email valid.
- `min`: Minimal panjang (untuk string).
- `max`: Maksimal panjang.
- `gte`: Lebih besar atau sama dengan.
- `lte`: Lebih kecil atau sama dengan.
- `len`: Panjang tepat.
- `numeric`: Harus angka.
- `alphanum`: Hanya huruf dan angka.
- Lainnya: Gunakan "default" untuk fallback.

## Cara Tests

### Menjalankan Unit Tests
```bash
# Dari root project
go test ./validator/

# Dengan verbose output
go test -v ./validator/

# Dengan coverage
go test -cover ./validator/
```

### Test Files
- `validator/validator_test.go`: Berisi unit tests untuk semua fungsi.
- Tests mencakup validasi valid/invalid, multi bahasa, custom bahasa, dan method error.

### Contoh Output Test
```
=== RUN   TestValidateStruct_Valid
--- PASS: TestValidateStruct_Valid (0.00s)
=== RUN   TestValidateStruct_Invalid_ID
--- PASS: TestValidateStruct_Invalid_ID (0.00s)
...
PASS
ok      github.com/budimanlai/go-pkg/validator  0.166s
```

## Dependencies
- `github.com/go-playground/validator/v10`
- `golang.org/x/text/cases`
- `golang.org/x/text/language`

## Lisensi
Sesuai dengan project utama.