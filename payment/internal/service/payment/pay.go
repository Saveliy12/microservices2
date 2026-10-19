package payment

import (
	"context"

	"github.com.Saveliy12/microservices2.git/payment/internal/model"
	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (*Service) PayOrder(ctx context.Context, _ model.Payment) (string, error) {
	transactionUUID := uuid.NewString()
	logger.Info(ctx, "Оплата успешно проведена", zap.String("transaction_uuid", transactionUUID))

	return transactionUUID, nil
}
