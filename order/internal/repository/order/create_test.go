package order_test

import (
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *OrderRepositorySuite) TestCreateOrder_Success() {
	orderUUID := "order-create-1"
	want := model.Order{
		OrderUUID:   orderUUID,
		UserUUID:    "user-1",
		TotalPrice:  99.99,
		Status:      model.StatusPendingPayment,
	}

	err := s.repo.CreateOrder(s.ctx, want)
	require.NoError(s.T(), err)

	got, err := s.repo.GetOrder(s.ctx, orderUUID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), want.OrderUUID, got.OrderUUID)
	assert.Equal(s.T(), want.UserUUID, got.UserUUID)
	assert.Equal(s.T(), want.TotalPrice, got.TotalPrice)
	assert.Equal(s.T(), want.Status, got.Status)
}

func (s *OrderRepositorySuite) TestCreateOrder_Overwrite() {
	orderUUID := "order-create-2"
	first := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "user-A",
		TotalPrice: 10,
		Status:     model.StatusPendingPayment,
	}
	second := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "user-A",
		TotalPrice: 20,
		Status:     model.StatusPendingPayment,
	}

	require.NoError(s.T(), s.repo.CreateOrder(s.ctx, first))
	require.NoError(s.T(), s.repo.CreateOrder(s.ctx, second))

	got, err := s.repo.GetOrder(s.ctx, orderUUID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), second.TotalPrice, got.TotalPrice)
}
