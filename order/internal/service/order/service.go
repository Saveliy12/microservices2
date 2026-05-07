package order

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repository "github.com.Saveliy12/microservices2.git/order/internal/repository"
	service "github.com.Saveliy12/microservices2.git/order/internal/service"
)

type InventoryClient interface {
	ListParts(ctx context.Context, partUUIDs []string) ([]model.PartQuote, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error)
}

var _ service.OrderService = (*Service)(nil)

type Service struct {
	orderRepository repository.OrderRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
}

func NewService(orderRepository repository.OrderRepository, inventoryClient InventoryClient, paymentClient PaymentClient) *Service {
	return &Service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
