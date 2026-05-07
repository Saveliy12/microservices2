package order

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
)

func (s *Service) CancelOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}

	if order.Status == model.StatusPaid {
		return model.Order{}, model.ErrOrderAlreadyPaid
	}

	if order.Status == model.StatusCancelled {
		return model.Order{}, model.ErrOrderCancelled
	}

	order.Status = model.StatusCancelled
	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}
