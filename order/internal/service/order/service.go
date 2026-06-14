package order

import (
	repository "github.com.Saveliy12/microservices2.git/order/internal/repository"
	service "github.com.Saveliy12/microservices2.git/order/internal/service"
)

var _ service.OrderService = (*Service)(nil)

type Service struct {
	orderRepository repository.OrderRepository
	inventoryClient service.InventoryClient
	paymentClient   service.PaymentClient
}

func NewService(orderRepository repository.OrderRepository, inventoryClient service.InventoryClient, paymentClient service.PaymentClient) *Service {
	return &Service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
