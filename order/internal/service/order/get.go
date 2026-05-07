package order

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
)

func (s *Service) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	return s.orderRepository.GetOrder(ctx, orderUUID)
}
