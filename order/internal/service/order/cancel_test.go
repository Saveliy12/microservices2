package order_test

import (
	"context"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repoMocks "github.com.Saveliy12/microservices2.git/order/internal/repository/mocks"
	serviceorder "github.com.Saveliy12/microservices2.git/order/internal/service/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderCancelServiceSuite struct {
	suite.Suite
	service  *serviceorder.Service
	repoMock *repoMocks.OrderRepository
}

func (s *OrderCancelServiceSuite) SetupTest() {
	s.repoMock = repoMocks.NewOrderRepository(s.T())
	s.service = serviceorder.NewService(s.repoMock, nil, nil)
}

func (s *OrderCancelServiceSuite) TestCancelOrder_Success() {
	orderUUID := "order-1"
	order := model.Order{OrderUUID: orderUUID, Status: model.StatusPendingPayment, TotalPrice: 100}

	s.repoMock.EXPECT().
		GetOrder(context.Background(), orderUUID).
		Return(order, nil).
		Once()

	s.repoMock.EXPECT().
		UpdateOrder(context.Background(), model.Order{OrderUUID: orderUUID, Status: model.StatusCancelled, TotalPrice: 100}).
		Return(nil).
		Once()

	got, err := s.service.CancelOrder(context.Background(), orderUUID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), model.StatusCancelled, got.Status)
}

func (s *OrderCancelServiceSuite) TestCancelOrder_NotFound() {
	s.repoMock.EXPECT().
		GetOrder(context.Background(), "missing").
		Return(model.Order{}, model.ErrOrderNotFound).
		Once()

	_, err := s.service.CancelOrder(context.Background(), "missing")
	require.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}

func (s *OrderCancelServiceSuite) TestCancelOrder_AlreadyPaid() {
	order := model.Order{OrderUUID: "order-1", Status: model.StatusPaid}

	s.repoMock.EXPECT().
		GetOrder(context.Background(), "order-1").
		Return(order, nil).
		Once()

	_, err := s.service.CancelOrder(context.Background(), "order-1")
	require.ErrorIs(s.T(), err, model.ErrOrderAlreadyPaid)
}

func (s *OrderCancelServiceSuite) TestCancelOrder_AlreadyCancelled() {
	order := model.Order{OrderUUID: "order-1", Status: model.StatusCancelled}

	s.repoMock.EXPECT().
		GetOrder(context.Background(), "order-1").
		Return(order, nil).
		Once()

	_, err := s.service.CancelOrder(context.Background(), "order-1")
	require.ErrorIs(s.T(), err, model.ErrOrderCancelled)
}

func TestOrderCancelServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderCancelServiceSuite))
}
