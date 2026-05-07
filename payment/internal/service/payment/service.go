package payment

import (
	paymentservice "github.com.Saveliy12/microservices2.git/payment/internal/service"
)

var _ paymentservice.PaymentService = (*Service)(nil)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
