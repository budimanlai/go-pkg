# go-pkg

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/budimanlai/go-pkg)](https://goreportcard.com/report/github.com/budimanlai/go-pkg)
[![GoDoc](https://godoc.org/github.com/budimanlai/go-pkg?status.svg)](https://godoc.org/github.com/budimanlai/go-pkg)

A comprehensive Go utility package providing essential tools for web development, including internationalization (i18n), HTTP response handling, custom types, and helper functions.

## Features

- **i18n**: Multi-language support with JSON-based locales and Fiber integration
- **Response**: Standardized HTTP response utilities with i18n support
- **Validator**: Struct validation with user-friendly, multilingual error messages
- **Security**: Password hashing and verification with bcrypt
- **Types**: Custom time types (UTCTime) for consistent UTC JSON serialization
- **Helpers**: Utility functions for pointers, JSON handling, string manipulation, and ID generation
- **Databases**: MySQL and PostgreSQL database utilities with GORM integration
- **Logger**: Logging utilities with timestamp support
- **Storage**: File storage abstraction supporting local filesystem and AWS S3
- **Middleware**: Authentication middleware for Fiber (Basic Auth, JWT, API Key, etc.)

## Installation

```bash
go get github.com/budimanlai/go-pkg
```

## Quick Start

### Basic Usage

```go
import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/response"
    "github.com/budimanlai/go-pkg/types"
    "github.com/gofiber/fiber/v2"
    "golang.org/x/text/language"
)

// Setup i18n
config := i18n.I18nConfig{
    DefaultLanguage: language.English,
    SupportedLangs:  []string{"en", "id"},
    LocalesPath:     "./locales",
}
i18nManager, _ := i18n.NewI18nManager(config)

// Use in Fiber app
app := fiber.New()
response.SetI18nManager(i18nManager)

app.Get("/", func(c *fiber.Ctx) error {
    return response.Success(c, "Welcome!", map[string]string{"version": "1.0"})
})
```

### Custom Time Type

```go
import "github.com/budimanlai/go-pkg/types"

type Event struct {
    CreatedAt types.UTCTime `json:"created_at"`
}

event := Event{CreatedAt: types.UTCTime(time.Now())}
// JSON output: {"created_at":"2025-10-15T12:30:45Z"}
```

## Documentation

Comprehensive documentation is available in the [docs/](docs/) folder:

### Main Packages

- **[databases](docs/databases.md)** - MySQL and PostgreSQL database management with GORM
- **[helpers](docs/helpers.md)** - JSON utilities, pointer operations, string helpers, ID generation
- **[i18n](docs/i18n.md)** - Internationalization with go-i18n and Fiber middleware
- **[logger](docs/logger.md)** - Logging utilities with timestamp support
- **[security](docs/security.md)** - Password hashing and verification with bcrypt
- **[types](docs/types.md)** - Custom UTCTime type for timezone-safe JSON handling
- **[storage](docs/storage.md)** - File storage abstraction for local filesystem and AWS S3
- **[middleware](docs/middleware.md)** - Authentication middleware for Fiber applications

### Middleware Package

The middleware package provides authentication solutions for Fiber applications:

- **[Middleware Overview](docs/middleware.md)** - Complete guide for all authentication methods
- **[JWT Auth](docs/jwt-auth.md)** - JSON Web Token authentication with multiple token sources
- **[Header Auth](docs/header-auth.md)** - Header-based API key authentication (X-API-Key)
- Basic Auth - HTTP Basic Authentication (username/password)
- QueryString Auth - Query parameter-based API key authentication
- Database API Key - Database-backed API key management with GORM

### Response Package (Multiple Files)

The response package documentation is organized in the [docs/response/](docs/response/) folder:

- **[README](docs/response/README.md)** - Overview and quick start
- **[Standard Responses](docs/response/standard-responses.md)** - Basic response functions (Success, Error, NotFound, BadRequest)
- **[I18n Responses](docs/response/i18n-responses.md)** - Internationalized response functions with automatic translation
- **[Error Handler](docs/response/error-handler.md)** - Custom Fiber error handler with i18n support
- **[Examples](docs/response/examples.md)** - Practical examples (REST API, authentication, file upload, pagination)

### Validator Package (Multiple Files)

The validator package documentation is organized in the [docs/validator/](docs/validator/) folder:

- **[README](docs/validator/README.md)** - Overview and quick start
- **[Validation Tags](docs/validator/validation-tags.md)** - Complete reference for all validation rules (required, email, min, max, etc.)
- **[Error Handling](docs/validator/error-handling.md)** - ValidationError type, error handling patterns, response formats
- **[I18n Integration](docs/validator/i18n-integration.md)** - Multilingual validation messages setup and configuration
- **[Examples](docs/validator/examples.md)** - Practical examples (user registration, complex structs, nested validation)

### Storage Package

Abstraction layer for file storage operations supporting multiple backends:

```go
import "github.com/budimanlai/go-pkg/storage"

// Local storage
localStorage, _ := storage.NewLocalStorage(storage.LocalStorageConfig{
    BasePath: "./uploads",
    BaseURL:  "http://localhost:3000/uploads",
})

// S3 storage
s3Storage, _ := storage.NewS3Storage(storage.S3StorageConfig{
    Region:      "us-east-1",
    Bucket:      "my-bucket",
    Endpoint:    "", // Optional for AWS S3
    AccessKeyID: "your-access-key",
    SecretKey:   "your-secret-key",
})

// Use storage interface
file, _ := os.Open("image.jpg")
url, _ := localStorage.Put(context.Background(), "images/photo.jpg", file, nil)
```

**Features:**
- Unified interface for local and S3 storage
- File operations: Put, Get, Delete, Exists, GetURL
- Stream support for large files
- Automatic directory creation (local)
- Public/private file support

### Middleware Package

Authentication middleware for Fiber applications with flexible configuration:

```go
import "github.com/budimanlai/go-pkg/middleware/auth"

// JWT Authentication
jwtAuth := auth.NewJWTAuth(auth.JWTConfig{
    SecretKey: "your-secret-key",
})
app.Use(jwtAuth.Middleware())

// API Key (Header-based)
headerAuth := auth.NewHeaderAuth(auth.HeaderAuthConfig{
    KeyProvider: keyProvider,
    HeaderName:  "X-API-Key",
})
app.Use(headerAuth.Middleware())

// Basic Auth
basicAuth := auth.NewBasicAuth(auth.BasicAuthConfig{
    KeyProvider: keyProvider,
})
app.Use(basicAuth.Middleware())
```

**Features:**
- JWT token authentication (header/query/cookie)
- Header-based API key authentication
- HTTP Basic Authentication
- Query string authentication
- Database-backed API key storage
- In-memory key provider
- Thread-safe operations
- Custom success/error handlers
- Full test coverage (87 tests)

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Generate coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Project Structure

```
go-pkg/
├── databases/          # Database utilities (MySQL, PostgreSQL)
├── docs/              # Documentation
├── helpers/           # General utility functions
├── i18n/              # Internationalization
├── locales/           # Translation files
├── logger/            # Logging utilities
├── middleware/        # Authentication middleware
│   └── auth/          # Auth implementations (JWT, Basic, Header, etc.)
├── response/          # HTTP response helpers
├── security/          # Password hashing utilities
├── storage/           # File storage abstraction (Local, S3)
├── types/             # Custom types
├── validator/         # Validation utilities
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please ensure all tests pass and add tests for new features.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you find this package helpful, please give it a ⭐ on GitHub!

For issues or questions, please open an issue on GitHub.