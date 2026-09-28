package app

import (
	apiv1 "github.com.Saveliy12/microservices2.git/payment/internal/api/payment/v1"
	"github.com.Saveliy12/microservices2.git/payment/internal/config"
	paymentservice "github.com.Saveliy12/microservices2.git/payment/internal/service"
	paymentimpl "github.com.Saveliy12/microservices2.git/payment/internal/service/payment"
	"github.com.Saveliy12/microservices2.git/platform/pkg/closer"
	"github.com.Saveliy12/microservices2.git/platform/pkg/grpc/health"
	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
	paymentv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type serviceProvider struct {
	config     config.PaymentConfig
	service    paymentservice.PaymentService
	api        *apiv1.API
	grpcServer *grpc.Server
	closer     *closer.Closer
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) Config() (config.PaymentConfig, error) {
	if sp.config == nil {
		cfg, err := config.New()
		if err != nil {
			return nil, err
		}

		sp.config = cfg
	}

	return sp.config, nil
}

func (sp *serviceProvider) PaymentService() paymentservice.PaymentService {
	if sp.service == nil {
		sp.service = paymentimpl.NewService()
	}

	return sp.service
}

func (sp *serviceProvider) PaymentAPI() *apiv1.API {
	if sp.api == nil {
		sp.api = apiv1.NewAPI(sp.PaymentService())
	}

	return sp.api
}

func (sp *serviceProvider) GRPCServer() *grpc.Server {
	if sp.grpcServer == nil {
		server := grpc.NewServer()
		paymentv1.RegisterPaymentServiceServer(server, sp.PaymentAPI())
		health.RegisterService(server)
		reflection.Register(server)

		sp.grpcServer = server
	}

	return sp.grpcServer
}

func (sp *serviceProvider) Closer() *closer.Closer {
	if sp.closer == nil {
		sp.closer = closer.NewWithLogger(logger.Logger())
	}

	return sp.closer
}
