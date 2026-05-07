package order

import (
	"sync"

	repository "github.com.Saveliy12/microservices2.git/order/internal/repository"
	repositoryModel "github.com.Saveliy12/microservices2.git/order/internal/repository/model"
)

var _ repository.OrderRepository = (*Repository)(nil)

type Repository struct {
	mu     sync.RWMutex
	orders map[string]repositoryModel.Order
}

func NewRepository() *Repository {
	return &Repository{
		orders: make(map[string]repositoryModel.Order),
	}
}
