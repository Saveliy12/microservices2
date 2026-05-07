package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	serviceMocks "github.com.Saveliy12/microservices2.git/order/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetOrderAPISuite struct {
	suite.Suite
	api         *API
	serviceMock *serviceMocks.OrderService
}

func (s *GetOrderAPISuite) SetupTest() {
	s.serviceMock = serviceMocks.NewOrderService(s.T())
	s.api = NewAPI(s.serviceMock)
}

func (s *GetOrderAPISuite) TestGetOrder_MethodNotAllowed() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/order-1", nil)
	rr := httptest.NewRecorder()

	s.api.GetOrder(rr, req)

	require.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
	assert.JSONEq(s.T(), `{"error":"method not allowed"}`, rr.Body.String())
}

func (s *GetOrderAPISuite) TestGetOrder_InvalidOrderUUID() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/", nil)
	rr := httptest.NewRecorder()

	s.api.GetOrder(rr, req)

	require.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.JSONEq(s.T(), `{"error":"invalid order uuid"}`, rr.Body.String())
}

func (s *GetOrderAPISuite) TestGetOrder_Success() {
	orderUUID := "order-1"
	txUUID := "tx-1"
	method := model.PaymentMethodCard

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        "user-1",
		PartUUIDs:       []string{"p-1", "p-2"},
		TotalPrice:      150,
		TransactionUUID: &txUUID,
		PaymentMethod:   &method,
		Status:          model.StatusPaid,
	}

	s.serviceMock.EXPECT().GetOrder(mock.Anything, orderUUID).Return(order, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+orderUUID, nil)
	rr := httptest.NewRecorder()

	s.api.GetOrder(rr, req)

	require.Equal(s.T(), http.StatusOK, rr.Code)

	var got map[string]any
	require.NoError(s.T(), json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(s.T(), orderUUID, got["order_uuid"])
	assert.Equal(s.T(), "user-1", got["user_uuid"])
	assert.Equal(s.T(), float64(150), got["total_price"])
	assert.Equal(s.T(), string(model.StatusPaid), got["status"])
	assert.Equal(s.T(), txUUID, got["transaction_uuid"])
	assert.Equal(s.T(), string(method), got["payment_method"])
}

func (s *GetOrderAPISuite) TestGetOrder_NotFound() {
	orderUUID := "missing"
	s.serviceMock.EXPECT().GetOrder(mock.Anything, orderUUID).Return(model.Order{}, model.ErrOrderNotFound).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+orderUUID, nil)
	rr := httptest.NewRecorder()

	s.api.GetOrder(rr, req)

	require.Equal(s.T(), http.StatusNotFound, rr.Code)
	assert.JSONEq(s.T(), `{"error":"order not found"}`, rr.Body.String())
}

func (s *GetOrderAPISuite) TestGetOrder_ServiceError() {
	orderUUID := "order-1"
	errBoom := errors.New("upstream failed")
	s.serviceMock.EXPECT().GetOrder(mock.Anything, orderUUID).Return(model.Order{}, errBoom).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+orderUUID, nil)
	rr := httptest.NewRecorder()

	s.api.GetOrder(rr, req)

	require.Equal(s.T(), http.StatusBadGateway, rr.Code)
	assert.JSONEq(s.T(), `{"error":"upstream failed"}`, rr.Body.String())
}

func TestGetOrderAPISuite(t *testing.T) {
	suite.Run(t, new(GetOrderAPISuite))
}
