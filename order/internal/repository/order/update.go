package order

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repositoryConverter "github.com.Saveliy12/microservices2.git/order/internal/repository/converter"
)

func (r *Repository) UpdateOrder(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.OrderUUID]; !ok {
		return model.ErrOrderNotFound
	}

	r.orders[order.OrderUUID] = repositoryConverter.OrderToRepoModel(order)
	return nil
}
