package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
)

const (
	defaultAppName        = "app"
	defaultAppPort        = "50051"
	defaultStartupTimeout = 1 * time.Minute
)

// Logger - минимальный интерфейс логгера
type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// Config — параметры контейнера сервиса
type Config struct {
	Name          string            // имя контейнера (host внутри сети)
	DockerfileDir string            // контекст сборки (обычно корень проекта)
	Dockerfile    string            // путь к Dockerfile относительно контекста
	Port          string            // порт, который слушает сервис
	Env           map[string]string // переменные окружения сервиса
	Networks      []string          // сети, к которым подключить контейнер
	LogOutput     io.Writer         // куда стримить логи сборки и работы
	StartupWait   wait.Strategy     // стратегия ожидания готовности
	Logger        Logger
}

// Container — запущенный контейнер сервиса
type Container struct {
	container    testcontainers.Container
	externalHost string
	externalPort string
	cfg          *Config
}

// NewContainer собирает образ из Dockerfile, запускает контейнер в нужной сети
// и дожидается готовности сервиса (по умолчанию — открытый порт)
func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := &Config{
		Name:          defaultAppName,
		Port:          defaultAppPort,
		Dockerfile:    "Dockerfile",
		DockerfileDir: ".",
		LogOutput:     io.Discard,
		StartupWait:   wait.ForListeningPort(defaultAppPort + "/tcp").WithStartupTimeout(defaultStartupTimeout),
		Env:           make(map[string]string),
		Logger:        &logger.NoopLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	req := testcontainers.ContainerRequest{
		Name: cfg.Name,
		FromDockerfile: testcontainers.FromDockerfile{
			Context:        cfg.DockerfileDir,
			Dockerfile:     cfg.Dockerfile,
			BuildLogWriter: cfg.LogOutput,
		},
		Networks:           cfg.Networks,
		Env:                cfg.Env,
		WaitingFor:         cfg.StartupWait,
		ExposedPorts:       []string{cfg.Port + "/tcp"},
		HostConfigModifier: DefaultHostConfig(),
	}

	genericContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start app genericContainer: %w", err)
	}

	// Порт сервиса мапится на случайный порт хоста — по нему тесты и подключаются
	mappedPort, err := genericContainer.MappedPort(ctx, cfg.Port+"/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped externalPort: %w", err)
	}

	host, err := genericContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get genericContainer externalHost: %w", err)
	}

	go streamContainerLogs(ctx, genericContainer, cfg.LogOutput)

	cfg.Logger.Info(ctx, "App container started", zap.String("uri:", net.JoinHostPort(host, mappedPort.Port())))

	return &Container{
		container:    genericContainer,
		externalHost: host,
		externalPort: mappedPort.Port(),
		cfg:          cfg,
	}, nil
}

// Address возвращает host:port, доступный с хоста, — для gRPC/HTTP клиента в тестах
func (a *Container) Address() string {
	return net.JoinHostPort(a.externalHost, a.externalPort)
}

// Terminate останавливает и удаляет контейнер сервиса
func (a *Container) Terminate(ctx context.Context) error {
	return a.container.Terminate(ctx)
}

// streamContainerLogs переливает логи контейнера в переданный writer (обычно os.Stdout)
func streamContainerLogs(ctx context.Context, container testcontainers.Container, out io.Writer) {
	logs, err := container.Logs(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get container logs", zap.Error(err))
		return
	}
	defer func() {
		err = logs.Close()
		if err != nil {
			logger.Error(ctx, "failed to close container logs", zap.Error(err))
		}
	}()

	go func() {
		_, err = io.Copy(out, logs)
		if err != nil && !errors.Is(err, io.EOF) {
			logger.Error(ctx, "error copying container logs", zap.Error(err))
		}
	}()
}

// DefaultHostConfig удаляет контейнер при остановке, чтобы не копить мусор
func DefaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}
