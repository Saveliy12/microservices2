package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	serviceMocks "github.com.Saveliy12/microservices2.git/order/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PayOrderAPISuite struct {
	suite.Suite
	api         *API
	serviceMock *serviceMocks.OrderService
}

func (s *PayOrderAPISuite) SetupTest() {
	s.serviceMock = serviceMocks.NewOrderService(s.T())
	s.api = NewAPI(s.serviceMock)
}

func (s *PayOrderAPISuite) TestPayOrder_MethodNotAllowed() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/order-1/pay", nil)
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
	assert.JSONEq(s.T(), `{"error":"method not allowed"}`, rr.Body.String())
}

func (s *PayOrderAPISuite) TestPayOrder_InvalidOrderUUID() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid order uuid"}`, rr.Body.String())
}

func (s *PayOrderAPISuite) TestPayOrder_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/order-1/pay", strings.NewReader("{"))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid json"}`, rr.Body.String())
}

func (s *PayOrderAPISuite) TestPayOrder_PaymentMethodRequired() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/order-1/pay", strings.NewReader(`{"payment_method":""}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"payment_method is required"}`, rr.Body.String())
}

func (s *PayOrderAPISuite) TestPayOrder_InvalidPaymentMethod() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/order-1/pay", strings.NewReader(`{"payment_method":"BTC"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid payment_method"}`, rr.Body.String())
}

func (s *PayOrderAPISuite) TestPayOrder_Success() {
	orderUUID := "order-1"
	tx := "tx-1"
	pm := model.PaymentMethodCard
	order := model.Order{OrderUUID: orderUUID, TransactionUUID: &tx, PaymentMethod: &pm, Status: model.StatusPaid}

	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(order, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusOK, rr.Code)
	assert.JSONEq(s.T(), `{"transaction_uuid":"tx-1"}`, rr.Body.String())
}

func (s *PayOrderAPISuite) TestPayOrder_NotFound() {
	orderUUID := "missing"
	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(model.Order{}, model.ErrOrderNotFound).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusNotFound, rr.Code)
	assert.JSONEq(s.T(), `{"error":"order not found"}`, rr.Body.String())
}

func (s *PayOrderAPISuite) TestPayOrder_AlreadyPaid() {
	orderUUID := "order-1"
	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(model.Order{}, model.ErrOrderAlreadyPaid).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusConflict, rr.Code)
	assert.JSONEq(s.T(), `{"error":"order already paid"}`, rr.Body.String())
}

func (s *PayOrderAPISuite) TestPayOrder_Cancelled() {
	orderUUID := "order-1"
	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(model.Order{}, model.ErrOrderCancelled).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusConflict, rr.Code)
	assert.JSONEq(s.T(), `{"error":"order is cancelled"}`, rr.Body.String())
}

func (s *PayOrderAPISuite) TestPayOrder_ServiceError() {
	orderUUID := "order-1"
	errBoom := errors.New("payment backend down")
	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(model.Order{}, errBoom).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadGateway, rr.Code)
	assert.JSONEq(s.T(), `{"error":"payment backend down"}`, rr.Body.String())
}

func (s *PayOrderAPISuite) TestPayOrder_ResponseConversionError() {
	orderUUID := "order-1"
	pm := model.PaymentMethodCard
	order := model.Order{OrderUUID: orderUUID, PaymentMethod: &pm, Status: model.StatusPaid}
	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(order, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadGateway, rr.Code)
	assert.JSONEq(s.T(), `{"error":"transaction uuid is empty"}`, rr.Body.String())
}

func TestPayOrderAPISuite(t *testing.T) {
	suite.Run(t, new(PayOrderAPISuite))
}
