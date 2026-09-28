package env

import (
	"os"
	"strconv"
)

const (
	paymentGRPCHostEnvName = "PAYMENT_GRPC_HOST"
	paymentGRPCPortEnvName = "PAYMENT_GRPC_PORT"
	paymentGRPCPortDefault = 50052
)

type paymentGRPCConfig struct {
	host string
	port int
}

// NewPaymentGRPCConfig считывает адрес Payment gRPC из переменных окружения
func NewPaymentGRPCConfig() (*paymentGRPCConfig, error) {
	portStr := os.Getenv(paymentGRPCPortEnvName)
	if portStr == "" {
		return &paymentGRPCConfig{port: paymentGRPCPortDefault}, nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &paymentGRPCConfig{
		host: os.Getenv(paymentGRPCHostEnvName),
		port: port,
	}, nil
}

func (cfg *paymentGRPCConfig) Host() string {
	return cfg.host
}

func (cfg *paymentGRPCConfig) Port() int {
	return cfg.port
}
