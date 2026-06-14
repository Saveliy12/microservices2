package order_test

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *OrderServiceSuite) TestPayOrder_Success() {
	orderUUID := "order-1"
	userUUID := "user-1"
	paymentMethod := model.PaymentMethodCreditCard
	transactionUUID := "txn-1"

	order := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		Status:     model.StatusPendingPayment,
		TotalPrice: 100,
	}

	s.repoMock.EXPECT().GetOrder(s.ctx, orderUUID).Return(order, nil).Once()

	s.paymentClientMock.EXPECT().
		PayOrder(s.ctx, orderUUID, userUUID, paymentMethod).
		Run(func(ctx context.Context, oUUID, uUUID string, method model.PaymentMethod) {
			require.Equal(s.T(), orderUUID, oUUID)
			require.Equal(s.T(), userUUID, uUUID)
			require.Equal(s.T(), paymentMethod, method)
		}).
		Return(transactionUUID, nil).
		Once()

	s.repoMock.EXPECT().UpdateOrder(s.ctx, mock.MatchedBy(func(o model.Order) bool {
		return o.OrderUUID == orderUUID && o.Status == model.StatusPaid && o.TransactionUUID != nil && *o.TransactionUUID == transactionUUID
	})).Return(nil).Once()

	got, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), model.StatusPaid, got.Status)
	assert.NotNil(s.T(), got.TransactionUUID)
	assert.Equal(s.T(), transactionUUID, *got.TransactionUUID)
}

func (s *OrderServiceSuite) TestPayOrder_AlreadyPaid() {
	orderUUID := "order-1"
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.StatusPaid,
	}

	s.repoMock.EXPECT().GetOrder(s.ctx, orderUUID).Return(order, nil).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, model.PaymentMethodCreditCard)
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, model.ErrOrderAlreadyPaid)
}

func (s *OrderServiceSuite) TestPayOrder_Cancelled() {
	orderUUID := "order-1"
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.StatusCancelled,
	}

	s.repoMock.EXPECT().GetOrder(s.ctx, orderUUID).Return(order, nil).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, model.PaymentMethodCreditCard)
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, model.ErrOrderCancelled)
}

func (s *OrderServiceSuite) TestPayOrder_RepositoryError() {
	orderUUID := "order-1"

	s.repoMock.EXPECT().GetOrder(s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderNotFound).Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, model.PaymentMethodCreditCard)
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, model.ErrOrderNotFound)
}

func (s *OrderServiceSuite) TestPayOrder_PaymentClientError() {
	orderUUID := "order-1"
	userUUID := "user-1"
	order := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		Status:     model.StatusPendingPayment,
		TotalPrice: 100,
	}

	s.repoMock.EXPECT().GetOrder(s.ctx, orderUUID).Return(order, nil).Once()

	paymentErr := model.ErrOrderNotFound
	s.paymentClientMock.EXPECT().
		PayOrder(s.ctx, orderUUID, userUUID, model.PaymentMethodCreditCard).
		Return("", paymentErr).
		Once()

	_, err := s.service.PayOrder(s.ctx, orderUUID, model.PaymentMethodCreditCard)
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, paymentErr)
}
