package v1

import (
	"context"

	"github.com.Saveliy12/microservices2.git/payment/internal/converter"
	paymentv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *API) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	transactionUUID, err := p.paymentService.PayOrder(ctx, converter.PayOrderRequestToModel(req))
	if err != nil {
		return nil, err
	}

	return converter.TransactionUUIDToPayOrderResponse(transactionUUID), nil
}
