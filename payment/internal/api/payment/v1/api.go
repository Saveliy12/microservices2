package v1

import (
	paymentservice "github.com.Saveliy12/microservices2.git/payment/internal/service"
	paymentv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/payment/v1"
)

var _ paymentv1.PaymentServiceServer = (*API)(nil)

type API struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService paymentservice.PaymentService
}

func NewAPI(paymentService paymentservice.PaymentService) *API {
	return &API{
		paymentService: paymentService,
	}
}
