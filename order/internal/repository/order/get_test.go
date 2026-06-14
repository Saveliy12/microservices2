package order_test

import (
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/stretchr/testify/require"
)

func (s *OrderRepositorySuite) TestGetOrder_NotFound() {
	_, err := s.repo.GetOrder(s.ctx, "missing-order")
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}
