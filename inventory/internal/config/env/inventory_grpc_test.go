package env_test

import (
	"testing"

	env "github.com.Saveliy12/microservices2.git/inventory/internal/config/env"
	"github.com/stretchr/testify/suite"
)

type GrpcConfigInventoryServiceTestSuite struct {
	suite.Suite
}

func (s *GrpcConfigInventoryServiceTestSuite) TestNewGRPCConfig_FromEnv() {
	t := s.T()
	t.Setenv("GRPC_PORT", "50055")

	cfg, err := env.NewGRPCConfig()

	s.Require().NoError(err)
	s.Require().Equal(50055, cfg.Port())
}

func (s *GrpcConfigInventoryServiceTestSuite) TestNewGRPCConfig_Default() {
	t := s.T()
	t.Setenv("GRPC_PORT", "")

	cfg, err := env.NewGRPCConfig()

	s.Require().NoError(err)
	s.Require().Equal(50051, cfg.Port())
}

func (s *GrpcConfigInventoryServiceTestSuite) TestNewGRPCConfig_Invalid() {
	t := s.T()
	t.Setenv("GRPC_PORT", "abc")

	cfg, err := env.NewGRPCConfig()

	s.Require().Error(err)
	s.Nil(cfg)
}

func TestGrpcConfigTestSuite(t *testing.T) {
	suite.Run(t, new(GrpcConfigInventoryServiceTestSuite))
}
