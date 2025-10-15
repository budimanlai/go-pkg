# Databases Package Documentation

Package `databases` menyediakan manajemen koneksi database menggunakan GORM, mendukung MySQL dan PostgreSQL.

## Fitur
- Manajemen koneksi database dengan GORM.
- Mendukung driver MySQL dan PostgreSQL.
- Konfigurasi database yang fleksibel.
- Method untuk open, close, dan get DB instance.

## Instalasi
Package ini adalah bagian dari `github.com/budimanlai/go-pkg`. Import sebagai:

```go
import "github.com/budimanlai/go-pkg/databases"
```

## Cara Pakai

### 1. Konfigurasi Database
```go
config := databases.DbConfig{
    Driver:   databases.MySQL, // atau databases.Postgres
    Host:     "localhost",
    Port:     "3306",
    Username: "user",
    Password: "password",
    Name:     "database_name",
    Charset:  "utf8mb4", // optional, default utf8mb4
}
```

### 2. Buat DbManager
```go
dbManager := databases.NewMySQLDb(config)
```

### 3. Open Koneksi
```go
// Dengan default config
dbManager.Open()

// Atau dengan custom gorm.Config
dbManager.OpenWithConfig(&gorm.Config{
    // custom config
})
```

### 4. Gunakan Database
```go
db := dbManager.GetDb()
// Gunakan db untuk query GORM
```

### 5. Close Koneksi
```go
dbManager.Close()
```

## Driver yang Didukung
- `MySQL`: Untuk database MySQL.
- `Postgres`: Untuk database PostgreSQL.

## Default Values
- Driver: MySQL (jika kosong).
- Charset: utf8mb4 (untuk MySQL).

## Cara Tests

### Menjalankan Unit Tests
```bash
# Dari root project
go test ./databases/

# Dengan verbose output
go test -v ./databases/

# Dengan coverage
go test -cover ./databases/
```

### Test Files
- `databases/mysql_test.go`: Berisi unit tests untuk DbManager.
- Tests mencakup konstruktor, getter, close, dll. (Test open di-skip karena log.Fatal).

### Contoh Output Test
```
=== RUN   TestNewMySQLDb
--- PASS: TestNewMySQLDb (0.00s)
=== RUN   TestGetDb_BeforeOpen
--- PASS: TestGetDb_BeforeOpen (0.00s)
...
PASS
ok      github.com/budimanlai/go-pkg/databases  0.309s
```

## Dependencies
- `gorm.io/gorm`
- `gorm.io/driver/mysql`
- `gorm.io/driver/postgres`
- `github.com/go-sql-driver/mysql`

## Lisensi
Sesuai dengan project utama.