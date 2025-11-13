package databases

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseDriver represents the type of database driver to use for connections.
// It provides a type-safe way to specify database drivers in configurations.
type DatabaseDriver string

const (
	// MySQL represents the MySQL database driver
	MySQL DatabaseDriver = "mysql"
	// Postgres represents the PostgreSQL database driver
	Postgres DatabaseDriver = "postgres"
)

// DbConfig holds the configuration parameters for establishing a database connection.
// It includes connection details, authentication credentials, and connection pool settings.
//
// Example usage:
//
//	config := DbConfig{
//	    Driver:          MySQL,
//	    Host:            "localhost",
//	    Port:            "3306",
//	    Username:        "root",
//	    Password:        "password",
//	    Name:            "mydb",
//	    Charset:         "utf8mb4",
//	    MaxIdleConns:    10,
//	    MaxOpenConns:    100,
//	    ConnMaxLifeTime: time.Hour,
//	}
type DbConfig struct {
	// Driver specifies the database driver type (MySQL or Postgres)
	Driver DatabaseDriver

	// Host is the address of the database server (e.g., "localhost" or "127.0.0.1")
	Host string

	// Port is the port number where the database server is listening (e.g., "3306" for MySQL)
	Port string

	// Username is the database user for authentication
	Username string

	// Password is the database password for authentication
	Password string

	// Name is the name of the database to connect to
	Name string

	// Charset defines the character set for the database connection (default: "utf8mb4")
	Charset string

	// MaxIdleConns sets the maximum number of connections in the idle connection pool.
	// Use 0 or negative value to skip setting this parameter.
	MaxIdleConns int

	// MaxOpenConns sets the maximum number of open connections to the database.
	// Use 0 or negative value to skip setting this parameter.
	MaxOpenConns int

	// ConnMaxLifeTime sets the maximum amount of time a connection may be reused.
	// Use 0 or negative value to skip setting this parameter.
	ConnMaxLifeTime time.Duration
}

// DbManager manages database connections and operations using GORM.
// It encapsulates the database instance and configuration, providing a clean interface
// for database operations and connection management.
//
// Fields:
//   - Db: The GORM database instance for executing queries and database operations
//   - Config: The database configuration used to establish the connection
//
// Example usage:
//
//	manager := NewDbManager(config)
//	manager.Open()
//	defer manager.Close()
//	db := manager.GetDb()
type DbManager struct {
	Db     *gorm.DB
	Config DbConfig
}

// NewDbManager creates and initializes a new DbManager instance with the provided configuration.
// The database connection is not established until Open() or OpenWithConfig() is called.
//
// Parameters:
//   - config: DbConfig containing all database connection parameters
//
// Returns:
//   - *DbManager: A new DbManager instance ready to establish a database connection
//
// Example:
//
//	config := DbConfig{
//	    Driver:   MySQL,
//	    Host:     "localhost",
//	    Port:     "3306",
//	    Username: "root",
//	    Password: "password",
//	    Name:     "mydb",
//	}
//	manager := NewDbManager(config)
func NewDbManager(config DbConfig) *DbManager {
	return &DbManager{
		Config: config,
	}
}

// GetDb returns the underlying GORM database instance.
// This instance can be used to perform database operations such as queries,
// migrations, transactions, and other GORM operations.
//
// Returns:
//   - *gorm.DB: The GORM database instance
//
// Example:
//
//	db := manager.GetDb()
//	db.AutoMigrate(&User{})
//	db.Create(&User{Name: "John"})
func (m *DbManager) GetDb() *gorm.DB {
	return m.Db
}

// OpenWithConfig establishes a database connection using the provided GORM configuration.
// It constructs the DSN (Data Source Name) based on the driver type and applies
// connection pool settings. The method will terminate the program if connection fails.
//
// Default values:
//   - If Driver is empty, defaults to MySQL
//   - If Charset is empty, defaults to "utf8mb4"
//
// Parameters:
//   - cfg: *gorm.Config for customizing GORM behavior (logger, naming strategy, etc.)
//
// Panics:
//   - Terminates the program with log.Fatalf if connection fails
//
// Example:
//
//	manager := NewDbManager(config)
//	gormConfig := &gorm.Config{
//	    Logger: logger.Default.LogMode(logger.Info),
//	}
//	manager.OpenWithConfig(gormConfig)
func (m *DbManager) OpenWithConfig(cfg *gorm.Config) {
	if m.Config.Driver == "" {
		m.Config.Driver = MySQL
	}

	if m.Config.Charset == "" {
		m.Config.Charset = "utf8mb4"
	}

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		m.Config.Username,
		m.Config.Password,
		m.Config.Host,
		m.Config.Port,
		m.Config.Name,
		m.Config.Charset,
	)

	switch m.Config.Driver {
	case MySQL:
		m.Db, err = gorm.Open(mysql.Open(dsn), cfg)
	case Postgres:
		m.Db, err = gorm.Open(postgres.Open(dsn), cfg)
	}

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := m.Db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	if m.Config.MaxIdleConns >= 0 {
		sqlDB.SetMaxIdleConns(m.Config.MaxIdleConns)
	}

	if m.Config.MaxOpenConns >= 0 {
		sqlDB.SetMaxOpenConns(m.Config.MaxOpenConns)
	}

	if m.Config.ConnMaxLifeTime >= 0 {
		sqlDB.SetConnMaxLifetime(m.Config.ConnMaxLifeTime)
	}
}

// Open establishes a database connection using the default GORM configuration.
// This is a convenience method that calls OpenWithConfig with an empty gorm.Config.
//
// Default GORM behaviors will be applied. For custom configuration such as
// custom logger or naming strategy, use OpenWithConfig instead.
//
// Panics:
//   - Terminates the program with log.Fatalf if connection fails
//
// Example:
//
//	manager := NewDbManager(config)
//	manager.Open()
//	defer manager.Close()
func (m *DbManager) Open() {
	m.OpenWithConfig(&gorm.Config{})
}

// Close gracefully closes the database connection and releases all resources.
// It safely handles cases where the database connection was never established.
// Any errors during closing are logged but do not cause the program to terminate.
//
// It's recommended to defer this method immediately after opening a connection.
//
// Example:
//
//	manager := NewDbManager(config)
//	manager.Open()
//	defer manager.Close()
func (m *DbManager) Close() {
	if m.Db == nil {
		return
	}
	sqlDB, err := m.Db.DB()
	if err != nil {
		log.Printf("failed to get sql.DB from gorm.DB: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("failed to close database connection: %v", err)
	}
}
