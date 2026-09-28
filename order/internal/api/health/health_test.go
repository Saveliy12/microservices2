package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com.Saveliy12/microservices2.git/order/internal/api/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HealthHandlerTestSuite struct {
	suite.Suite
	handler *health.Handler
}

func (s *HealthHandlerTestSuite) SetupTest() {
	s.handler = health.NewHandler()
}

func (s *HealthHandlerTestSuite) TestHealth_Get() {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	s.handler.ServeHTTP(rr, req)

	require.Equal(s.T(), http.StatusOK, rr.Code)
	assert.JSONEq(s.T(), `{"status":"SERVING"}`, rr.Body.String())
	assert.Equal(s.T(), "application/json", rr.Header().Get("Content-Type"))
}

func (s *HealthHandlerTestSuite) TestHealth_AnyMethod() {
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rr := httptest.NewRecorder()

	s.handler.ServeHTTP(rr, req)

	require.Equal(s.T(), http.StatusOK, rr.Code)
	assert.JSONEq(s.T(), `{"status":"SERVING"}`, rr.Body.String())
}

func TestHealthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HealthHandlerTestSuite))
}
