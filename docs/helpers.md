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

### Pointer Functions

#### Pointer
```go
func Pointer[T any](v T) *T
```
Returns a pointer to the given value of any type T. This is a generic helper function useful for creating pointers to literals or values where taking the address directly is not possible.

**Example:**
```go
// String pointer
name := helpers.Pointer("John")        // *string
fmt.Println(*name)                     // Output: John

// Integer pointer
age := helpers.Pointer(30)             // *int
fmt.Println(*age)                      // Output: 30

// Boolean pointer
active := helpers.Pointer(true)        // *bool
fmt.Println(*active)                   // Output: true

// Float pointer
price := helpers.Pointer(99.99)        // *float64
fmt.Println(*price)                    // Output: 99.99

// Struct pointer
type User struct {
    Name string
    Age  int
}
user := helpers.Pointer(User{Name: "Alice", Age: 25})
fmt.Println(user.Name)                 // Output: Alice
```

#### DerefPointer
```go
func DerefPointer[T any](p *T, defaultValue T) T
```
Safely dereferences a pointer and returns its value. If the pointer is nil, it returns the provided defaultValue instead.

**Example:**
```go
// With nil pointer - returns default value
var name *string
result := helpers.DerefPointer(name, "default")
fmt.Println(result) // Output: "default" (safe, no panic)

// With non-nil pointer - returns actual value
name = helpers.Pointer("John")
result = helpers.DerefPointer(name, "default")
fmt.Println(result) // Output: "John"

// With numbers
var age *int
result := helpers.DerefPointer(age, 18)
fmt.Println(result) // Output: 18

num := helpers.Pointer(30)
result = helpers.DerefPointer(num, 18)
fmt.Println(result) // Output: 30
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

#### GenerateRandomString
```go
func GenerateRandomString(length int) string
```
Generates a random alphanumeric string of the specified length.

**Parameters:**
- `length` (int): Desired length of the random string

**Returns:**
- Random string consisting of uppercase letters, lowercase letters, and digits

**Example:**
```go
randomStr := helpers.GenerateRandomString(10)
// Output: "aZ3bC9dE1F" (random each time)

// Use cases
apiKey := helpers.GenerateRandomString(32)
token := helpers.GenerateRandomString(16)
code := helpers.GenerateRandomString(6)
```

---

### Date Functions

#### StringToDate
```go
func StringToDate(dateStr string) (time.Time, error)
```
Converts a string to a time.Time object using the "YYYY-MM-DD" format.

**Parameters:**
- `dateStr` (string): Date string in "YYYY-MM-DD" format

**Returns:**
- `time.Time`: Parsed time object
- `error`: Error if parsing fails

**Example:**
```go
date, err := helpers.StringToDate("2024-11-13")
if err != nil {
    log.Fatal(err)
}
fmt.Println(date.Format("2006-01-02")) // Output: 2024-11-13

// Invalid format
_, err = helpers.StringToDate("13-11-2024")
if err != nil {
    fmt.Println("Invalid date format")
}
```

## Usage Examples

### Working with JSON

```go
type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

// Parse JSON string to struct
jsonStr := `{"id":1,"name":"Laptop","price":999.99}`
product, err := helpers.UnmarshalTo[Product](jsonStr)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Product: %s, Price: %.2f\n", product.Name, product.Price)

// Convert map to struct
dataMap := map[string]interface{}{
    "id":    2,
    "name":  "Mouse",
    "price": 29.99,
}
product2, err := helpers.UnmarshalFromMap[Product](dataMap)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Product: %s, Price: %.2f\n", product2.Name, product2.Price)
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

// Create user with optional fields
user := User{
    Name:     "John Doe",
    Email:    "john@example.com",
    Age:      helpers.Pointer(30),
    Bio:      helpers.Pointer("Software developer"),
    Verified: helpers.Pointer(true),
}

// Safely access optional fields with defaults
age := helpers.DerefPointer(user.Age, 0)
bio := helpers.DerefPointer(user.Bio, "No bio")
verified := helpers.DerefPointer(user.Verified, false)

fmt.Printf("Age: %d, Bio: %s, Verified: %t\n", age, bio, verified)

// Handling nil pointers
var optionalUser User
optionalUser.Name = "Jane"
optionalUser.Email = "jane@example.com"
// Age, Bio, Verified are nil

age = helpers.DerefPointer(optionalUser.Age, 18)       // Returns 18 (default)
bio = helpers.DerefPointer(optionalUser.Bio, "N/A")    // Returns "N/A" (default)
verified = helpers.DerefPointer(optionalUser.Verified, false) // Returns false (default)
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

1. **JSON Unmarshaling**: Always check errors returned by `UnmarshalTo()` and `UnmarshalFromMap()`
2. **Nil Safety**: Use `DerefPointer()` when dereferencing potentially nil pointers with safe defaults
3. **Optional Fields**: Use pointer types for optional struct fields and `Pointer()` to create them
4. **Phone Normalization**: Always normalize phone numbers before storage or validation
5. **Unique IDs**: Use appropriate ID generator based on your use case (transaction, message, or general unique ID)
6. **Date Parsing**: Use `StringToDate()` for consistent date parsing in "YYYY-MM-DD" format and always handle errors
7. **Random Strings**: Use `GenerateRandomString()` for tokens, API keys, or verification codes with appropriate length for security

## Testing

Run tests with:
```bash
go test ./helpers/ -v
```

## License

This package is part of the go-pkg project and follows the same license.
