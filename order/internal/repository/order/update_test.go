package order_test

import (
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *OrderRepositorySuite) TestUpdateOrder_Success() {
	orderUUID := "order-update-1"
	original := model.Order{OrderUUID: orderUUID, UserUUID: "user-1", TotalPrice: 10}

	require.NoError(s.T(), s.repo.CreateOrder(s.ctx, original))

	updated := original
	updated.TotalPrice = 99.99
	updated.Status = model.StatusPaid

	require.NoError(s.T(), s.repo.UpdateOrder(s.ctx, updated))

	got, err := s.repo.GetOrder(s.ctx, orderUUID)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), updated.OrderUUID, got.OrderUUID)
	assert.Equal(s.T(), updated.UserUUID, got.UserUUID)
	assert.Equal(s.T(), updated.TotalPrice, got.TotalPrice)
	assert.Equal(s.T(), updated.Status, got.Status)
}

func (s *OrderRepositorySuite) TestUpdateOrder_NotFound() {
	missing := model.Order{OrderUUID: "missing-order", UserUUID: "u", TotalPrice: 1}

	err := s.repo.UpdateOrder(s.ctx, missing)
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}
