package config

type OrderConfig interface {
	HTTP() HTTPConfig          // :8080
	InventoryGRPC() GRPCConfig // localhost:50051
	PaymentGRPC() GRPCConfig   // localhost:50052
	Postgres() PostgresConfig  // DB_URI
	Logger() LoggerConfig
}

type HTTPConfig interface {
	Host() string
	Port() int
}
type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}

type GRPCConfig interface {
	Host() string
	Port() int
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}
