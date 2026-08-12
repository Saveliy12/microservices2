package v1

import (
	"testing"

	serviceMocks "github.com.Saveliy12/microservices2.git/order/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type OrderAPISuite struct {
	suite.Suite
	api         *API
	serviceMock *serviceMocks.OrderServiceMock
}

func (s *OrderAPISuite) SetupTest() {
	s.serviceMock = serviceMocks.NewOrderServiceMock(s.T())
	s.api = NewAPI(s.serviceMock)
}

func TestOrderAPISuite(t *testing.T) {
	suite.Run(t, new(OrderAPISuite))
}
