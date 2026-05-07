package v1

import (
	service "github.com.Saveliy12/microservices2.git/inventory/internal/service"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
)

var _ inventoryv1.InventoryServiceServer = (*API)(nil)

type API struct {
	inventoryv1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewAPI(inventoryService service.InventoryService) *API {
	return &API{
		inventoryService: inventoryService,
	}
}
