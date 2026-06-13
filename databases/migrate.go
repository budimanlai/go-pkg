package databases

import (
	"fmt"
	"time"

	"github.com/budimanlai/go-pkg/v3/logger"
	"gorm.io/gorm"
)

// Migration is the interface every migration must implement.
// Each migration represents a single, versioned change to the database schema.
//
// Implementations must ensure that Up and Down are idempotent-safe and that
// Version returns a unique, lexicographically sortable string so migrations
// are applied in the correct order.
//
// Example implementation:
//
//	type CreateUsersTable struct{}
//
//	func (m *CreateUsersTable) Version() string     { return "20250101_000001" }
//	func (m *CreateUsersTable) Description() string { return "create users table" }
//
//	func (m *CreateUsersTable) Up(db *gorm.DB) error {
//	    return db.Exec(`CREATE TABLE users (id BIGINT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255))`).Error
//	}
//
//	func (m *CreateUsersTable) Down(db *gorm.DB) error {
//	    return db.Exec(`DROP TABLE IF EXISTS users`).Error
//	}
type Migration interface {
	// Version returns a unique sortable string identifier, e.g. "20250101_000001".
	// Migrations are applied in ascending Version order.
	Version() string

	// Description returns a human-readable summary of what this migration does.
	// It is stored in schema_migrations and displayed in Status output.
	Description() string

	// Up applies the migration (schema change, data seeding, etc.).
	Up(db *gorm.DB) error

	// Down rolls back the migration, undoing what Up applied.
	Down(db *gorm.DB) error
}

// Runner executes ordered migrations and tracks applied versions in the database.
// It maintains a schema_migrations table that records which migrations have been
// applied, allowing Up and Down to be called safely at any time.
//
// Example usage:
//
//	migrations := []databases.Migration{
//	    &CreateUsersTable{},
//	    &AddEmailToUsers{},
//	}
//
//	runner := databases.NewRunner(db, migrations)
//
//	// Apply all pending migrations
//	if err := runner.Up(); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Check migration status
//	runner.Status()
type Runner struct {
	db         *gorm.DB
	migrations []Migration
}

// NewRunner creates a Runner with the provided ordered migration list.
//
// Parameters:
//   - db: A connected *gorm.DB instance
//   - migrations: Ordered slice of Migration implementations; they will be applied
//     in the order given, so sort by Version before passing.
//
// Returns:
//   - *Runner: Ready to call Up, Down, or Status on.
//
// Example:
//
//	runner := databases.NewRunner(manager.GetDb(), []databases.Migration{
//	    &m001CreateUsers{},
//	    &m002AddEmailIndex{},
//	})
func NewRunner(db *gorm.DB, migrations []Migration) *Runner {
	return &Runner{db: db, migrations: migrations}
}

// ensureTable creates schema_migrations if it does not exist.
func (r *Runner) ensureTable() error {
	return r.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version     VARCHAR(50) PRIMARY KEY,
			description TEXT        NOT NULL DEFAULT '',
			applied_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

// appliedVersions returns a set of already-applied migration versions.
func (r *Runner) appliedVersions() (map[string]bool, error) {
	var versions []string
	if err := r.db.Raw("SELECT version FROM schema_migrations ORDER BY version").Scan(&versions).Error; err != nil {
		return nil, err
	}
	m := make(map[string]bool, len(versions))
	for _, v := range versions {
		m[v] = true
	}
	return m, nil
}

// Up runs all pending migrations in ascending Version order.
// It creates the schema_migrations table if it does not already exist, then
// applies every migration whose Version is not yet recorded, inserting a row
// for each one upon success.
//
// Returns nil when all pending migrations have been applied (or when there are
// none). Returns an error if the table cannot be created, if any migration's
// Up method fails, or if the version record cannot be inserted.
//
// Example:
//
//	if err := runner.Up(); err != nil {
//	    log.Fatalf("migration failed: %v", err)
//	}
func (r *Runner) Up() error {
	if err := r.ensureTable(); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}
	applied, err := r.appliedVersions()
	if err != nil {
		return fmt.Errorf("load applied versions: %w", err)
	}

	pending := 0
	for _, m := range r.migrations {
		if applied[m.Version()] {
			continue
		}
		pending++
		logger.Infof("[migrate] applying %s — %s", m.Version(), m.Description())
		if err := m.Up(r.db); err != nil {
			return fmt.Errorf("migration %s Up: %w", m.Version(), err)
		}
		if err := r.db.Exec(
			"INSERT INTO schema_migrations (version, description, applied_at) VALUES (?, ?, ?)",
			m.Version(), m.Description(), time.Now(),
		).Error; err != nil {
			return fmt.Errorf("record migration %s: %w", m.Version(), err)
		}
		logger.Infof("[migrate] %s applied", m.Version())
	}

	if pending == 0 {
		logger.Infof("[migrate] nothing to migrate — already up to date")
	} else {
		logger.Infof("[migrate] %d migration(s) applied", pending)
	}
	return nil
}

// Down rolls back the single most-recently applied migration.
// It iterates the migration list in reverse order, finds the last Version that
// is recorded in schema_migrations, calls its Down method, and removes the row.
//
// Only one migration is rolled back per call. To roll back multiple migrations
// call Down repeatedly.
//
// Returns nil if the rollback succeeds or if there is nothing to roll back.
// Returns an error if the schema_migrations table cannot be created, if the
// Down method fails, or if the version record cannot be deleted.
//
// Example:
//
//	if err := runner.Down(); err != nil {
//	    log.Fatalf("rollback failed: %v", err)
//	}
func (r *Runner) Down() error {
	if err := r.ensureTable(); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}
	applied, err := r.appliedVersions()
	if err != nil {
		return fmt.Errorf("load applied versions: %w", err)
	}

	for i := len(r.migrations) - 1; i >= 0; i-- {
		m := r.migrations[i]
		if !applied[m.Version()] {
			continue
		}
		logger.Infof("[migrate] rolling back %s — %s", m.Version(), m.Description())
		if err := m.Down(r.db); err != nil {
			return fmt.Errorf("migration %s Down: %w", m.Version(), err)
		}
		if err := r.db.Exec(
			"DELETE FROM schema_migrations WHERE version = ?", m.Version(),
		).Error; err != nil {
			return fmt.Errorf("remove migration record %s: %w", m.Version(), err)
		}
		logger.Infof("[migrate] %s rolled back", m.Version())
		return nil
	}

	logger.Infof("[migrate] nothing to roll back")
	return nil
}

// Status prints a formatted table of all registered migrations and their
// current state (applied / pending) to standard output.
//
// Output format:
//
//	VERSION               STATUS    DESCRIPTION
//	--------------------  --------  -----------
//	20250101_000001       applied   create users table
//	20250201_000002       pending   add email index
//
// Returns an error if the schema_migrations table cannot be created or if
// applied versions cannot be fetched.
//
// Example:
//
//	if err := runner.Status(); err != nil {
//	    log.Fatal(err)
//	}
func (r *Runner) Status() error {
	if err := r.ensureTable(); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}
	applied, err := r.appliedVersions()
	if err != nil {
		return fmt.Errorf("load applied versions: %w", err)
	}

	fmt.Printf("%-20s  %-8s  %s\n", "VERSION", "STATUS", "DESCRIPTION")
	fmt.Printf("%-20s  %-8s  %s\n", "--------------------", "--------", "-----------")
	for _, m := range r.migrations {
		status := "pending"
		if applied[m.Version()] {
			status = "applied"
		}
		fmt.Printf("%-20s  %-8s  %s\n", m.Version(), status, m.Description())
	}
	return nil
}
