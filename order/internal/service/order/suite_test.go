package order_test

import (
	"context"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/repository/mocks"
	servicemocks "github.com.Saveliy12/microservices2.git/order/internal/service/mocks"
	serviceorder "github.com.Saveliy12/microservices2.git/order/internal/service/order"
	"github.com/stretchr/testify/suite"
)

type OrderServiceSuite struct {
	suite.Suite
	ctx               context.Context
	service           *serviceorder.Service
	repoMock          *mocks.OrderRepositoryMock
	inventoryMock     *servicemocks.InventoryClientMock
	paymentClientMock *servicemocks.PaymentClientMock
}

func (s *OrderServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repoMock = mocks.NewOrderRepositoryMock(s.T())
	s.inventoryMock = servicemocks.NewInventoryClientMock(s.T())
	s.paymentClientMock = servicemocks.NewPaymentClientMock(s.T())
	s.service = serviceorder.NewService(s.repoMock, s.inventoryMock, s.paymentClientMock)
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceSuite))
}
