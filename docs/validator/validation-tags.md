# Validation Tags

Comprehensive reference for all validation tags supported by the validator package. These tags are based on `go-playground/validator/v10`.

## String Validation

### required
Field cannot be empty or zero value.

```go
type User struct {
    Name string `validate:"required"`
}
```

**Valid:** `"John"`, `"A"`, `""`  
**Invalid:** `""` (empty string)

### email
Must be a valid email address format.

```go
type User struct {
    Email string `validate:"required,email"`
}
```

**Valid:** `john@example.com`, `user.name+tag@example.co.uk`  
**Invalid:** `invalid-email`, `@example.com`, `user@`

### url
Must be a valid URL.

```go
type Website struct {
    URL string `validate:"url"`
}
```

**Valid:** `https://example.com`, `http://localhost:8080`  
**Invalid:** `not-a-url`, `example.com` (missing protocol)

### uri
Must be a valid URI.

```go
type Resource struct {
    URI string `validate:"uri"`
}
```

**Valid:** `https://example.com`, `tel:+1-234-567-8900`, `mailto:user@example.com`  
**Invalid:** `not a uri`

### alpha
Only alphabetic characters (a-z, A-Z).

```go
type Name struct {
    FirstName string `validate:"alpha"`
}
```

**Valid:** `John`, `abc`, `XYZ`  
**Invalid:** `John123`, `abc-def`, `John Doe`

### alphanum
Only alphanumeric characters (a-z, A-Z, 0-9).

```go
type Username struct {
    Value string `validate:"alphanum"`
}
```

**Valid:** `John123`, `abc`, `XYZ999`  
**Invalid:** `John-123`, `abc_def`, `user@name`

### alphanumunicode
Alphanumeric Unicode characters.

```go
type Name struct {
    Value string `validate:"alphanumunicode"`
}
```

**Valid:** `John123`, `名前123`, `usuario123`  
**Invalid:** `John-123`, `abc_def`

### numeric
Only numeric characters (0-9).

```go
type Code struct {
    Value string `validate:"numeric"`
}
```

**Valid:** `123`, `0`, `999999`  
**Invalid:** `12.34`, `abc`, `12a`

### hexadecimal
Valid hexadecimal string.

```go
type Color struct {
    Hex string `validate:"hexadecimal"`
}
```

**Valid:** `FF0000`, `abc123`, `0x1234`  
**Invalid:** `GG0000`, `xyz`

### lowercase
All characters must be lowercase.

```go
type Tag struct {
    Value string `validate:"lowercase"`
}
```

**Valid:** `abc`, `hello`, `test123`  
**Invalid:** `ABC`, `Hello`, `Test`

### uppercase
All characters must be uppercase.

```go
type Code struct {
    Value string `validate:"uppercase"`
}
```

**Valid:** `ABC`, `HELLO`, `TEST123`  
**Invalid:** `abc`, `Hello`, `test`

## Length Validation

### min
Minimum length for strings, minimum value for numbers.

```go
type User struct {
    Password string `validate:"min=8"`
    Age      int    `validate:"min=18"`
}
```

**String Valid:** `"12345678"`, `"password123"`  
**String Invalid:** `"123"`, `"short"`

**Number Valid:** `18`, `25`, `100`  
**Number Invalid:** `17`, `0`, `10`

### max
Maximum length for strings, maximum value for numbers.

```go
type User struct {
    Name string `validate:"max=50"`
    Age  int    `validate:"max=100"`
}
```

**String Valid:** `"John"`, `"A"`, `"12345"`  
**String Invalid:** `"Very long string that exceeds 50 characters..."`

**Number Valid:** `100`, `50`, `0`  
**Number Invalid:** `101`, `200`

### len
Exact length for strings, exact value for numbers.

```go
type Code struct {
    Value string `validate:"len=10"`
}
```

**Valid:** `"1234567890"`  
**Invalid:** `"123"`, `"12345678901"`

## Number Validation

### gt
Greater than (>).

```go
type Product struct {
    Price float64 `validate:"gt=0"`
}
```

**Valid:** `0.01`, `10.5`, `999`  
**Invalid:** `0`, `-10`

### gte
Greater than or equal (≥).

```go
type User struct {
    Age int `validate:"gte=18"`
}
```

**Valid:** `18`, `25`, `100`  
**Invalid:** `17`, `0`, `-1`

### lt
Less than (<).

```go
type Percentage struct {
    Value int `validate:"lt=100"`
}
```

**Valid:** `99`, `50`, `0`  
**Invalid:** `100`, `101`

### lte
Less than or equal (≤).

```go
type Score struct {
    Value int `validate:"lte=100"`
}
```

**Valid:** `100`, `50`, `0`  
**Invalid:** `101`, `200`

### eq
Equal to (==).

```go
type Status struct {
    Code int `validate:"eq=200"`
}
```

**Valid:** `200`  
**Invalid:** `201`, `404`, `0`

### ne
Not equal to (!=).

```go
type Status struct {
    Code int `validate:"ne=0"`
}
```

**Valid:** `1`, `200`, `-1`  
**Invalid:** `0`

### oneof
Value must be one of the specified values.

```go
type Status struct {
    Value string `validate:"oneof=pending approved rejected"`
}
```

**Valid:** `"pending"`, `"approved"`, `"rejected"`  
**Invalid:** `"processing"`, `"cancelled"`

## Comparison Validation

### eqfield
Field value must equal another field.

```go
type User struct {
    Password        string `validate:"required,min=8"`
    PasswordConfirm string `validate:"required,eqfield=Password"`
}
```

**Valid:** Both fields have same value  
**Invalid:** Fields have different values

### nefield
Field value must not equal another field.

```go
type PasswordChange struct {
    OldPassword string `validate:"required"`
    NewPassword string `validate:"required,nefield=OldPassword"`
}
```

**Valid:** Fields have different values  
**Invalid:** Both fields have same value

### gtfield
Field value must be greater than another field.

```go
type DateRange struct {
    StartDate time.Time `validate:"required"`
    EndDate   time.Time `validate:"required,gtfield=StartDate"`
}
```

### ltefield
Field value must be less than or equal to another field.

```go
type Range struct {
    Min int `validate:"required"`
    Max int `validate:"required,gtefield=Min"`
}
```

## Format Validation

### datetime
Must be a valid datetime string with specified format.

```go
type Event struct {
    Date string `validate:"datetime=2006-01-02"`
    Time string `validate:"datetime=15:04:05"`
}
```

**Valid:** `"2024-01-15"`, `"14:30:00"`  
**Invalid:** `"15-01-2024"`, `"2:30 PM"`

### isbn
Valid ISBN (ISBN-10 or ISBN-13).

```go
type Book struct {
    ISBN string `validate:"isbn"`
}
```

**Valid:** `978-3-16-148410-0`, `0-306-40615-2`  
**Invalid:** `123-456`, `invalid`

### uuid
Valid UUID string.

```go
type Resource struct {
    ID string `validate:"uuid"`
}
```

**Valid:** `550e8400-e29b-41d4-a716-446655440000`  
**Invalid:** `not-a-uuid`, `12345`

### uuid4
Valid UUID version 4.

```go
type Resource struct {
    ID string `validate:"uuid4"`
}
```

### ipv4
Valid IPv4 address.

```go
type Server struct {
    IP string `validate:"ipv4"`
}
```

**Valid:** `192.168.1.1`, `127.0.0.1`  
**Invalid:** `999.999.999.999`, `192.168.1`

### ipv6
Valid IPv6 address.

```go
type Server struct {
    IP string `validate:"ipv6"`
}
```

**Valid:** `2001:0db8:85a3:0000:0000:8a2e:0370:7334`  
**Invalid:** `192.168.1.1`

### mac
Valid MAC address.

```go
type Device struct {
    MAC string `validate:"mac"`
}
```

**Valid:** `01:23:45:67:89:ab`, `01-23-45-67-89-AB`  
**Invalid:** `invalid-mac`

### latitude
Valid latitude (-90 to 90).

```go
type Location struct {
    Lat string `validate:"latitude"`
}
```

**Valid:** `"-90.0"`, `"45.5"`, `"90.0"`  
**Invalid:** `"91.0"`, `"-91.0"`

### longitude
Valid longitude (-180 to 180).

```go
type Location struct {
    Lng string `validate:"longitude"`
}
```

**Valid:** `"-180.0"`, `"90.5"`, `"180.0"`  
**Invalid:** `"181.0"`, `"-181.0"`

## Collection Validation

### dive
Validates each element in a slice/array/map.

```go
type UserList struct {
    Emails []string `validate:"required,dive,email"`
}
```

### unique
All elements must be unique.

```go
type Tags struct {
    Values []string `validate:"unique"`
}
```

**Valid:** `[]string{"a", "b", "c"}`  
**Invalid:** `[]string{"a", "b", "a"}`

### contains
String must contain substring.

```go
type URL struct {
    Value string `validate:"contains=example.com"`
}
```

**Valid:** `"https://example.com/path"`  
**Invalid:** `"https://other.com"`

### excludes
String must not contain substring.

```go
type Username struct {
    Value string `validate:"excludes=admin"`
}
```

**Valid:** `"john_doe"`, `"user123"`  
**Invalid:** `"admin"`, `"superadmin"`

### startswith
String must start with prefix.

```go
type Code struct {
    Value string `validate:"startswith=PRE"`
}
```

**Valid:** `"PRE-123"`, `"PREFIX"`  
**Invalid:** `"123-PRE"`, `"code"`

### endswith
String must end with suffix.

```go
type File struct {
    Name string `validate:"endswith=.pdf"`
}
```

**Valid:** `"document.pdf"`, `"file.pdf"`  
**Invalid:** `"document.doc"`, `"pdf"`

## Special Validation

### omitempty
Skip validation if field is empty.

```go
type User struct {
    Email string `validate:"omitempty,email"`
}
```

Field is validated only if it has a value.

### required_if
Required if another field has a specific value.

```go
type Payment struct {
    Method      string `validate:"required,oneof=card paypal"`
    CardNumber  string `validate:"required_if=Method card"`
}
```

### required_without
Required if another field is empty.

```go
type Contact struct {
    Email string `validate:"required_without=Phone"`
    Phone string `validate:"required_without=Email"`
}
```

At least one field must be filled.

## Custom Validation Messages

All validation tags have default messages that can be translated via i18n:

```json
{
  "validator.required": "{{.FieldName}} is required",
  "validator.email": "{{.FieldName}} must be a valid email address",
  "validator.min": "{{.FieldName}} must be at least {{.Param}} characters",
  "validator.max": "{{.FieldName}} must be at most {{.Param}} characters"
}
```

See [I18n Integration](i18n-integration.md) for more details.

## Combining Multiple Tags

Use comma to combine multiple validation rules:

```go
type User struct {
    Email    string `validate:"required,email,max=100"`
    Password string `validate:"required,min=8,max=64,alphanum"`
    Age      int    `validate:"required,gte=18,lte=100"`
    Website  string `validate:"omitempty,url"`
}
```

## Related Documentation

- [README](README.md) - Package overview
- [Error Handling](error-handling.md) - ValidationError type and patterns
- [I18n Integration](i18n-integration.md) - Multilingual validation messages
- [Examples](examples.md) - Practical usage examples
