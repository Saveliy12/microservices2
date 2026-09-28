package env_test

import (
	"testing"

	env "github.com.Saveliy12/microservices2.git/order/internal/config/env"
	"github.com/stretchr/testify/suite"
)

type PaymentGRPCConfigTestSuite struct {
	suite.Suite
}

func (s *PaymentGRPCConfigTestSuite) TestNewPaymentGRPCConfig_FromEnv() {
	t := s.T()
	t.Setenv("PAYMENT_GRPC_PORT", "50053")

	cfg, err := env.NewPaymentGRPCConfig()

	s.Require().NoError(err)
	s.Require().Equal(50053, cfg.Port())
}

func (s *PaymentGRPCConfigTestSuite) TestNewPaymentGRPCConfig_Default() {
	t := s.T()
	t.Setenv("PAYMENT_GRPC_PORT", "")

	cfg, err := env.NewPaymentGRPCConfig()

	s.Require().NoError(err)
	s.Require().Equal(50052, cfg.Port())
}

func (s *PaymentGRPCConfigTestSuite) TestNewPaymentGRPCConfig_Invalid() {
	t := s.T()
	t.Setenv("PAYMENT_GRPC_PORT", "abc")

	cfg, err := env.NewPaymentGRPCConfig()

	s.Require().Error(err)
	s.Nil(cfg)
}

func TestPaymentGRPCConfigTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentGRPCConfigTestSuite))
}
