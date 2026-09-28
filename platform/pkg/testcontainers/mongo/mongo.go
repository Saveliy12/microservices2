package mongo

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

const (
	mongoPort           = "27017"
	mongoStartupTimeout = 1 * time.Minute

	mongoEnvUsernameKey = "MONGO_INITDB_ROOT_USERNAME"
	mongoEnvPasswordKey = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
)

// Container — запущенный MongoDB-контейнер вместе с готовым клиентом
type Container struct {
	container testcontainers.Container
	client    *mongo.Client
	cfg       *Config
}

// NewContainer поднимает контейнер, узнаёт его адрес, подключает клиент и проверяет Ping
// При любой ошибке контейнер удаляется, чтобы не оставлять мусор
func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	container, err := startMongoContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	uri := buildMongoURI(cfg)

	client, err := connectMongoClient(ctx, uri)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info(ctx, "Mongo container started", zap.String("uri", uri))
	success = true

	return &Container{
		container: container,
		client:    client,
		cfg:       cfg,
	}, nil
}

// Client отдаёт mongo-клиент для подготовки тестовых данных
func (c *Container) Client() *mongo.Client {
	return c.client
}

// Config отдаёт параметры контейнера
func (c *Container) Config() *Config {
	return c.cfg
}

// Terminate закрывает клиент и удаляет контейнер
func (c *Container) Terminate(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to disconnect mongo client", zap.Error(err))
	}

	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
	}

	c.cfg.Logger.Info(ctx, "Mongo container terminated")

	return nil
}
