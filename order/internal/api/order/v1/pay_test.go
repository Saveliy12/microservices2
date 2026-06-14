package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *OrderAPISuite) TestPayOrder_MethodNotAllowed() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/order-1/pay", nil)
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
	assert.JSONEq(s.T(), `{"error":"method not allowed"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestPayOrder_InvalidOrderUUID() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid order uuid"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestPayOrder_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/order-1/pay", strings.NewReader("{"))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid json"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestPayOrder_PaymentMethodRequired() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/order-1/pay", strings.NewReader(`{"payment_method":""}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"payment_method is required"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestPayOrder_InvalidPaymentMethod() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/order-1/pay", strings.NewReader(`{"payment_method":"BTC"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid payment_method"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestPayOrder_Success() {
	orderUUID := "order-1"
	tx := "tx-1"
	pm := model.PaymentMethodCard
	order := model.Order{
		OrderUUID:       orderUUID,
		TransactionUUID: &tx,
		PaymentMethod:   &pm,
		Status:          model.StatusPaid,
	}

	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(order, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusOK, rr.Code)
	assert.JSONEq(s.T(), `{"transaction_uuid":"tx-1"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestPayOrder_NotFound() {
	orderUUID := "missing"
	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(model.Order{}, model.ErrOrderNotFound).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusNotFound, rr.Code)
	assert.JSONEq(s.T(), `{"error":"order not found"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestPayOrder_AlreadyPaid() {
	orderUUID := "order-1"
	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(model.Order{}, model.ErrOrderAlreadyPaid).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusConflict, rr.Code)
	assert.JSONEq(s.T(), `{"error":"order already paid"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestPayOrder_Cancelled() {
	orderUUID := "order-1"
	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(model.Order{}, model.ErrOrderCancelled).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusConflict, rr.Code)
	assert.JSONEq(s.T(), `{"error":"order is cancelled"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestPayOrder_ServiceError() {
	orderUUID := "order-1"
	errBoom := errors.New("unexpected error")
	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(model.Order{}, errBoom).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadGateway, rr.Code)
	assert.JSONEq(s.T(), `{"error":"unexpected error"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestPayOrder_ResponseConversionError() {
	orderUUID := "order-1"
	pm := model.PaymentMethodCard
	order := model.Order{
		OrderUUID:     orderUUID,
		PaymentMethod: &pm,
		Status:        model.StatusPaid,
	}
	s.serviceMock.EXPECT().PayOrder(mock.Anything, orderUUID, model.PaymentMethodCard).Return(order, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/pay", strings.NewReader(`{"payment_method":"CARD"}`))
	rr := httptest.NewRecorder()

	s.api.PayOrder(rr, req)

	require.Equal(s.T(), http.StatusBadGateway, rr.Code)
	assert.JSONEq(s.T(), `{"error":"transaction uuid is empty"}`, rr.Body.String())
}
