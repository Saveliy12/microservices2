package order_test

import (
	"context"
	"testing"

	repoorder "github.com.Saveliy12/microservices2.git/order/internal/repository/order"
	"github.com/stretchr/testify/suite"
)

type OrderRepositorySuite struct {
	suite.Suite
	ctx  context.Context
	repo *repoorder.Repository
}

func (s *OrderRepositorySuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = repoorder.NewRepository()
}

func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrderRepositorySuite))
}
