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

type OrderCreateSuite struct {
	suite.Suite
	repo *repoorder.Repository
}

func (s *OrderCreateSuite) SetupTest() {
	s.repo = repoorder.NewRepository()
}

func (s *OrderCreateSuite) TestCreateOrder_Success() {
	ctx := context.Background()
	orderUUID := "order-create-1"
	want := model.Order{OrderUUID: orderUUID, UserUUID: "user-1", TotalPrice: 99.99}

	err := s.repo.CreateOrder(ctx, want)
	require.NoError(s.T(), err)

	got, err := s.repo.GetOrder(ctx, orderUUID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), want.OrderUUID, got.OrderUUID)
	assert.Equal(s.T(), want.UserUUID, got.UserUUID)
	assert.Equal(s.T(), want.TotalPrice, got.TotalPrice)
}

func (s *OrderCreateSuite) TestCreateOrder_Overwrite() {
	ctx := context.Background()
	orderUUID := "order-create-2"
	first := model.Order{OrderUUID: orderUUID, UserUUID: "user-A", TotalPrice: 10}
	second := model.Order{OrderUUID: orderUUID, UserUUID: "user-A", TotalPrice: 20}

	require.NoError(s.T(), s.repo.CreateOrder(ctx, first))
	require.NoError(s.T(), s.repo.CreateOrder(ctx, second))

	got, err := s.repo.GetOrder(ctx, orderUUID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), second.TotalPrice, got.TotalPrice)
}

func TestOrderCreateSuite(t *testing.T) {
	suite.Run(t, new(OrderCreateSuite))
}
