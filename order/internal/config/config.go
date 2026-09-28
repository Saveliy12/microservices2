package config

import (
	"github.com.Saveliy12/microservices2.git/order/internal/config/env"
)

var _ OrderConfig = (*config)(nil)

type config struct {
	httpConf          HTTPConfig
	inventoryGRPCConf GRPCConfig
	paymentGRPCConf   GRPCConfig
	postgresConf      PostgresConfig
	loggerConf        LoggerConfig
}

func New() (*config, error) {
	httpConf, err := env.NewHTTPConfig()
	if err != nil {
		return nil, err
	}

	inventoryGRPCConf, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return nil, err
	}

	paymentGRPCConf, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return nil, err
	}

	postgresConf, err := env.NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	loggerConf, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	return &config{
		httpConf:          httpConf,
		inventoryGRPCConf: inventoryGRPCConf,
		paymentGRPCConf:   paymentGRPCConf,
		postgresConf:      postgresConf,
		loggerConf:        loggerConf,
	}, nil
}

func (c *config) HTTP() HTTPConfig {
	return c.httpConf
}

func (c *config) InventoryGRPC() GRPCConfig {
	return c.inventoryGRPCConf
}

func (c *config) PaymentGRPC() GRPCConfig {
	return c.paymentGRPCConf
}

func (c *config) Postgres() PostgresConfig {
	return c.postgresConf
}

func (c *config) Logger() LoggerConfig {
	return c.loggerConf
}
