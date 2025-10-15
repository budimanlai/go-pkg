# Helpers Package Documentation

Paket helpers menyediakan berbagai utility functions untuk operasi sehari-hari dalam pengembangan Go, termasuk pointer manipulation, JSON handling, dan string utilities.

## Pointer Utilities

### Pointer

Membuat pointer dari value apapun.

```go
func Pointer[T any](v T) *T
```

**Contoh:**

```go
val := 42
ptr := helpers.Pointer(val)
// ptr adalah *int dengan value 42

name := "Alice"
namePtr := helpers.Pointer(name)
// namePtr adalah *string dengan value "Alice"
```

### DerefPointer

Dereference pointer dengan fallback default value jika nil.

```go
func DerefPointer[T any](p *T, defaultValue T) T
```

**Contoh:**

```go
var ptr *int
val := helpers.DerefPointer(ptr, 100)
// val = 100 (karena ptr nil)

ptr = helpers.Pointer(42)
val = helpers.DerefPointer(ptr, 100)
// val = 42
```

## JSON Utilities

### UnmarshalTo

Unmarshal JSON string ke struct atau type apapun.

```go
func UnmarshalTo[T any](jsonString string) (T, error)
```

**Contoh:**

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

jsonStr := `{"name":"Bob","age":30}`
person, err := helpers.UnmarshalTo[Person](jsonStr)
if err != nil {
    panic(err)
}
// person.Name = "Bob", person.Age = 30

// Untuk map
data, err := helpers.UnmarshalTo[map[string]interface{}](jsonStr)
// data["name"] = "Bob", data["age"] = 30.0
```

## String Utilities

### GenerateTrxID

Generate transaction ID unik berdasarkan timestamp dan random number.

```go
func GenerateTrxID() string
```

**Format:** YYMMDDHHMMSS + 4 digit random (total 16 karakter)

**Contoh:**

```go
id := helpers.GenerateTrxID()
// Output: "2510151430521234" (contoh)
```

### GenerateTrxIDWithPrefix

Generate transaction ID dengan prefix.

```go
func GenerateTrxIDWithPrefix(prefix string) string
```

**Contoh:**

```go
id := helpers.GenerateTrxIDWithPrefix("TXN")
// Output: "TXN2510151430521234"
```

### GenerateTrxIDWithSuffix

Generate transaction ID dengan suffix.

```go
func GenerateTrxIDWithSuffix(suffix string) string
```

**Contoh:**

```go
id := helpers.GenerateTrxIDWithSuffix("END")
// Output: "2510151430521234END"
```

### GenerateMessageID

Generate UUID untuk message ID.

```go
func GenerateMessageID() string
```

**Contoh:**

```go
id := helpers.GenerateMessageID()
// Output: "550e8400-e29b-41d4-a716-446655440000"
```

### GenerateUniqueID

Generate unique ID pendek (8 karakter pertama dari UUID).

```go
func GenerateUniqueID() string
```

**Contoh:**

```go
id := helpers.GenerateUniqueID()
// Output: "550e8400"
```

### NormalizePhoneNumber

Normalize nomor telepon dengan country code tanpa tanda +.

```go
func NormalizePhoneNumber(phone string) string
```

**Rules:**
- Remove prefix `+`
- Jika start dengan `0`, ganti dengan `62` (Indonesia)
- Jika tidak ada country code (62, 1, 65), default ke `62` (Indonesia)
- Support country codes: 62 (Indonesia), 1 (US), 65 (Singapore)

**Contoh:**

```go
// Indonesia
helpers.NormalizePhoneNumber("+628123456789") // "628123456789"
helpers.NormalizePhoneNumber("08123456789")   // "628123456789"
helpers.NormalizePhoneNumber("8123456789")    // "628123456789"

// Singapore
helpers.NormalizePhoneNumber("+658123456789") // "658123456789"

// US
helpers.NormalizePhoneNumber("+18123456789")  // "18123456789"

// No country code
helpers.NormalizePhoneNumber("23456789")      // "6223456789"
```

## Contoh Lengkap

```go
package main

import (
    "fmt"
    "github.com/budimanlai/go-pkg/helpers"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Phone string `json:"phone"`
}

func main() {
    // Generate IDs
    trxID := helpers.GenerateTrxIDWithPrefix("ORDER")
    msgID := helpers.GenerateMessageID()
    uniqueID := helpers.GenerateUniqueID()

    fmt.Printf("Transaction ID: %s\n", trxID)
    fmt.Printf("Message ID: %s\n", msgID)
    fmt.Printf("Unique ID: %s\n", uniqueID)

    // Normalize phone
    phone := helpers.NormalizePhoneNumber("08123456789")
    fmt.Printf("Normalized phone: %s\n", phone)

    // Pointer utilities
    name := "Alice"
    namePtr := helpers.Pointer(name)
    derefName := helpers.DerefPointer(namePtr, "Unknown")
    fmt.Printf("Dereferenced name: %s\n", derefName)

    // JSON unmarshal
    jsonStr := `{"id":"123","name":"Bob","phone":"08123456789"}`
    user, err := helpers.UnmarshalTo[User](jsonStr)
    if err != nil {
        panic(err)
    }

    user.Phone = helpers.NormalizePhoneNumber(user.Phone)
    fmt.Printf("User: %+v\n", user)
}
```

## Testing

Jalankan unit tests dengan:

```bash
go test ./helpers
```

Tests mencakup:
- Pointer dan DerefPointer dengan berbagai types
- UnmarshalTo dengan valid dan invalid JSON
- GenerateTrxID dan variants (format dan uniqueness)
- GenerateMessageID dan GenerateUniqueID (format dan uniqueness)
- NormalizePhoneNumber dengan berbagai input cases

## Dependencies

- `github.com/google/uuid` (untuk GenerateMessageID dan GenerateUniqueID)
- Standard library: `strings`, `math/rand`, `time`, `encoding/json`