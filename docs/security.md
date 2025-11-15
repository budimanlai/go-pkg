# Security Package

The `security` package provides secure password hashing and verification utilities using bcrypt.

## Overview

This package offers simple and secure password handling functions that use the bcrypt algorithm, which is specifically designed for password hashing with built-in salt generation and configurable cost factor.

## Installation

```bash
go get golang.org/x/crypto/bcrypt
```

## Functions

### HashPassword

Hashes a password using bcrypt with the default cost factor.

**Signature:**
```go
func HashPassword(password string) string
```

**Parameters:**
- `password` (string): The plain text password to hash

**Returns:**
- (string): The bcrypt hash of the password, or empty string if hashing fails

**Example:**
```go
import "github.com/budimanlai/go-pkg/security"

password := "mySecurePassword123"
hash := security.HashPassword(password)
fmt.Println("Hashed password:", hash)
// Output: $2a$10$... (60-character bcrypt hash)
```

**Notes:**
- Uses `bcrypt.DefaultCost` (currently 10) as the cost factor
- Returns empty string if an error occurs during hashing
- Each call generates a unique hash due to random salt

---

### CheckPasswordHash

Compares a plain text password with a bcrypt hash to verify if they match.

**Signature:**
```go
func CheckPasswordHash(password, hash string) (bool, error)
```

**Parameters:**
- `password` (string): The plain text password to verify
- `hash` (string): The bcrypt hash to compare against

**Returns:**
- (bool): `true` if password matches the hash, `false` otherwise
- (error): Error if something went wrong (not including mismatch)

**Example:**
```go
import "github.com/budimanlai/go-pkg/security"

password := "mySecurePassword123"
hash := security.HashPassword(password)

// Correct password
valid, err := security.CheckPasswordHash(password, hash)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Password valid:", valid) // true

// Wrong password
valid, err = security.CheckPasswordHash("wrongPassword", hash)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Password valid:", valid) // false
```

**Error Handling:**
- Returns `(false, nil)` if the password doesn't match (expected case)
- Returns `(false, error)` only for unexpected errors (corrupted hash, etc.)

## Usage Examples

### User Registration

```go
type User struct {
    Email        string
    PasswordHash string
}

func RegisterUser(email, password string) (*User, error) {
    // Hash the password
    hash := security.HashPassword(password)
    if hash == "" {
        return nil, errors.New("failed to hash password")
    }
    
    user := &User{
        Email:        email,
        PasswordHash: hash,
    }
    
    // Save user to database...
    
    return user, nil
}
```

### User Login

```go
func LoginUser(email, password string) (*User, error) {
    // Fetch user from database...
    user, err := findUserByEmail(email)
    if err != nil {
        return nil, err
    }
    
    // Verify password
    valid, err := security.CheckPasswordHash(password, user.PasswordHash)
    if err != nil {
        return nil, err
    }
    
    if !valid {
        return nil, errors.New("invalid credentials")
    }
    
    return user, nil
}
```

### Password Change

```go
func ChangePassword(userID int, oldPassword, newPassword string) error {
    // Fetch user from database...
    user, err := findUserByID(userID)
    if err != nil {
        return err
    }
    
    // Verify old password
    valid, err := security.CheckPasswordHash(oldPassword, user.PasswordHash)
    if err != nil {
        return err
    }
    
    if !valid {
        return errors.New("invalid old password")
    }
    
    // Hash new password
    newHash := security.HashPassword(newPassword)
    if newHash == "" {
        return errors.New("failed to hash new password")
    }
    
    // Update user password hash in database...
    user.PasswordHash = newHash
    
    return nil
}
```

### Password Reset with Token

```go
func ResetPassword(token, newPassword string) error {
    // Verify reset token and get user...
    user, err := findUserByResetToken(token)
    if err != nil {
        return err
    }
    
    // Hash new password
    hash := security.HashPassword(newPassword)
    if hash == "" {
        return errors.New("failed to hash password")
    }
    
    // Update password and invalidate token
    user.PasswordHash = hash
    user.ResetToken = ""
    user.ResetTokenExpiry = nil
    
    // Save to database...
    
    return nil
}
```

## Security Considerations

1. **Cost Factor**: The package uses `bcrypt.DefaultCost` (10), which provides a good balance between security and performance. Higher values increase security but also increase computation time.

2. **Salt**: Bcrypt automatically generates a random salt for each password, so identical passwords will have different hashes.

3. **Hash Storage**: Always store the full bcrypt hash (60 characters). Never truncate it.

4. **Timing Attacks**: The `CheckPasswordHash` function is resistant to timing attacks as it uses bcrypt's constant-time comparison.

5. **Password Policy**: This package only handles hashing. Implement password strength requirements in your application layer.

6. **Rate Limiting**: Implement rate limiting on login endpoints to prevent brute force attacks.

## Best Practices

1. **Never Log Passwords**: Never log plain text passwords or hashes
2. **Error Handling**: Check if `HashPassword()` returns empty string and handle appropriately
3. **Database Field**: Use `VARCHAR(60)` or `TEXT` for storing bcrypt hashes
4. **Password Requirements**: Enforce minimum password length (8+ characters recommended)
5. **Use HTTPS**: Always transmit passwords over secure connections
6. **Constant-Time Comparison**: Always use `CheckPasswordHash()` for verification, never compare hashes directly

## Testing

Run tests with:
```bash
go test ./security/ -v
```

## Dependencies

- `golang.org/x/crypto/bcrypt`: Industry-standard bcrypt implementation

## License

This package is part of the go-pkg project and follows the same license.
