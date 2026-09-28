package config

import (
	"github.com.Saveliy12/microservices2.git/payment/internal/config/env"
)

var _ PaymentConfig = (*config)(nil)

type config struct {
	grpcConf   GRPCConfig
	loggerConf LoggerConfig
}

// New собирает конфигурацию payment-сервиса из переменных окружения
func New() (*config, error) {
	grpcConf, err := env.NewGRPCConfig()
	if err != nil {
		return nil, err
	}

	loggerConf, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	return &config{
		grpcConf:   grpcConf,
		loggerConf: loggerConf,
	}, nil
}

func (cfg *config) GRPC() GRPCConfig {
	return cfg.grpcConf
}

func (cfg *config) Logger() LoggerConfig {
	return cfg.loggerConf
}
