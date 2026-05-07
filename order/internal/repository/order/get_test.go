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

type OrderRepositorySuite struct {
	suite.Suite
	repo *repoorder.Repository
}

func (s *OrderRepositorySuite) SetupTest() {
	s.repo = repoorder.NewRepository()
}

func (s *OrderRepositorySuite) TestGetOrder_Success() {
	orderUUID := "order-123"
	want := model.Order{OrderUUID: orderUUID, UserUUID: "user-1", TotalPrice: 123.45}

	// seed
	err := s.repo.CreateOrder(context.Background(), want)
	require.NoError(s.T(), err)

	got, err := s.repo.GetOrder(context.Background(), orderUUID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), want.OrderUUID, got.OrderUUID)
	assert.Equal(s.T(), want.UserUUID, got.UserUUID)
	assert.Equal(s.T(), want.TotalPrice, got.TotalPrice)
}

func (s *OrderRepositorySuite) TestGetOrder_NotFound() {
	_, err := s.repo.GetOrder(context.Background(), "missing-order")
	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}

func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrderRepositorySuite))
}
