package order_test

import (
	"errors"
	"reflect"

	"github.com.Saveliy12/microservices2.git/order/internal/client/db/mocks"
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *OrderRepositorySuite) TestUpdateOrder_Success() {
	orderUUID := "order-update-1"
	original := model.Order{OrderUUID: orderUUID, UserUUID: "user-1", TotalPrice: 10}

	// Мокаем Exec для CreateOrder
	s.dbMock.On("Exec", s.ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil)

	require.NoError(s.T(), s.repo.CreateOrder(s.ctx, original))

	updated := original
	updated.TotalPrice = 99.99
	updated.Status = model.StatusPaid
	updated.PartUUIDs = []string{"part-1", "part-2"}

	// Мокаем Exec для UpdateOrder
	s.dbMock.On("Exec", s.ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil)

	require.NoError(s.T(), s.repo.UpdateOrder(s.ctx, updated))

	// Мокаем QueryRow для GetOrder
	rowMock := &mocks.RowMock{}
	rowMock.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*string) = updated.OrderUUID
		*dest[1].(*string) = updated.UserUUID
		*dest[2].(*float64) = updated.TotalPrice
		reflect.ValueOf(dest[3]).Elem().SetString(string(updated.Status))
		*dest[6].(*[]string) = updated.PartUUIDs
	}).Return(nil)

	s.dbMock.On("QueryRow", s.ctx, mock.Anything, orderUUID).Return(rowMock)

	got, err := s.repo.GetOrder(s.ctx, orderUUID)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), updated.OrderUUID, got.OrderUUID)
	assert.Equal(s.T(), updated.UserUUID, got.UserUUID)
	assert.Equal(s.T(), updated.TotalPrice, got.TotalPrice)
	assert.Equal(s.T(), updated.Status, got.Status)
	assert.Equal(s.T(), updated.PartUUIDs, got.PartUUIDs)
}

func (s *OrderRepositorySuite) TestUpdateOrder_NotFound() {
	missing := model.Order{OrderUUID: "missing-order", UserUUID: "u", TotalPrice: 1}

	// Мокаем Exec для UpdateOrder
	s.dbMock.On("Exec", s.ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, errors.New("no rows affected"))

	err := s.repo.UpdateOrder(s.ctx, missing)
	require.Error(s.T(), err)
}
