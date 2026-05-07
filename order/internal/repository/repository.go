package repository

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
)

type OrderRepository interface {
	GetOrder(ctx context.Context, orderUUID string) (model.Order, error)
	CreateOrder(ctx context.Context, order model.Order) error
	UpdateOrder(ctx context.Context, order model.Order) error
}
