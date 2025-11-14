# Helpers Package

The `helpers` package provides a collection of utility functions for common operations in Go applications, including JSON manipulation, pointer operations, and string utilities.

## Features

- 🔤 JSON string manipulation and validation
- 👉 Safe pointer operations for primitive types
- 🔧 Common string utilities and ID generation
- ⚡ Type-safe and efficient implementations
- 🧪 Well-tested and production-ready

## Installation

This package is part of `github.com/budimanlai/go-pkg`. Import it as:

```go
import "github.com/budimanlai/go-pkg/helpers"
```

## API Reference

### JSON Functions

#### UnmarshalTo
```go
func UnmarshalTo[T any](jsonString string) (T, error)
```
Deserializes a JSON string into a value of type T using generics.

**Example:**
```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

jsonStr := `{"name":"John","email":"john@example.com"}`
user, err := helpers.UnmarshalTo[User](jsonStr)
if err != nil {
    log.Fatal(err)
}
fmt.Println(user.Name) // Output: John
```

#### UnmarshalFromMap
```go
func UnmarshalFromMap[T any](dataMap map[string]interface{}) (T, error)
```
Deserializes a map[string]interface{} into a value of type T using generics.

**Example:**
```go
type User struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
}

dataMap := map[string]interface{}{
    "name": "John",
    "age":  30,
}

user, err := helpers.UnmarshalFromMap[User](dataMap)
if err != nil {
    log.Fatal(err)
}
fmt.Println(user.Name) // Output: John
```
```

### Pointer Functions

#### StringPtr / IntPtr / BoolPtr / Float64Ptr
```go
func StringPtr(s string) *string
func IntPtr(i int) *int
func BoolPtr(b bool) *bool
func Float64Ptr(f float64) *float64
```
Returns a pointer to the given value.

**Example:**
```go
name := helpers.StringPtr("John")
age := helpers.IntPtr(30)
active := helpers.BoolPtr(true)
price := helpers.Float64Ptr(99.99)
```

#### StringValue / IntValue / BoolValue / Float64Value
```go
func StringValue(s *string) string
func IntValue(i *int) int
func BoolValue(b *bool) bool
func Float64Value(f *float64) float64
```
Safely dereferences a pointer, returning zero value if nil.

**Example:**
```go
var name *string
fmt.Println(helpers.StringValue(name)) // Output: "" (safe, no panic)

name = helpers.StringPtr("John")
fmt.Println(helpers.StringValue(name)) // Output: John
```

### String & ID Generation Functions

#### GenerateTrxID
```go
func GenerateTrxID() string
```
Generates a unique transaction ID based on timestamp and random number.

**Format:** YYMMDDHHMMSS + 4 random digits (16 characters total)

**Example:**
```go
id := helpers.GenerateTrxID()
// Output: "2511150430521234"
```

#### GenerateTrxIDWithPrefix
```go
func GenerateTrxIDWithPrefix(prefix string) string
```
Generates transaction ID with a custom prefix.

**Example:**
```go
id := helpers.GenerateTrxIDWithPrefix("TXN")
// Output: "TXN2511150430521234"
```

#### GenerateTrxIDWithSuffix
```go
func GenerateTrxIDWithSuffix(suffix string) string
```
Generates transaction ID with a custom suffix.

**Example:**
```go
id := helpers.GenerateTrxIDWithSuffix("END")
// Output: "2511150430521234END"
```

#### GenerateMessageID
```go
func GenerateMessageID() string
```
Generates a UUID v4 for message identification.

**Example:**
```go
id := helpers.GenerateMessageID()
// Output: "550e8400-e29b-41d4-a716-446655440000"
```

#### GenerateUniqueID
```go
func GenerateUniqueID() string
```
Generates a short unique ID (first 8 characters of UUID).

**Example:**
```go
id := helpers.GenerateUniqueID()
// Output: "550e8400"
```

#### NormalizePhoneNumber
```go
func NormalizePhoneNumber(phone string) string
```
Normalizes phone numbers to international format without the + prefix.

**Rules:**
- Removes + prefix if present
- Converts leading 0 to 62 (Indonesia)
- Adds 62 prefix if no country code detected
- Supports country codes: 62 (Indonesia), 1 (US), 65 (Singapore)

**Example:**
```go
// Indonesia
helpers.NormalizePhoneNumber("+628123456789") // "628123456789"
helpers.NormalizePhoneNumber("08123456789")   // "628123456789"
helpers.NormalizePhoneNumber("8123456789")    // "628123456789"

// Singapore
helpers.NormalizePhoneNumber("+658123456789") // "658123456789"

// US
helpers.NormalizePhoneNumber("+18123456789")  // "18123456789"
```

## Usage Examples

### Working with JSON

```go
type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

// Convert to JSON
product := Product{ID: 1, Name: "Laptop", Price: 999.99}
jsonStr := helpers.ToJSON(product)

// Validate JSON
if helpers.IsJSON(jsonStr) {
    fmt.Println("Valid JSON")
}

// Parse JSON
var newProduct Product
if err := helpers.FromJSON(jsonStr, &newProduct); err != nil {
    log.Fatal(err)
}
```

### Optional Fields with Pointers

```go
type User struct {
    Name     string
    Email    string
    Age      *int
    Bio      *string
    Verified *bool
}

user := User{
    Name:     "John Doe",
    Email:    "john@example.com",
    Age:      helpers.IntPtr(30),
    Bio:      helpers.StringPtr("Software developer"),
    Verified: helpers.BoolPtr(true),
}

// Safely access optional fields
fmt.Println("User age:", helpers.IntValue(user.Age))
fmt.Println("User verified:", helpers.BoolValue(user.Verified))
```

### Transaction ID Generation

```go
// Simple transaction ID
trxID := helpers.GenerateTrxID()
fmt.Println("Transaction ID:", trxID)

// With prefix
paymentID := helpers.GenerateTrxIDWithPrefix("PAY")
fmt.Println("Payment ID:", paymentID)

// With suffix
orderID := helpers.GenerateTrxIDWithSuffix("ORD")
fmt.Println("Order ID:", orderID)

// UUID for messaging
msgID := helpers.GenerateMessageID()
fmt.Println("Message ID:", msgID)
```

## Best Practices

1. **JSON Validation**: Always use `IsJSON()` before parsing untrusted JSON strings
2. **Nil Safety**: Use pointer value functions when dereferencing potentially nil pointers
3. **Optional Fields**: Use pointer types for optional struct fields
4. **Error Handling**: Check errors returned by `FromJSON()`
5. **Unique IDs**: Use appropriate ID generator based on your use case

## Testing

Run tests with:
```bash
go test ./helpers/ -v
```

## License

This package is part of the go-pkg project and follows the same license.
