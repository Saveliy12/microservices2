package converter

import (
	"github.com.Saveliy12/microservices2.git/payment/internal/model"
	paymentv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/payment/v1"
)

// api-модель -> service-модель
func PayOrderRequestToModel(req *paymentv1.PayOrderRequest) model.Payment {
	if req == nil {
		return model.Payment{}
	}

	return model.Payment{
		OrderUUID:     req.GetOrderUuid(),
		UserUUID:      req.GetUserUuid(),
		PaymentMethod: req.GetPaymentMethod().String(),
	}
}

// service-модель -> api-модель (для uuid)
func TransactionUUIDToPayOrderResponse(transactionUUID string) *paymentv1.PayOrderResponse {
	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}
}
