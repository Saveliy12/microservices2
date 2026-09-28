// Package network создаёт изолированную Docker-сеть для тестовых контейнеров
// Внутри одной сети контейнеры видят друг друга по именам
package network

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	tcnetwork "github.com/testcontainers/testcontainers-go/network"
)

// Network — обёртка над Docker-сетью, созданной testcontainers
type Network struct {
	network *testcontainers.DockerNetwork
}

// NewNetwork создаёт новую bridge-сеть с меткой проекта
// Имя сети генерируется автоматически и уникально для каждого запуска тестов
func NewNetwork(ctx context.Context, projectName string) (*Network, error) {
	net, err := tcnetwork.New(ctx,
		tcnetwork.WithDriver(testcontainers.Bridge),
		tcnetwork.WithAttachable(),
		tcnetwork.WithLabels(map[string]string{
			"project": projectName,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker network: %w", err)
	}

	return &Network{network: net}, nil
}

// Name возвращает сгенерированное имя сети — его передают контейнерам
func (n *Network) Name() string {
	return n.network.Name
}

// Remove удаляет сеть после завершения тестов.
func (n *Network) Remove(ctx context.Context) error {
	return n.network.Remove(ctx)
}
