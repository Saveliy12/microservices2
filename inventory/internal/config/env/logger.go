package env

import (
	"os"
	"strconv"
)

const (
	loggerLevelEnvName  = "LOGGER_LEVEL"
	loggerLevelDefault  = "info"
	loggerAsJSONEnvName = "LOGGER_AS_JSON"
	loggerAsJSONDefault = false
)

type loggerConfig struct {
	level  string
	asJSON bool
}

// NewLoggerConfig считывает настройки логгера из переменных окружения
func NewLoggerConfig() (*loggerConfig, error) {
	level := os.Getenv(loggerLevelEnvName)
	if level == "" {
		level = loggerLevelDefault
	}

	asJSONStr := os.Getenv(loggerAsJSONEnvName)
	if asJSONStr == "" {
		return &loggerConfig{level: level, asJSON: loggerAsJSONDefault}, nil
	}

	asJSON, err := strconv.ParseBool(asJSONStr)
	if err != nil {
		return nil, err
	}

	return &loggerConfig{level: level, asJSON: asJSON}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.level
}

func (cfg *loggerConfig) AsJSON() bool {
	return cfg.asJSON
}
