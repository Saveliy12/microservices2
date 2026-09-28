package app

import (
	"io"
	"maps"

	"github.com/testcontainers/testcontainers-go/wait"
)

// Option — функциональная опция для настройки Config
type Option func(*Config)

// WithName задаёт имя контейнера (оно же host внутри Docker-сети)
func WithName(name string) Option {
	return func(c *Config) {
		c.Name = name
	}
}

// WithDockerfile задаёт контекст сборки и путь к Dockerfile относительно него
func WithDockerfile(dir, file string) Option {
	return func(c *Config) {
		c.DockerfileDir = dir
		c.Dockerfile = file
	}
}

// WithPort задаёт порт, который слушает сервис внутри контейнера
func WithPort(port string) Option {
	return func(c *Config) {
		c.Port = port
	}
}

// WithNetwork подключает контейнер к Docker-сети (можно вызывать несколько раз)
func WithNetwork(name string) Option {
	return func(c *Config) {
		c.Networks = append(c.Networks, name)
	}
}

// WithEnv добавляет переменные окружения сервису
func WithEnv(env map[string]string) Option {
	return func(c *Config) {
		maps.Copy(c.Env, env)
	}
}

// WithLogOutput задаёт writer для логов сборки и работы контейнера
func WithLogOutput(out io.Writer) Option {
	return func(c *Config) {
		c.LogOutput = out
	}
}

// WithStartupWait переопределяет стратегию ожидания готовности сервиса
func WithStartupWait(strategy wait.Strategy) Option {
	return func(c *Config) {
		c.StartupWait = strategy
	}
}

// WithLogger подключает логгер для сообщений обёртки.
func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}
