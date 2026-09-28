package config

import (
	"github.com.Saveliy12/microservices2.git/inventory/internal/config/env"
)

var _ InventoryConfig = (*config)(nil)

type config struct {
	grpcConf   GRPCConfig
	loggerConf LoggerConfig
	mongoConf  MongoConfig
}

// New собирает конфигурацию inventory-сервиса из переменных окружения
func New() (*config, error) {
	grpcConf, err := env.NewGRPCConfig()
	if err != nil {
		return nil, err
	}

	loggerConf, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	MongoConf, err := env.NewMongoConfig()
	if err != nil {
		return nil, err
	}

	return &config{
		grpcConf:   grpcConf,
		loggerConf: loggerConf,
		mongoConf:  MongoConf,
	}, nil
}

func (cfg *config) GRPC() GRPCConfig {
	return cfg.grpcConf
}

func (cfg *config) Logger() LoggerConfig {
	return cfg.loggerConf
}

func (cfg *config) Mongo() MongoConfig {
	return cfg.mongoConf
}
