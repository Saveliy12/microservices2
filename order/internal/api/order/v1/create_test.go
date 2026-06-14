package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *OrderAPISuite) TestCreateOrder_MethodNotAllowed() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
	assert.JSONEq(s.T(), `{"error":"method not allowed"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestCreateOrder_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader("{"))
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid json"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestCreateOrder_InvalidInput() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{"user_uuid":"","part_uuids":[]}`))
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"user_uuid and part_uuids are required"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestCreateOrder_Success() {
	userUUID := "user-1"
	partUUIDs := []string{"p-1", "p-2"}
	order := model.Order{
		OrderUUID:  "order-1",
		TotalPrice: 150,
		Status:     model.StatusPendingPayment,
	}

	s.serviceMock.EXPECT().CreateOrder(mock.Anything, userUUID, partUUIDs).Return(order, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{"user_uuid":"user-1","part_uuids":["p-1","p-2"]}`))
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusCreated, rr.Code)
	var got map[string]any
	require.NoError(s.T(), json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(s.T(), "order-1", got["uuid"])
	assert.Equal(s.T(), float64(150), got["total_price"])
}

func (s *OrderAPISuite) TestCreateOrder_ServiceInvalidInput() {
	userUUID := "user-1"
	partUUIDs := []string{"p-1"}

	s.serviceMock.EXPECT().CreateOrder(mock.Anything, userUUID, partUUIDs).Return(model.Order{}, model.ErrInvalidOrderInput).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{"user_uuid":"user-1","part_uuids":["p-1"]}`))
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid order input"}`, rr.Body.String())
}

func (s *OrderAPISuite) TestCreateOrder_ServiceError() {
	userUUID := "user-1"
	partUUIDs := []string{"p-1"}
	errBoom := errors.New("unexpected error")

	s.serviceMock.EXPECT().CreateOrder(mock.Anything, userUUID, partUUIDs).Return(model.Order{}, errBoom).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{"user_uuid":"user-1","part_uuids":["p-1"]}`))
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusBadGateway, rr.Code)
	assert.JSONEq(s.T(), `{"error":"unexpected error"}`, rr.Body.String())
}
