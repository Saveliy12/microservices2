package order_test

import (
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *OrderServiceSuite) TestGetOrder_Success() {
	orderUUID := "order-123"
	want := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "user-1",
		TotalPrice: 99.99,
		Status:     model.StatusPendingPayment,
	}

	s.repoMock.EXPECT().
		GetOrder(s.ctx, orderUUID).
		Return(want, nil).
		Once()

	got, err := s.service.GetOrder(s.ctx, orderUUID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), want.OrderUUID, got.OrderUUID)
	assert.Equal(s.T(), want.UserUUID, got.UserUUID)
	assert.Equal(s.T(), want.TotalPrice, got.TotalPrice)
	assert.Equal(s.T(), want.Status, got.Status)
}

func (s *OrderServiceSuite) TestGetOrder_ReturnsError() {
	orderUUID := "missing-order"
	wantErr := model.ErrOrderNotFound

	s.repoMock.EXPECT().
		GetOrder(s.ctx, orderUUID).
		Return(model.Order{}, wantErr).
		Once()

	got, err := s.service.GetOrder(s.ctx, orderUUID)

	require.ErrorIs(s.T(), err, wantErr)
	assert.Equal(s.T(), model.Order{}, got)
}
