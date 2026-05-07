package part

import (
	repository "github.com.Saveliy12/microservices2.git/inventory/internal/repository"
	service "github.com.Saveliy12/microservices2.git/inventory/internal/service"
)

var _ service.InventoryService = (*Service)(nil)

type Service struct {
	partRepository repository.PartRepository
}

func NewService(partRepository repository.PartRepository) *Service {
	return &Service{
		partRepository: partRepository,
	}
}
