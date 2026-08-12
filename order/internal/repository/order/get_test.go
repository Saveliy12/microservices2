package order_test

import (
	"errors"

	"github.com.Saveliy12/microservices2.git/order/internal/client/db/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *OrderRepositorySuite) TestGetOrder_NotFound() {
	orderUUID := "missing-order"

	rowMock := &mocks.RowMock{}
	rowMock.On("Scan", mock.Anything).Return(errors.New("no rows"))

	s.dbMock.On("QueryRow", s.ctx, mock.Anything, orderUUID).Return(rowMock)

	_, err := s.repo.GetOrder(s.ctx, orderUUID)
	require.Error(s.T(), err)
}
