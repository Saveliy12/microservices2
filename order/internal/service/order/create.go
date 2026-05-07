package order

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com/google/uuid"
)

func (s *Service) CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (model.Order, error) {
	if userUUID == "" || len(partUUIDs) == 0 {
		return model.Order{}, model.ErrInvalidOrderInput
	}

	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		return model.Order{}, err
	}

	partsByUUID := make(map[string]model.PartQuote, len(parts))
	for _, part := range parts {
		partsByUUID[part.UUID] = part
	}

	if len(partsByUUID) != len(partUUIDs) {
		return model.Order{}, model.ErrInvalidOrderInput
	}

	totalPrice := calculateTotalPrice(partUUIDs, partsByUUID)
	order := model.Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   userUUID,
		PartUUIDs:  append([]string(nil), partUUIDs...),
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	if err := s.orderRepository.CreateOrder(ctx, order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func calculateTotalPrice(partUUIDs []string, partsByUUID map[string]model.PartQuote) float64 {
	totalPrice := 0.0
	for _, partUUID := range partUUIDs {
		totalPrice += partsByUUID[partUUID].Price
	}

	return totalPrice
}
