package v1

import (
	"context"
	"testing"

	serviceMocks "github.com.Saveliy12/microservices2.git/inventory/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type InventoryAPISuite struct {
	suite.Suite
	ctx           context.Context
	api           *API
	inventoryMock *serviceMocks.InventoryService
}

func (s *InventoryAPISuite) SetupTest() {
	s.ctx = context.Background()
	s.inventoryMock = serviceMocks.NewInventoryService(s.T())
	s.api = NewAPI(s.inventoryMock)
}

func TestInventoryAPISuite(t *testing.T) {
	suite.Run(t, new(InventoryAPISuite))
}
