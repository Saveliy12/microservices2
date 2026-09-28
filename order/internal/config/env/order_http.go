package env

import (
	"os"
	"strconv"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
	httpPortDefault = 8080
)

type httpConfig struct {
	host string
	port int
}

// NewHTTPConfig считывает HTTP настройки из переменных окружения
func NewHTTPConfig() (*httpConfig, error) {
	portStr := os.Getenv(httpPortEnvName)
	if portStr == "" {
		return &httpConfig{port: httpPortDefault}, nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &httpConfig{
		host: os.Getenv(httpHostEnvName),
		port: port,
	}, nil
}

func (cfg *httpConfig) Host() string {
	return cfg.host
}

func (cfg *httpConfig) Port() int {
	return cfg.port
}
