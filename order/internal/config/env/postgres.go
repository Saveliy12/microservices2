package env

import (
	"errors"
	"os"
)

const (
	postgresURIEnvName           = "DB_URI"
	postgresMigrationsDirEnvName = "MIGRATIONS_DIR"
	postgresMigrationsDirDefault = "order/migrations"
)

type postgresConfig struct {
	uri           string
	migrationsDir string
}

// NewPostgresConfig считывает DSN Postgres и путь до миграций из переменных окружения
func NewPostgresConfig() (*postgresConfig, error) {
	uri := os.Getenv(postgresURIEnvName)
	if uri == "" {
		return nil, errors.New("DB_URI is not set")
	}

	migrationsDir := os.Getenv(postgresMigrationsDirEnvName)
	if migrationsDir == "" {
		migrationsDir = postgresMigrationsDirDefault
	}

	return &postgresConfig{uri: uri, migrationsDir: migrationsDir}, nil
}

func (cfg *postgresConfig) URI() string {
	return cfg.uri
}

func (cfg *postgresConfig) MigrationsDir() string {
	return cfg.migrationsDir
}
