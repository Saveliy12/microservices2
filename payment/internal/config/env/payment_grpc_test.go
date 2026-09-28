package env_test

import (
	"testing"

	env "github.com.Saveliy12/microservices2.git/payment/internal/config/env"
	"github.com/stretchr/testify/suite"
)

type GrpcConfigPaymentServiceTestSuite struct {
	suite.Suite
}

func (s *GrpcConfigPaymentServiceTestSuite) TestNewGRPCConfig_FromEnv() {
	t := s.T()
	t.Setenv("GRPC_PORT", "50053")

	cfg, err := env.NewGRPCConfig()

	s.Require().NoError(err)
	s.Require().Equal(50053, cfg.Port())
}

func (s *GrpcConfigPaymentServiceTestSuite) TestNewGRPCConfig_Default() {
	t := s.T()
	t.Setenv("GRPC_PORT", "")

	cfg, err := env.NewGRPCConfig()

	s.Require().NoError(err)
	s.Require().Equal(50052, cfg.Port())
}

func (s *GrpcConfigPaymentServiceTestSuite) TestNewGRPCConfig_Invalid() {
	t := s.T()
	t.Setenv("GRPC_PORT", "abc")

	cfg, err := env.NewGRPCConfig()

	s.Require().Error(err)
	s.Nil(cfg)
}

func TestGrpcConfigTestSuite(t *testing.T) {
	suite.Run(t, new(GrpcConfigPaymentServiceTestSuite))
}
