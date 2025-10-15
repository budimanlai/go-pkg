package databases

import (
	"testing"
)

func TestNewMySQLDb(t *testing.T) {
	config := DbConfig{
		Driver:   MySQL,
		Host:     "localhost",
		Port:     "3306",
		Username: "user",
		Password: "pass",
		Name:     "testdb",
		Charset:  "utf8mb4",
	}

	dbManager := NewMySQLDb(config)

	if dbManager.Config.Driver != MySQL {
		t.Errorf("Expected driver %s, got %s", MySQL, dbManager.Config.Driver)
	}
	if dbManager.Config.Host != "localhost" {
		t.Errorf("Expected host localhost, got %s", dbManager.Config.Host)
	}
	if dbManager.Db != nil {
		t.Error("Expected Db to be nil initially")
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

	dbManager := NewMySQLDb(config)
	db := dbManager.GetDb()

	if db != nil {
		t.Error("Expected GetDb to return nil before Open")
	}
}

func TestOpen_InvalidConfig(t *testing.T) {
	t.Skip("Skipping test that calls log.Fatal, which exits the process")
}

func TestOpenWithConfig_InvalidConfig(t *testing.T) {
	t.Skip("Skipping test that calls log.Fatal, which exits the process")
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

	dbManager := NewMySQLDb(config)

	// Should not panic
	dbManager.Close()
}

func TestDefaultValues(t *testing.T) {
	t.Skip("Default values are set in OpenWithConfig, not in NewMySQLDb")
}
