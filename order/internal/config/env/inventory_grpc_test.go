package env_test

import (
	"testing"

	env "github.com.Saveliy12/microservices2.git/order/internal/config/env"
	"github.com/stretchr/testify/suite"
)

type InventoryGRPCConfigTestSuite struct {
	suite.Suite
}

func (s *InventoryGRPCConfigTestSuite) TestNewInventoryGRPCConfig_FromEnv() {
	t := s.T()
	t.Setenv("INVENTORY_GRPC_PORT", "50055")

	cfg, err := env.NewInventoryGRPCConfig()

	s.Require().NoError(err)
	s.Require().Equal(50055, cfg.Port())
}

func (s *InventoryGRPCConfigTestSuite) TestNewInventoryGRPCConfig_Default() {
	t := s.T()
	t.Setenv("INVENTORY_GRPC_PORT", "")

	cfg, err := env.NewInventoryGRPCConfig()

	s.Require().NoError(err)
	s.Require().Equal(50051, cfg.Port())
}

func (s *InventoryGRPCConfigTestSuite) TestNewInventoryGRPCConfig_Invalid() {
	t := s.T()
	t.Setenv("INVENTORY_GRPC_PORT", "abc")

	cfg, err := env.NewInventoryGRPCConfig()

	s.Require().Error(err)
	s.Nil(cfg)
}

func TestInventoryGRPCConfigTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryGRPCConfigTestSuite))
}
