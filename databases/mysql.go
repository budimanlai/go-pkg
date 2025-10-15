package databases

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseDriver string

const (
	MySQL    DatabaseDriver = "mysql"
	Postgres DatabaseDriver = "postgres"
)

type DbConfig struct {
	Driver   DatabaseDriver
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Charset  string
}

// DbConfig holds the configuration for the database connection.
type DbManager struct {
	Db     *gorm.DB
	Config DbConfig
}

// NewMySQLDb creates a new DbManager instance with the provided configuration.
func NewMySQLDb(config DbConfig) *DbManager {
	return &DbManager{
		Config: config,
	}
}

// GetDb returns the gorm.DB instance.
func (m *DbManager) GetDb() *gorm.DB {
	return m.Db
}

// OpenWithConfig opens the database connection with the provided gorm.Config.
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
}

// Open opens the database connection with default gorm.Config.
func (m *DbManager) Open() {
	m.OpenWithConfig(&gorm.Config{})
}

// Close closes the database connection.
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
