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

type OrderGetServiceSuite struct {
	suite.Suite
	service  *serviceorder.Service
	repoMock *repoMocks.OrderRepository
}

func (s *OrderGetServiceSuite) SetupTest() {
	s.repoMock = repoMocks.NewOrderRepository(s.T())
	s.service = serviceorder.NewService(s.repoMock, nil, nil)
}

func (s *OrderGetServiceSuite) TestGetOrder_Success() {
	orderUUID := "order-123"
	want := model.Order{OrderUUID: orderUUID, UserUUID: "user-1", TotalPrice: 99.99}

	s.repoMock.EXPECT().
		GetOrder(context.Background(), orderUUID).
		Return(want, nil).
		Once()

	got, err := s.service.GetOrder(context.Background(), orderUUID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), want.OrderUUID, got.OrderUUID)
	assert.Equal(s.T(), want.UserUUID, got.UserUUID)
	assert.Equal(s.T(), want.TotalPrice, got.TotalPrice)
}

func (s *OrderGetServiceSuite) TestGetOrder_ReturnsError() {
	orderUUID := "missing-order"
	wantErr := model.ErrOrderNotFound

	s.repoMock.EXPECT().
		GetOrder(context.Background(), orderUUID).
		Return(model.Order{}, wantErr).
		Once()

	got, err := s.service.GetOrder(context.Background(), orderUUID)

	require.ErrorIs(s.T(), err, wantErr)
	assert.Equal(s.T(), model.Order{}, got)
}

func TestOrderGetServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderGetServiceSuite))
}
