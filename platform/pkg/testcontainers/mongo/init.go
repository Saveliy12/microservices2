package mongo

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// startMongoContainer создаёт и запускает контейнер MongoDB
// Ждём готовности по открытому порту 27017 (wait.ForListeningPort), а не по sleep
func startMongoContainer(ctx context.Context, cfg *Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:     cfg.ContainerName,
		Image:    cfg.ImageName,
		Networks: []string{cfg.NetworkName},
		Env: map[string]string{
			mongoEnvUsernameKey: cfg.Username,
			mongoEnvPasswordKey: cfg.Password,
		},
		WaitingFor:         wait.ForListeningPort(mongoPort + "/tcp").WithStartupTimeout(mongoStartupTimeout),
		HostConfigModifier: defaultHostConfig(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start mongo container: %w", err)
	}

	return container, nil
}

// getContainerHostPort возвращает адрес и порт, по которым контейнер доступен с хоста
// Docker мапит внутренний 27017 на случайный свободный порт — конфликтов с локальными сервисами нет
func getContainerHostPort(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, mongoPort+"/tcp")
	if err != nil {
		return "", "", fmt.Errorf("failed to get mapped port: %w", err)
	}

	return host, port.Port(), nil
}

// buildMongoURI собирает строку подключения для драйвера MongoDB
func buildMongoURI(cfg *Config) string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.AuthDB,
	)
}
