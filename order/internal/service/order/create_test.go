package order_test

import (
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *OrderServiceSuite) TestCreateOrder_Success() {
	userUUID := "user-1"
	partUUIDs := []string{"part-1", "part-2"}
	parts := []model.PartQuote{
		{UUID: "part-1", Price: 100},
		{UUID: "part-2", Price: 50},
	}

	s.inventoryMock.EXPECT().
		ListParts(s.ctx, partUUIDs).
		Return(parts, nil).
		Once()

	s.repoMock.EXPECT().
		CreateOrder(s.ctx, mock.MatchedBy(func(order model.Order) bool {
			return order.UserUUID == userUUID &&
				len(order.PartUUIDs) == 2 &&
				order.TotalPrice == 150 &&
				order.Status == model.StatusPendingPayment
		})).
		Return(nil).
		Once()

	got, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), userUUID, got.UserUUID)
	assert.Equal(s.T(), float64(150), got.TotalPrice)
	assert.Equal(s.T(), model.StatusPendingPayment, got.Status)
	assert.NotEmpty(s.T(), got.OrderUUID)
}

func (s *OrderServiceSuite) TestCreateOrder_InvalidInput() {
	_, err := s.service.CreateOrder(s.ctx, "", []string{"part-1"})
	require.ErrorIs(s.T(), err, model.ErrInvalidOrderInput)

	_, err = s.service.CreateOrder(s.ctx, "user-1", []string{})
	require.ErrorIs(s.T(), err, model.ErrInvalidOrderInput)
}

func (s *OrderServiceSuite) TestCreateOrder_InventoryError() {
	userUUID := "user-1"
	partUUIDs := []string{"part-1"}
	inventoryErr := model.ErrOrderNotFound

	s.inventoryMock.EXPECT().
		ListParts(s.ctx, partUUIDs).
		Return(nil, inventoryErr).
		Once()

	_, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, inventoryErr)
}
