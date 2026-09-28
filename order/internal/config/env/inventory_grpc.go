package env

import (
	"os"
	"strconv"
)

const (
	inventoryGRPCHostEnvName = "INVENTORY_GRPC_HOST"
	inventoryGRPCPortEnvName = "INVENTORY_GRPC_PORT"
	inventoryGRPCPortDefault = 50051
)

type inventoryGRPCConfig struct {
	host string
	port int
}

// NewInventoryGRPCConfig считывает адрес Inventory gRPC из переменных окружения
func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	portStr := os.Getenv(inventoryGRPCPortEnvName)
	if portStr == "" {
		return &inventoryGRPCConfig{port: inventoryGRPCPortDefault}, nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &inventoryGRPCConfig{
		host: os.Getenv(inventoryGRPCHostEnvName),
		port: port,
	}, nil
}

func (cfg *inventoryGRPCConfig) Host() string {
	return cfg.host
}

func (cfg *inventoryGRPCConfig) Port() int {
	return cfg.port
}
