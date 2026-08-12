package order_test

import (
	"reflect"

	"github.com.Saveliy12/microservices2.git/order/internal/client/db/mocks"
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *OrderRepositorySuite) TestCreateOrder_Success() {
	orderUUID := "order-create-1"
	want := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "user-1",
		TotalPrice: 99.99,
		Status:     model.StatusPendingPayment,
	}

	s.dbMock.On("Exec", s.ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil)

	err := s.repo.CreateOrder(s.ctx, want)
	require.NoError(s.T(), err)
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

	s.dbMock.On("Exec", s.ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil).Twice()

	require.NoError(s.T(), s.repo.CreateOrder(s.ctx, first))
	require.NoError(s.T(), s.repo.CreateOrder(s.ctx, second))

	rowMock := &mocks.RowMock{}
	rowMock.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*string) = second.OrderUUID
		*dest[1].(*string) = second.UserUUID
		*dest[2].(*float64) = second.TotalPrice
		reflect.ValueOf(dest[3]).Elem().SetString(string(second.Status))
		*dest[6].(*[]string) = second.PartUUIDs
	}).Return(nil)

	s.dbMock.On("QueryRow", s.ctx, mock.Anything, orderUUID).Return(rowMock)

	got, err := s.repo.GetOrder(s.ctx, orderUUID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), second.TotalPrice, got.TotalPrice)
}
