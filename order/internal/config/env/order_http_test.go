package env_test

import (
	"testing"

	env "github.com.Saveliy12/microservices2.git/order/internal/config/env"
	"github.com/stretchr/testify/suite"
)

type HTTPConfigTestSuite struct {
	suite.Suite
}

func (s *HTTPConfigTestSuite) TestNewHTTPConfig_FromEnv() {
	t := s.T()
	t.Setenv("HTTP_PORT", "9090")

	cfg, err := env.NewHTTPConfig()

	s.Require().NoError(err)
	s.Require().Equal(9090, cfg.Port())
}

func (s *HTTPConfigTestSuite) TestNewHTTPConfig_Default() {
	t := s.T()
	t.Setenv("HTTP_PORT", "")

	cfg, err := env.NewHTTPConfig()

	s.Require().NoError(err)
	s.Require().Equal(8080, cfg.Port())
}

func (s *HTTPConfigTestSuite) TestNewHTTPConfig_Invalid() {
	t := s.T()
	t.Setenv("HTTP_PORT", "abc")

	cfg, err := env.NewHTTPConfig()

	s.Require().Error(err)
	s.Nil(cfg)
}

func TestHTTPConfigTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPConfigTestSuite))
}
