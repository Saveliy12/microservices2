package payment

import (
	"context"
	"log"

	"github.com.Saveliy12/microservices2.git/payment/internal/model"
	"github.com/google/uuid"
)

func (*Service) PayOrder(_ context.Context, payment model.Payment) (string, error) {

	transactionUUID := uuid.NewString()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return transactionUUID, nil
}
