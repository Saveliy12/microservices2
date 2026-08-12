package v1

import (
	"context"
	"errors"
	"testing"

	"github.com.Saveliy12/microservices2.git/payment/internal/model"
	"github.com.Saveliy12/microservices2.git/payment/internal/service/mocks"
	paymentv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/payment/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type APITestSuite struct {
	suite.Suite
	api                *API
	paymentServiceMock *mocks.PaymentServiceMock
}

func (s *APITestSuite) SetupTest() {
	s.paymentServiceMock = mocks.NewPaymentServiceMock(s.T())
	s.api = NewAPI(s.paymentServiceMock)
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) TestPayOrder_Success() {
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "660e8400-e29b-41d4-a716-446655440001"
	transactionUUID := "770e8400-e29b-41d4-a716-446655440002"

	s.paymentServiceMock.EXPECT().
		PayOrder(
			context.Background(),
			model.Payment{
				OrderUUID:     orderUUID,
				UserUUID:      userUUID,
				PaymentMethod: "PAYMENT_METHOD_CARD",
			},
		).
		Return(transactionUUID, nil).
		Once()

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	resp, err := s.api.PayOrder(context.Background(), req)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), transactionUUID, resp.TransactionUuid)
}

func (s *APITestSuite) TestPayOrder_NilRequest() {
	resp, err := s.api.PayOrder(context.Background(), nil)

	assert.Nil(s.T(), resp)
	assert.Error(s.T(), err)
	st, _ := status.FromError(err)
	assert.Equal(s.T(), codes.InvalidArgument, st.Code())
}

func (s *APITestSuite) TestPayOrder_ServiceError() {
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"

	s.paymentServiceMock.EXPECT().
		PayOrder(context.Background(), model.Payment{
			OrderUUID:     orderUUID,
			PaymentMethod: "PAYMENT_METHOD_SBP",
		}).
		Return("", errors.New("payment failed")).
		Once()

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
	}

	resp, err := s.api.PayOrder(context.Background(), req)

	assert.Nil(s.T(), resp)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), "payment failed", err.Error())
}
