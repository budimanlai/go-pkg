# Logger Package

The `logger` package provides simple logging utilities for Go applications with timestamp support and flexible output control.

## Features

- 📝 Formatted logging with timestamps
- 🐛 Debug logging with toggle control
- 🔍 Variable dump with JSON formatting
- 📊 Hexadecimal data output
- ⚙️ Global output control flags
- 🎨 Clean and readable output format

## Installation

This package is part of `github.com/budimanlai/go-pkg`. Import it as:

```go
import "github.com/budimanlai/go-pkg/logger"
```

## Configuration

### Global Flags

```go
var (
    ShowOutput = true  // Controls Printf, PrintHex, and Vardump
    ShowDebug  = true  // Controls Debugf output
)
```

You can disable logging globally:

```go
logger.ShowOutput = false  // Disable regular logging
logger.ShowDebug = false   // Disable debug logging
```

## API Reference

### Vardump
```go
func Vardump(v any)
```
Prints a formatted JSON representation of any value.

**Example:**
```go
type User struct {
    Name string
    Age  int
}
user := User{Name: "John", Age: 30}
logger.Vardump(user)
// Output:
// {
//   "Name": "John",
//   "Age": 30
// }
```

### Printf
```go
func Printf(format string, args ...interface{})
```
Formats and prints a log message with timestamp prefix.

**Output Format:** `[YYYY-MM-DD HH:MM:SS] message`

**Example:**
```go
logger.Printf("User %s logged in", "john@example.com")
// Output: [2025-11-15 04:56:56] User john@example.com logged in

logger.Printf("Processing order #%d with amount $%.2f", 1234, 99.99)
// Output: [2025-11-15 04:56:56] Processing order #1234 with amount $99.99
```

### PrintHex
```go
func PrintHex(data []byte)
```
Prints hexadecimal representation of byte data with timestamp.

**Example:**
```go
data := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
logger.PrintHex(data)
// Output: [2025-11-15 04:56:56] 48656c6c6f
```

### Debugf
```go
func Debugf(format string, args ...interface{})
```
Formats and prints a debug message with timestamp if ShowDebug is enabled.

**Output Format:** `[YYYY-MM-DD HH:MM:SS] DEBUG: message`

**Example:**
```go
logger.Debugf("Variable value: %v", someVar)
// Output: [2025-11-15 04:56:56] DEBUG: Variable value: 42
```

### Error
```go
func Error(msg string)
```
Logs an error message with timestamp and ERROR prefix.

**Output Format:** `[YYYY-MM-DD HH:MM:SS] ERROR: message`

**Example:**
```go
logger.Error("Failed to connect to database")
// Output: [2025-11-15 04:56:56] ERROR: Failed to connect to database
```

### Errorf
```go
func Errorf(format string, args ...interface{})
```
Formats and logs an error message with timestamp.

**Example:**
```go
logger.Errorf("Database error: %v", err)
// Output: [2025-11-15 04:56:56] ERROR: Database error: connection timeout
```

### Info
```go
func Info(msg string)
```
Logs an informational message with timestamp and INFO prefix.

**Example:**
```go
logger.Info("Server started successfully")
// Output: [2025-11-15 04:56:56] INFO: Server started successfully
```

### Infof
```go
func Infof(format string, args ...interface{})
```
Formats and logs an informational message.

**Example:**
```go
logger.Infof("Server listening on port %d", 8080)
// Output: [2025-11-15 04:56:56] INFO: Server listening on port 8080
```

### Fatal
```go
func Fatal(msg string)
```
Logs a fatal error message and terminates the program with exit code 1.

**Example:**
```go
logger.Fatal("Critical error: Unable to start server")
// Output: [2025-11-15 04:56:56] FATAL: Critical error: Unable to start server
// Program exits with code 1
```

### Fatalf
```go
func Fatalf(format string, args ...interface{})
```
Formats a fatal error message and terminates the program.

**Example:**
```go
logger.Fatalf("Failed to load config: %v", err)
// Output: [2025-11-15 04:56:56] FATAL: Failed to load config: file not found
// Program exits with code 1
```

## Usage Examples

### Basic Logging

```go
package main

import "github.com/budimanlai/go-pkg/logger"

func main() {
    logger.Printf("Application started")
    logger.Printf("User %s logged in from IP %s", "john", "192.168.1.1")
    
    logger.Debugf("Processing request with ID: %s", requestID)
    
    logger.Info("Database connection established")
    logger.Infof("Connected to %s database", dbType)
    
    logger.Error("Failed to process payment")
    logger.Errorf("Payment failed: %v", err)
}
```

### Variable Debugging

```go
type Config struct {
    Host     string
    Port     int
    Database string
    Debug    bool
}

config := Config{
    Host:     "localhost",
    Port:     5432,
    Database: "myapp",
    Debug:    true,
}

logger.Vardump(config)
// Output:
// {
//   "Host": "localhost",
//   "Port": 5432,
//   "Database": "myapp",
//   "Debug": true
// }
```

### Conditional Logging

```go
// Disable debug logging in production
if os.Getenv("ENV") == "production" {
    logger.ShowDebug = false
}

logger.Printf("Application started")
logger.Debugf("Debug mode: %v", logger.ShowDebug) // Won't show in production

// Disable all logging for tests
if testing.Testing() {
    logger.ShowOutput = false
    logger.ShowDebug = false
}
```

### HTTP Request Logging

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        logger.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
        logger.Debugf("Headers: %v", r.Header)
        
        next.ServeHTTP(w, r)
        
        duration := time.Since(start)
        logger.Printf("Request completed in %v", duration)
    })
}
```

### Error Handling with Logging

```go
func connectDatabase(dsn string) (*sql.DB, error) {
    logger.Infof("Connecting to database...")
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        logger.Errorf("Database connection failed: %v", err)
        return nil, err
    }
    
    if err := db.Ping(); err != nil {
        logger.Errorf("Database ping failed: %v", err)
        return nil, err
    }
    
    logger.Info("Database connection established")
    return db, nil
}

func main() {
    db, err := connectDatabase("postgres://localhost/mydb")
    if err != nil {
        logger.Fatalf("Unable to start application: %v", err)
    }
    defer db.Close()
}
```

## Best Practices

1. **Production Logging**: Disable debug logging in production with `logger.ShowDebug = false`
2. **Structured Data**: Use `Vardump()` for complex objects that need inspection
3. **Error Context**: Include error details with `Errorf()` for better debugging
4. **Performance**: Disable logging in tests with `logger.ShowOutput = false`
5. **Fatal Errors**: Use `Fatal()` only for unrecoverable errors that should terminate the application
6. **Consistent Format**: Use formatted logging functions for consistent timestamp and prefix formatting

## Testing

Run tests with:
```bash
go test ./logger/ -v
```

## License

This package is part of the go-pkg project and follows the same license.
