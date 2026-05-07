package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repoMocks "github.com.Saveliy12/microservices2.git/order/internal/repository/mocks"
	serviceorder "github.com.Saveliy12/microservices2.git/order/internal/service/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type mockInventoryClient struct {
	parts []model.PartQuote
	err   error
}

func (m *mockInventoryClient) ListParts(ctx context.Context, partUUIDs []string) ([]model.PartQuote, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.parts, nil
}

type mockPaymentClient struct{}

func (m *mockPaymentClient) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	return "", nil
}

type OrderCreateServiceSuite struct {
	suite.Suite
	service  *serviceorder.Service
	repoMock *repoMocks.OrderRepository
}

func (s *OrderCreateServiceSuite) SetupTest() {
	s.repoMock = repoMocks.NewOrderRepository(s.T())
}

func (s *OrderCreateServiceSuite) TestCreateOrder_Success() {
	invMock := &mockInventoryClient{
		parts: []model.PartQuote{
			{UUID: "part-1", Price: 100},
			{UUID: "part-2", Price: 50},
		},
	}
	s.service = serviceorder.NewService(s.repoMock, invMock, &mockPaymentClient{})

	userUUID := "user-1"
	partUUIDs := []string{"part-1", "part-2"}

	s.repoMock.EXPECT().
		CreateOrder(context.Background(), mock.MatchedBy(func(order model.Order) bool {
			return order.UserUUID == userUUID &&
				len(order.PartUUIDs) == 2 &&
				order.TotalPrice == 150 &&
				order.Status == model.StatusPendingPayment
		})).
		Return(nil).
		Once()

	got, err := s.service.CreateOrder(context.Background(), userUUID, partUUIDs)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), userUUID, got.UserUUID)
	assert.Equal(s.T(), float64(150), got.TotalPrice)
	assert.Equal(s.T(), model.StatusPendingPayment, got.Status)
	assert.NotEmpty(s.T(), got.OrderUUID)
}

func (s *OrderCreateServiceSuite) TestCreateOrder_InvalidInput() {
	invMock := &mockInventoryClient{}
	s.service = serviceorder.NewService(s.repoMock, invMock, &mockPaymentClient{})

	_, err := s.service.CreateOrder(context.Background(), "", []string{"part-1"})
	require.ErrorIs(s.T(), err, model.ErrInvalidOrderInput)

	_, err = s.service.CreateOrder(context.Background(), "user-1", []string{})
	require.ErrorIs(s.T(), err, model.ErrInvalidOrderInput)
}

func (s *OrderCreateServiceSuite) TestCreateOrder_InventoryError() {
	invMock := &mockInventoryClient{
		err: errors.New("inventory service down"),
	}
	s.service = serviceorder.NewService(s.repoMock, invMock, &mockPaymentClient{})

	_, err := s.service.CreateOrder(context.Background(), "user-1", []string{"part-1"})
	require.Error(s.T(), err)
	assert.Equal(s.T(), "inventory service down", err.Error())
}

func TestOrderCreateServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderCreateServiceSuite))
}
