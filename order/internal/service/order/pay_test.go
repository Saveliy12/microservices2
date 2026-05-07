package order

import (
	"context"
	"errors"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repoMocks "github.com.Saveliy12/microservices2.git/order/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderPayServiceSuite struct {
	suite.Suite
	service           *Service
	repoMock          *repoMocks.OrderRepository
	paymentClientMock mockPaymentClient
}

type mockPaymentClient struct {
	payFn func(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error)
}

func (m *mockPaymentClient) PayOrder(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error) {
	return m.payFn(ctx, orderUUID, userUUID, method)
}

func (s *OrderPayServiceSuite) SetupTest() {
	s.repoMock = &repoMocks.OrderRepository{}
	s.paymentClientMock = mockPaymentClient{}
	s.service = &Service{
		orderRepository: s.repoMock,
		paymentClient:   &s.paymentClientMock,
	}
}

func (s *OrderPayServiceSuite) TestPayOrder_Success() {
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

	s.repoMock.EXPECT().GetOrder(context.Background(), orderUUID).Return(order, nil).Once()

	paymentCalled := false
	s.paymentClientMock.payFn = func(ctx context.Context, oUUID, uUUID string, method model.PaymentMethod) (string, error) {
		paymentCalled = true
		require.Equal(s.T(), orderUUID, oUUID)
		require.Equal(s.T(), userUUID, uUUID)
		require.Equal(s.T(), paymentMethod, method)
		return transactionUUID, nil
	}

	s.repoMock.EXPECT().UpdateOrder(context.Background(), mock.MatchedBy(func(o model.Order) bool {
		require.True(s.T(), paymentCalled, "PayOrder must be called before UpdateOrder")
		return o.OrderUUID == orderUUID && o.Status == model.StatusPaid && o.TransactionUUID != nil && *o.TransactionUUID == transactionUUID
	})).Return(nil).Once()

	got, err := s.service.PayOrder(context.Background(), orderUUID, paymentMethod)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), model.StatusPaid, got.Status)
	assert.NotNil(s.T(), got.TransactionUUID)
	assert.Equal(s.T(), transactionUUID, *got.TransactionUUID)
}

func (s *OrderPayServiceSuite) TestPayOrder_AlreadyPaid() {
	orderUUID := "order-1"
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.StatusPaid,
	}

	s.repoMock.EXPECT().GetOrder(context.Background(), orderUUID).Return(order, nil).Once()

	_, err := s.service.PayOrder(context.Background(), orderUUID, model.PaymentMethodCreditCard)
	require.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrOrderAlreadyPaid, err)
}

func (s *OrderPayServiceSuite) TestPayOrder_Cancelled() {
	orderUUID := "order-1"
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.StatusCancelled,
	}

	s.repoMock.EXPECT().GetOrder(context.Background(), orderUUID).Return(order, nil).Once()

	_, err := s.service.PayOrder(context.Background(), orderUUID, model.PaymentMethodCreditCard)
	require.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrOrderCancelled, err)
}

func (s *OrderPayServiceSuite) TestPayOrder_RepositoryError() {
	orderUUID := "order-1"

	s.repoMock.EXPECT().GetOrder(context.Background(), orderUUID).Return(model.Order{}, model.ErrOrderNotFound).Once()

	_, err := s.service.PayOrder(context.Background(), orderUUID, model.PaymentMethodCreditCard)
	require.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrOrderNotFound, err)
}

func (s *OrderPayServiceSuite) TestPayOrder_PaymentClientError() {
	orderUUID := "order-1"
	userUUID := "user-1"
	order := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		Status:     model.StatusPendingPayment,
		TotalPrice: 100,
	}

	s.repoMock.EXPECT().GetOrder(context.Background(), orderUUID).Return(order, nil).Once()

	paymentErr := errors.New("payment service error")
	s.paymentClientMock.payFn = func(ctx context.Context, oUUID, uUUID string, method model.PaymentMethod) (string, error) {
		return "", paymentErr
	}

	_, err := s.service.PayOrder(context.Background(), orderUUID, model.PaymentMethodCreditCard)
	require.Error(s.T(), err)
	assert.Equal(s.T(), paymentErr, err)
}

func TestOrderPayServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderPayServiceSuite))
}
