package payment

import (
	"context"
	"errors"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	service "github.com.Saveliy12/microservices2.git/order/internal/service"
	paymentv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

var _ service.PaymentClient = (*Client)(nil)

type Client struct {
	client paymentv1.PaymentServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{
		client: paymentv1.NewPaymentServiceClient(conn),
	}
}

func (c *Client) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	pmEnum, err := paymentMethodToProto(paymentMethod)
	if err != nil {
		return "", err
	}

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: pmEnum,
	}

	resp, err := c.client.PayOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetTransactionUuid(), nil
}

// paymentMethodToProto преобразует доменный enum PaymentMethod в protobuf enum PaymentMethod
func paymentMethodToProto(method model.PaymentMethod) (paymentv1.PaymentMethod, error) {
	switch method {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case model.PaymentMethodSbp:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, errors.New("invalid payment method")
	}
}
