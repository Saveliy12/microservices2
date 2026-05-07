package converter

import (
	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repositoryModel "github.com.Saveliy12/microservices2.git/order/internal/repository/model"
	"github.com/samber/lo"
)

func OrderToRepoModel(order model.Order) repositoryModel.Order {
	partUUIDs := make([]string, len(order.PartUUIDs))
	copy(partUUIDs, order.PartUUIDs)

	return repositoryModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: cloneStringPtr(order.TransactionUUID),
		PaymentMethod:   toRepoPaymentMethodPtr(order.PaymentMethod),
		Status:          repositoryModel.Status(order.Status),
	}
}

func OrderFromRepoModel(order repositoryModel.Order) model.Order {
	partUUIDs := make([]string, len(order.PartUUIDs))
	copy(partUUIDs, order.PartUUIDs)

	return model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: cloneStringPtr(order.TransactionUUID),
		PaymentMethod:   fromRepoPaymentMethodPtr(order.PaymentMethod),
		Status:          model.Status(order.Status),
	}
}

func cloneStringPtr(v *string) *string {
	if v == nil {
		return nil
	}

	return lo.ToPtr(*v)
}

func toRepoPaymentMethodPtr(v *model.PaymentMethod) *repositoryModel.PaymentMethod {
	if v == nil {
		return nil
	}

	return lo.ToPtr(repositoryModel.PaymentMethod(*v))
}

func fromRepoPaymentMethodPtr(v *repositoryModel.PaymentMethod) *model.PaymentMethod {
	if v == nil {
		return nil
	}

	return lo.ToPtr(model.PaymentMethod(*v))
}
