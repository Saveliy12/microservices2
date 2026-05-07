package v1

import (
	"encoding/json"
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

type CreateOrderAPISuite struct {
	suite.Suite
	api         *API
	serviceMock *serviceMocks.OrderService
}

func (s *CreateOrderAPISuite) SetupTest() {
	s.serviceMock = serviceMocks.NewOrderService(s.T())
	s.api = NewAPI(s.serviceMock)
}

func (s *CreateOrderAPISuite) TestCreateOrder_MethodNotAllowed() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
	assert.JSONEq(s.T(), `{"error":"method not allowed"}`, rr.Body.String())
}

func (s *CreateOrderAPISuite) TestCreateOrder_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader("{"))
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid json"}`, rr.Body.String())
}

func (s *CreateOrderAPISuite) TestCreateOrder_InvalidInput() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{"user_uuid":"","part_uuids":[]}`))
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"user_uuid and part_uuids are required"}`, rr.Body.String())
}

func (s *CreateOrderAPISuite) TestCreateOrder_Success() {
	userUUID := "user-1"
	partUUIDs := []string{"p-1", "p-2"}
	order := model.Order{OrderUUID: "order-1", TotalPrice: 150}

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

func (s *CreateOrderAPISuite) TestCreateOrder_ServiceInvalidInput() {
	userUUID := "user-1"
	partUUIDs := []string{"p-1"}

	s.serviceMock.EXPECT().CreateOrder(mock.Anything, userUUID, partUUIDs).Return(model.Order{}, model.ErrInvalidOrderInput).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{"user_uuid":"user-1","part_uuids":["p-1"]}`))
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid order input"}`, rr.Body.String())
}

func (s *CreateOrderAPISuite) TestCreateOrder_ServiceError() {
	userUUID := "user-1"
	partUUIDs := []string{"p-1"}
	errBoom := errors.New("inventory unavailable")

	s.serviceMock.EXPECT().CreateOrder(mock.Anything, userUUID, partUUIDs).Return(model.Order{}, errBoom).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(`{"user_uuid":"user-1","part_uuids":["p-1"]}`))
	rr := httptest.NewRecorder()

	s.api.CreateOrder(rr, req)

	require.Equal(s.T(), http.StatusBadGateway, rr.Code)
	assert.JSONEq(s.T(), `{"error":"inventory unavailable"}`, rr.Body.String())
}

func TestCreateOrderAPISuite(t *testing.T) {
	suite.Run(t, new(CreateOrderAPISuite))
}
