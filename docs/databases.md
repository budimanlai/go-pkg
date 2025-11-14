# Databases Package

The `databases` package provides database connection management using GORM, supporting both MySQL and PostgreSQL databases with a simple and flexible API.

## Features

- 🔌 Database connection management with GORM
- 🗄️ Support for MySQL and PostgreSQL drivers
- ⚙️ Flexible database configuration
- 🔒 Connection pooling and lifecycle management
- 🧪 Easy testing with mock databases

## Installation

This package is part of `github.com/budimanlai/go-pkg`. Import it as:

```go
import "github.com/budimanlai/go-pkg/databases"
```

## Quick Start

### 1. Configure Database

```go
config := databases.DbConfig{
    Driver:          databases.MySQL, // or databases.Postgres
    Host:            "localhost",
    Port:            "3306",
    Username:        "root",
    Password:        "password",
    Name:            "myapp_db",
    Charset:         "utf8mb4", // optional, default: utf8mb4
    MaxIdleConns:    10,
    MaxOpenConns:    100,
    ConnMaxLifeTime: time.Hour,
}
```

### 2. Create Database Manager

```go
dbManager := databases.NewDbManager(config)
```

### 3. Open Connection

```go
// Open with default configuration
err := dbManager.Open()
if err != nil {
    log.Fatal(err)
}
defer dbManager.Close()
```

Or with custom GORM configuration:
```go
err := dbManager.OpenWithConfig(&gorm.Config{
    NowFunc: func() time.Time {
        return time.Now().UTC()
    },
    PrepareStmt: true,
})
```

### 4. Use Database

```go
db := dbManager.GetDb()

// Example: Auto migrate
db.AutoMigrate(&User{})

// Example: Create record
user := User{Name: "John", Email: "john@example.com"}
db.Create(&user)

// Example: Query
var users []User
db.Find(&users)
```

## API Reference

### Types

#### DatabaseDriver
```go
type DatabaseDriver string

const (
    MySQL    DatabaseDriver = "mysql"
    Postgres DatabaseDriver = "postgres"
)
```

#### DbConfig
```go
type DbConfig struct {
    Driver          DatabaseDriver // Database driver (MySQL or Postgres)
    Host            string         // Database host
    Port            string         // Database port
    Username        string         // Database username
    Password        string         // Database password
    Name            string         // Database name
    Charset         string         // Character set (optional, default: utf8mb4)
    MaxIdleConns    int            // Maximum idle connections (0 or negative to skip)
    MaxOpenConns    int            // Maximum open connections (0 or negative to skip)
    ConnMaxLifeTime time.Duration  // Connection max lifetime (0 or negative to skip)
}
```

#### DbManager
```go
type DbManager struct {
    Db     *gorm.DB
    Config DbConfig
}
```

### Functions

#### NewDbManager
```go
func NewDbManager(config DbConfig) *DbManager
```
Creates a new database manager instance for both MySQL and PostgreSQL.

**Example:**
```go
// MySQL
dbManager := databases.NewDbManager(databases.DbConfig{
    Driver:          databases.MySQL,
    Host:            "localhost",
    Port:            "3306",
    Username:        "root",
    Password:        "secret",
    Name:            "mydb",
    MaxIdleConns:    10,
    MaxOpenConns:    100,
    ConnMaxLifeTime: time.Hour,
})

// PostgreSQL
dbManager := databases.NewDbManager(databases.DbConfig{
    Driver:   databases.Postgres,
    Host:     "localhost",
    Port:     "5432",
    Username: "postgres",
    Password: "secret",
    Name:     "mydb",
})
```

## Usage Examples

### Basic CRUD Operations

```go
package main

import (
    "log"
    "time"
    "github.com/budimanlai/go-pkg/databases"
)

type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string
    Email string `gorm:"uniqueIndex"`
}

func main() {
    config := databases.DbConfig{
        Driver:          databases.MySQL,
        Host:            "localhost",
        Port:            "3306",
        Username:        "root",
        Password:        "password",
        Name:            "testdb",
        MaxIdleConns:    10,
        MaxOpenConns:    100,
        ConnMaxLifeTime: time.Hour,
    }

    dbManager := databases.NewDbManager(config)
    if err := dbManager.Open(); err != nil {
        log.Fatal(err)
    }
    defer dbManager.Close()

    db := dbManager.GetDb()
    db.AutoMigrate(&User{})

    // Create
    user := User{Name: "Alice", Email: "alice@example.com"}
    db.Create(&user)

    // Read
    var fetchedUser User
    db.First(&fetchedUser, user.ID)

    // Update
    db.Model(&fetchedUser).Update("Name", "Alice Smith")

    // Delete
    db.Delete(&fetchedUser)
}
```

### PostgreSQL Example

```go
config := databases.DbConfig{
    Driver:   databases.Postgres,
    Host:     "localhost",
    Port:     "5432",
    Username: "postgres",
    Password: "password",
    Name:     "testdb",
}

dbManager := databases.NewDbManager(config)
if err := dbManager.Open(); err != nil {
    log.Fatal(err)
}
defer dbManager.Close()

db := dbManager.GetDb()
// Use db for operations
```

### Custom GORM Configuration

```go
dbManager := databases.NewDbManager(config)
err := dbManager.OpenWithConfig(&gorm.Config{
    PrepareStmt: true,
    NowFunc: func() time.Time {
        return time.Now().UTC()
    },
})
if err != nil {
    log.Fatal(err)
}
defer dbManager.Close()
```

## Best Practices

1. **Always Close Connections**: Use `defer dbManager.Close()` to ensure connections are properly closed
2. **Error Handling**: Always check errors returned by `Open()` and `OpenWithConfig()`
3. **Connection Pooling**: Configure connection pool settings for production use
4. **Environment Variables**: Store database credentials in environment variables, not in code
5. **Single Instance**: Create one database manager instance and reuse it throughout your application

## Testing

Run tests with:
```bash
go test ./databases/ -v
```

## Dependencies

- [GORM](https://gorm.io/) - The fantastic ORM library for Go
- [MySQL Driver](https://github.com/go-sql-driver/mysql) - MySQL driver for Go
- [PostgreSQL Driver](https://github.com/lib/pq) - PostgreSQL driver for Go

## License

This package is part of the go-pkg project and follows the same license.
