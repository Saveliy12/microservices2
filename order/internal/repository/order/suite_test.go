package order_test

import (
	"context"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/client/db/mocks"
	repoorder "github.com.Saveliy12/microservices2.git/order/internal/repository/order"
	"github.com/stretchr/testify/suite"
)

type OrderRepositorySuite struct {
	suite.Suite
	ctx    context.Context
	repo   *repoorder.Repository
	dbMock *mocks.ClientMock
}

func (s *OrderRepositorySuite) SetupTest() {
	s.ctx = context.Background()
	s.dbMock = &mocks.ClientMock{}
	s.repo = repoorder.NewRepository(s.dbMock)
}

func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrderRepositorySuite))
}
