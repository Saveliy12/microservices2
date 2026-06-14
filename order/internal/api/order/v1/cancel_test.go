package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *OrderAPISuite) TestCancelOrder_MethodNotAllowed() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/order-1/cancel", nil)
	rr := httptest.NewRecorder()

	s.api.CancelOrder(rr, req)

	require.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
	assert.JSONEq(s.T(), `{"error":"method not allowed"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestCancelOrder_InvalidOrderUUID() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/cancel", nil)
	rr := httptest.NewRecorder()

	s.api.CancelOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid order uuid"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestCancelOrder_Success() {
	orderUUID := "order-1"
	s.serviceMock.EXPECT().CancelOrder(mock.Anything, orderUUID).Return(model.Order{OrderUUID: orderUUID}, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/cancel", nil)
	rr := httptest.NewRecorder()

	s.api.CancelOrder(rr, req)

	require.Equal(s.T(), http.StatusNoContent, rr.Code)
	assert.Empty(s.T(), rr.Body.String())
}

func (s *OrderAPISuite) TestCancelOrder_NotFound() {
	orderUUID := "missing"
	s.serviceMock.EXPECT().CancelOrder(mock.Anything, orderUUID).Return(model.Order{}, model.ErrOrderNotFound).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/cancel", nil)
	rr := httptest.NewRecorder()

	s.api.CancelOrder(rr, req)

	require.Equal(s.T(), http.StatusNotFound, rr.Code)
	assert.JSONEq(s.T(), `{"error":"order not found"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestCancelOrder_AlreadyPaid() {
	orderUUID := "order-1"
	s.serviceMock.EXPECT().CancelOrder(mock.Anything, orderUUID).Return(model.Order{}, model.ErrOrderAlreadyPaid).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/cancel", nil)
	rr := httptest.NewRecorder()

	s.api.CancelOrder(rr, req)

	require.Equal(s.T(), http.StatusConflict, rr.Code)
	assert.JSONEq(s.T(), `{"error":"paid order cannot be cancelled"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestCancelOrder_AlreadyCancelled() {
	orderUUID := "order-1"
	s.serviceMock.EXPECT().CancelOrder(mock.Anything, orderUUID).Return(model.Order{}, model.ErrOrderCancelled).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/cancel", nil)
	rr := httptest.NewRecorder()

	s.api.CancelOrder(rr, req)

	require.Equal(s.T(), http.StatusConflict, rr.Code)
	assert.JSONEq(s.T(), `{"error":"order is cancelled"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestCancelOrder_ServiceError() {
	orderUUID := "order-1"
	errBoom := errors.New("unexpected error")
	s.serviceMock.EXPECT().CancelOrder(mock.Anything, orderUUID).Return(model.Order{}, errBoom).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderUUID+"/cancel", nil)
	rr := httptest.NewRecorder()

	s.api.CancelOrder(rr, req)

	require.Equal(s.T(), http.StatusBadGateway, rr.Code)
	assert.JSONEq(s.T(), `{"error":"unexpected error"}`, rr.Body.String())
}
