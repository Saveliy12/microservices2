package service

import (
	"context"

	"github.com.Saveliy12/microservices2.git/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, payment model.Payment) (transactionUUID string, err error)
}
