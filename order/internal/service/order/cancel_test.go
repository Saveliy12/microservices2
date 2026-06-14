package order_test

import (
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *OrderServiceSuite) TestCancelOrder_Success() {
	orderUUID := "order-1"
	order := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "user-1",
		Status:     model.StatusPendingPayment,
		TotalPrice: 100,
	}

	s.repoMock.EXPECT().
		GetOrder(s.ctx, orderUUID).
		Return(order, nil).
		Once()

	s.repoMock.EXPECT().
		UpdateOrder(s.ctx, model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   order.UserUUID,
			TotalPrice: order.TotalPrice,
			Status:     model.StatusCancelled,
		}).
		Return(nil).
		Once()

	got, err := s.service.CancelOrder(s.ctx, orderUUID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), model.StatusCancelled, got.Status)
	assert.Equal(s.T(), order.UserUUID, got.UserUUID)
	assert.Equal(s.T(), order.TotalPrice, got.TotalPrice)
}

func (s *OrderServiceSuite) TestCancelOrder_NotFound() {
	s.repoMock.EXPECT().
		GetOrder(s.ctx, "missing").
		Return(model.Order{}, model.ErrOrderNotFound).
		Once()

	_, err := s.service.CancelOrder(s.ctx, "missing")
	require.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}

func (s *OrderServiceSuite) TestCancelOrder_AlreadyPaid() {
	order := model.Order{OrderUUID: "order-1", Status: model.StatusPaid}

	s.repoMock.EXPECT().
		GetOrder(s.ctx, "order-1").
		Return(order, nil).
		Once()

	_, err := s.service.CancelOrder(s.ctx, "order-1")
	require.ErrorIs(s.T(), err, model.ErrOrderAlreadyPaid)
}

func (s *OrderServiceSuite) TestCancelOrder_AlreadyCancelled() {
	order := model.Order{OrderUUID: "order-1", Status: model.StatusCancelled}

	s.repoMock.EXPECT().
		GetOrder(s.ctx, "order-1").
		Return(order, nil).
		Once()

	_, err := s.service.CancelOrder(s.ctx, "order-1")
	require.ErrorIs(s.T(), err, model.ErrOrderCancelled)
}
