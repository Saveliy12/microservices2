package env

import (
	"os"
	"strconv"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
	grpcPortDefault = 50051
)

type grpcConfig struct {
	host string
	port int
}

// NewGRPCConfig считывает настройки gRPC сервера из переменных окружения
func NewGRPCConfig() (*grpcConfig, error) {
	portStr := os.Getenv(grpcPortEnvName)
	if portStr == "" {
		return &grpcConfig{port: grpcPortDefault}, nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &grpcConfig{host: os.Getenv(grpcHostEnvName), port: port}, nil
}

func (cfg *grpcConfig) Host() string {
	return cfg.host
}

func (cfg *grpcConfig) Port() int {
	return cfg.port
}
