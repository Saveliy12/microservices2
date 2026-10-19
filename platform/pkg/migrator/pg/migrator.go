package pg

import (
	"database/sql"
	"fmt"

	"github.com.Saveliy12/microservices2.git/platform/pkg/migrator"
	"github.com/pressly/goose/v3"
)

var _ migrator.Migrator = (*Migrator)(nil)

// Migrator применяет миграции PostgreSQL через goose
type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func New(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Up() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(m.db, m.migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
