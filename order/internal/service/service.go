package service

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (model.Order, error)
	GetOrder(ctx context.Context, orderUUID string) (model.Order, error)
	PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (model.Order, error)
	CancelOrder(ctx context.Context, orderUUID string) (model.Order, error)
}

type InventoryClient interface {
	ListParts(ctx context.Context, partUUIDs []string) ([]model.PartQuote, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error)
}
