package order_test

import (
	"context"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repoorder "github.com.Saveliy12/microservices2.git/order/internal/repository/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderUpdateSuite struct {
	suite.Suite
	repo *repoorder.Repository
}

func (s *OrderUpdateSuite) SetupTest() {
	s.repo = repoorder.NewRepository()
}

func (s *OrderUpdateSuite) TestUpdateOrder_Success() {
	ctx := context.Background()
	orderUUID := "order-update-1"
	original := model.Order{OrderUUID: orderUUID, UserUUID: "user-1", TotalPrice: 10}

	require.NoError(s.T(), s.repo.CreateOrder(ctx, original))

	// modify and update
	updated := original
	updated.TotalPrice = 99.99
	updated.Status = model.StatusPaid

	require.NoError(s.T(), s.repo.UpdateOrder(ctx, updated))

	got, err := s.repo.GetOrder(ctx, orderUUID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), updated.TotalPrice, got.TotalPrice)
	assert.Equal(s.T(), updated.Status, got.Status)
}

func (s *OrderUpdateSuite) TestUpdateOrder_NotFound() {
	ctx := context.Background()
	missing := model.Order{OrderUUID: "missing-order", UserUUID: "u", TotalPrice: 1}

	err := s.repo.UpdateOrder(ctx, missing)
	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}

func TestOrderUpdateSuite(t *testing.T) {
	suite.Run(t, new(OrderUpdateSuite))
}
