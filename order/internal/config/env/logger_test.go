package env_test

import (
	"testing"

	env "github.com.Saveliy12/microservices2.git/order/internal/config/env"
	"github.com/stretchr/testify/suite"
)

type LoggerConfigTestSuite struct {
	suite.Suite
}

func (s *LoggerConfigTestSuite) TestNewLoggerConfig_Default() {
	t := s.T()
	t.Setenv("LOGGER_LEVEL", "")
	t.Setenv("LOGGER_AS_JSON", "")

	cfg, err := env.NewLoggerConfig()

	s.Require().NoError(err)
	s.Require().Equal("info", cfg.Level())
	s.Require().Equal(false, cfg.AsJSON())
}

func (s *LoggerConfigTestSuite) TestNewLoggerConfig_Custom() {
	t := s.T()
	t.Setenv("LOGGER_LEVEL", "debug")
	t.Setenv("LOGGER_AS_JSON", "true")

	cfg, err := env.NewLoggerConfig()

	s.Require().NoError(err)
	s.Require().Equal("debug", cfg.Level())
	s.Require().Equal(true, cfg.AsJSON())
}

func (s *LoggerConfigTestSuite) TestNewLoggerConfig_InvalidJSON() {
	t := s.T()
	t.Setenv("LOGGER_AS_JSON", "abc")

	cfg, err := env.NewLoggerConfig()

	s.Require().Error(err)
	s.Nil(cfg)
}

func TestLoggerConfigTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerConfigTestSuite))
}
