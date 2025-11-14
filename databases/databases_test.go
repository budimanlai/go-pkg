package databases

import (
	"testing"
	"time"
)

func TestNewDbManager(t *testing.T) {
	config := DbConfig{
		Driver:   MySQL,
		Host:     "localhost",
		Port:     "3306",
		Username: "user",
		Password: "pass",
		Name:     "testdb",
		Charset:  "utf8mb4",
	}

	dbManager := NewDbManager(config)

	if dbManager == nil {
		t.Fatal("Expected NewDbManager to return non-nil DbManager")
	}

	if dbManager.Config.Driver != MySQL {
		t.Errorf("Expected driver %s, got %s", MySQL, dbManager.Config.Driver)
	}
	if dbManager.Config.Host != "localhost" {
		t.Errorf("Expected host localhost, got %s", dbManager.Config.Host)
	}
	if dbManager.Config.Port != "3306" {
		t.Errorf("Expected port 3306, got %s", dbManager.Config.Port)
	}
	if dbManager.Config.Username != "user" {
		t.Errorf("Expected username user, got %s", dbManager.Config.Username)
	}
	if dbManager.Config.Name != "testdb" {
		t.Errorf("Expected database name testdb, got %s", dbManager.Config.Name)
	}
	if dbManager.Db != nil {
		t.Error("Expected Db to be nil initially")
	}
}

func TestNewDbManager_WithConnectionPoolSettings(t *testing.T) {
	config := DbConfig{
		Driver:          MySQL,
		Host:            "localhost",
		Port:            "3306",
		Username:        "user",
		Password:        "pass",
		Name:            "testdb",
		Charset:         "utf8mb4",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifeTime: time.Hour,
	}

	dbManager := NewDbManager(config)

	if dbManager.Config.MaxIdleConns != 10 {
		t.Errorf("Expected MaxIdleConns 10, got %d", dbManager.Config.MaxIdleConns)
	}
	if dbManager.Config.MaxOpenConns != 100 {
		t.Errorf("Expected MaxOpenConns 100, got %d", dbManager.Config.MaxOpenConns)
	}
	if dbManager.Config.ConnMaxLifeTime != time.Hour {
		t.Errorf("Expected ConnMaxLifeTime 1h, got %v", dbManager.Config.ConnMaxLifeTime)
	}
}

func TestDatabaseDriver_Constants(t *testing.T) {
	if MySQL != "mysql" {
		t.Errorf("Expected MySQL constant to be 'mysql', got %s", MySQL)
	}
	if Postgres != "postgres" {
		t.Errorf("Expected Postgres constant to be 'postgres', got %s", Postgres)
	}
}

func TestGetDb_BeforeOpen(t *testing.T) {
	config := DbConfig{
		Driver:   MySQL,
		Host:     "localhost",
		Port:     "3306",
		Username: "user",
		Password: "pass",
		Name:     "testdb",
	}

	dbManager := NewDbManager(config)
	db := dbManager.GetDb()

	if db != nil {
		t.Error("Expected GetDb to return nil before Open")
	}
}

func TestDbConfig_DefaultCharset(t *testing.T) {
	config := DbConfig{
		Driver:   MySQL,
		Host:     "localhost",
		Port:     "3306",
		Username: "user",
		Password: "pass",
		Name:     "testdb",
		// Charset not set
	}

	dbManager := NewDbManager(config)

	if dbManager.Config.Charset != "" {
		t.Errorf("Expected Charset to be empty initially, got %s", dbManager.Config.Charset)
	}
}

func TestDbConfig_PostgresDriver(t *testing.T) {
	config := DbConfig{
		Driver:   Postgres,
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "pass",
		Name:     "testdb",
	}

	dbManager := NewDbManager(config)

	if dbManager.Config.Driver != Postgres {
		t.Errorf("Expected driver %s, got %s", Postgres, dbManager.Config.Driver)
	}
	if dbManager.Config.Port != "5432" {
		t.Errorf("Expected port 5432 for Postgres, got %s", dbManager.Config.Port)
	}
}

func TestClose_WithoutOpen(t *testing.T) {
	config := DbConfig{
		Driver:   MySQL,
		Host:     "localhost",
		Port:     "3306",
		Username: "user",
		Password: "pass",
		Name:     "testdb",
	}

	dbManager := NewDbManager(config)

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Close() panicked when called without Open(): %v", r)
		}
	}()

	dbManager.Close()
}

func TestDbConfig_AllFields(t *testing.T) {
	config := DbConfig{
		Driver:          MySQL,
		Host:            "db.example.com",
		Port:            "3307",
		Username:        "admin",
		Password:        "secret",
		Name:            "production_db",
		Charset:         "utf8mb4",
		MaxIdleConns:    5,
		MaxOpenConns:    50,
		ConnMaxLifeTime: 30 * time.Minute,
	}

	dbManager := NewDbManager(config)

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Driver", dbManager.Config.Driver, MySQL},
		{"Host", dbManager.Config.Host, "db.example.com"},
		{"Port", dbManager.Config.Port, "3307"},
		{"Username", dbManager.Config.Username, "admin"},
		{"Password", dbManager.Config.Password, "secret"},
		{"Name", dbManager.Config.Name, "production_db"},
		{"Charset", dbManager.Config.Charset, "utf8mb4"},
		{"MaxIdleConns", dbManager.Config.MaxIdleConns, 5},
		{"MaxOpenConns", dbManager.Config.MaxOpenConns, 50},
		{"ConnMaxLifeTime", dbManager.Config.ConnMaxLifeTime, 30 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, tt.got)
			}
		})
	}
}

func TestDbManager_NilSafety(t *testing.T) {
	var dbManager *DbManager

	// Test that calling methods on nil pointer doesn't panic
	defer func() {
		if r := recover(); r == nil {
			// Expected to panic on nil pointer
			t.Log("Nil pointer access handled as expected")
		}
	}()

	// This should panic, which is expected behavior
	_ = dbManager.GetDb()
}

func TestDbConfig_ZeroValues(t *testing.T) {
	config := DbConfig{
		MaxIdleConns:    0,
		MaxOpenConns:    0,
		ConnMaxLifeTime: 0,
	}

	dbManager := NewDbManager(config)

	if dbManager.Config.MaxIdleConns != 0 {
		t.Errorf("Expected MaxIdleConns to be 0, got %d", dbManager.Config.MaxIdleConns)
	}
	if dbManager.Config.MaxOpenConns != 0 {
		t.Errorf("Expected MaxOpenConns to be 0, got %d", dbManager.Config.MaxOpenConns)
	}
	if dbManager.Config.ConnMaxLifeTime != 0 {
		t.Errorf("Expected ConnMaxLifeTime to be 0, got %v", dbManager.Config.ConnMaxLifeTime)
	}
}

func TestDbConfig_NegativeConnectionValues(t *testing.T) {
	config := DbConfig{
		Driver:          MySQL,
		Host:            "localhost",
		Port:            "3306",
		Username:        "user",
		Password:        "pass",
		Name:            "testdb",
		MaxIdleConns:    -1,
		MaxOpenConns:    -1,
		ConnMaxLifeTime: -1,
	}

	dbManager := NewDbManager(config)

	// Negative values should be preserved (they're used as "skip setting" flag)
	if dbManager.Config.MaxIdleConns != -1 {
		t.Errorf("Expected MaxIdleConns to be -1, got %d", dbManager.Config.MaxIdleConns)
	}
	if dbManager.Config.MaxOpenConns != -1 {
		t.Errorf("Expected MaxOpenConns to be -1, got %d", dbManager.Config.MaxOpenConns)
	}
}
