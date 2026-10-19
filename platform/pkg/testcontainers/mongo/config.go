package mongo

import (
	"context"

	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
	"github.com/moby/moby/api/types/container"
	"go.uber.org/zap"
)

// Logger — минимальный интерфейс логгера
type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// Config — параметры MongoDB-контейнера
// Host и Port заполняются после старта контейнера (Docker мапит 27017 на случайный порт хоста)
type Config struct {
	NetworkName   string
	ContainerName string
	ImageName     string
	Database      string
	Username      string
	Password      string
	AuthDB        string
	Logger        Logger

	Host string
	Port string
}

// buildConfig собирает конфиг из дефолтов и применённых опций
func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		NetworkName:   "test-network",
		ContainerName: "mongo-inventory",
		ImageName:     "mongo:7.0.5",
		Database:      "example",
		Username:      "demo",
		Password:      "demo",
		AuthDB:        "admin",
		Logger:        &logger.NoopLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// defaultHostConfig удаляет контейнер при остановке, чтобы не копить мусор
func defaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}
