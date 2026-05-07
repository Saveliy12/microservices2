package order

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
)

func (s *Service) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (model.Order, error) {
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

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.OrderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		return model.Order{}, err
	}

	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &paymentMethod
	order.Status = model.StatusPaid

	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}
