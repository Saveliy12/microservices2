package env_test

import (
	"testing"

	env "github.com.Saveliy12/microservices2.git/payment/internal/config/env"
	"github.com/stretchr/testify/suite"
)

type LoggerConfigTestSuite struct {
	suite.Suite
}

func (s *LoggerConfigTestSuite) TestNewLoggerConfig_Default() {
	t := s.T()
	t.Setenv("LOGGER_LEVEL", "")
	t.Setenv("LOGGER_AS_JSON", "")

	logger, err := env.NewLoggerConfig()

	s.Require().NoError(err)
	s.Require().Equal("info", logger.Level())
	s.Require().Equal(false, logger.AsJSON())
}

func (s *LoggerConfigTestSuite) TestNewLoggerConfig_Custom() {
	t := s.T()
	t.Setenv("LOGGER_LEVEL", "debug")
	t.Setenv("LOGGER_AS_JSON", "true")

	logger, err := env.NewLoggerConfig()

	s.Require().NoError(err)
	s.Require().Equal("debug", logger.Level())
	s.Require().Equal(true, logger.AsJSON())
}

func (s *LoggerConfigTestSuite) TestNewLoggerConfig_InvalidJSON() {
	t := s.T()
	t.Setenv("LOGGER_AS_JSON", "abc")

	logger, err := env.NewLoggerConfig()

	s.Require().Error(err)
	s.Nil(logger)
}

func TestLoggerConfigTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerConfigTestSuite))
}
