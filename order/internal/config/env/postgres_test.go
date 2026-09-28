package env_test

import (
	"testing"

	env "github.com.Saveliy12/microservices2.git/order/internal/config/env"
	"github.com/stretchr/testify/suite"
)

type PostgresConfigTestSuite struct {
	suite.Suite
}

func (s *PostgresConfigTestSuite) TestNewPostgresConfig_Valid() {
	t := s.T()
	t.Setenv("DB_URI", "postgres://demo:demo@localhost:5432/order_db?sslmode=disable")
	t.Setenv("MIGRATIONS_DIR", "")

	cfg, err := env.NewPostgresConfig()

	s.Require().NoError(err)
	s.Require().Equal("postgres://demo:demo@localhost:5432/order_db?sslmode=disable", cfg.URI())
	s.Require().Equal("order/migrations", cfg.MigrationsDir())
}

func (s *PostgresConfigTestSuite) TestNewPostgresConfig_MissingURI() {
	t := s.T()
	t.Setenv("DB_URI", "")

	cfg, err := env.NewPostgresConfig()

	s.Require().Error(err)
	s.Nil(cfg)
}

func (s *PostgresConfigTestSuite) TestNewPostgresConfig_CustomMigrationsDir() {
	t := s.T()
	t.Setenv("DB_URI", "postgres://user:pass@postgres:5432/test?sslmode=disable")
	t.Setenv("MIGRATIONS_DIR", "custom/migrations")

	cfg, err := env.NewPostgresConfig()

	s.Require().NoError(err)
	s.Require().Equal("postgres://user:pass@postgres:5432/test?sslmode=disable", cfg.URI())
	s.Require().Equal("custom/migrations", cfg.MigrationsDir())
}

func TestPostgresConfigTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresConfigTestSuite))
}
