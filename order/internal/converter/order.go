package converter

import (
	"fmt"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
)

type GetOrderResponse struct {
	OrderUUID       string   `json:"order_uuid"`
	UserUUID        string   `json:"user_uuid"`
	PartUUIDs       []string `json:"part_uuids"`
	TotalPrice      float64  `json:"total_price"`
	TransactionUUID *string  `json:"transaction_uuid,omitempty"`
	PaymentMethod   *string  `json:"payment_method,omitempty"`
	Status          string   `json:"status"`
}

type PayOrderResponse struct {
	TransactionUUID string `json:"transaction_uuid"`
}

func ModelToGetResponse(o model.Order) GetOrderResponse {
	status := string(o.Status)
	return GetOrderResponse{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUUIDs:       o.PartUUIDs,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   paymentMethodToStringPtr(o.PaymentMethod),
		Status:          status,
	}
}

func ModelToPayResponse(o model.Order) (PayOrderResponse, error) {
	if o.TransactionUUID == nil || *o.TransactionUUID == "" {
		return PayOrderResponse{}, fmt.Errorf("transaction uuid is empty")
	}

	return PayOrderResponse{TransactionUUID: *o.TransactionUUID}, nil
}

func PaymentMethodFromString(value string) (model.PaymentMethod, error) {
	switch value {
	case "CARD", "PAYMENT_METHOD_CARD":
		return model.PaymentMethodCard, nil
	case "SBP", "PAYMENT_METHOD_SBP":
		return model.PaymentMethodSbp, nil
	case "CREDIT_CARD", "PAYMENT_METHOD_CREDIT_CARD":
		return model.PaymentMethodCreditCard, nil
	case "INVESTOR_MONEY", "PAYMENT_METHOD_INVESTOR_MONEY":
		return model.PaymentMethodInvestorMoney, nil
	default:
		return model.PaymentMethodUnspecified, fmt.Errorf("invalid payment_method")
	}
}

func paymentMethodToStringPtr(method *model.PaymentMethod) *string {
	if method == nil {
		return nil
	}

	v := string(*method)
	return &v
}
