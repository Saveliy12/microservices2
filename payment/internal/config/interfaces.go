package config

type PaymentConfig interface {
	GRPC() GRPCConfig
	Logger() LoggerConfig
}

type GRPCConfig interface {
	Host() string
	Port() int
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}
