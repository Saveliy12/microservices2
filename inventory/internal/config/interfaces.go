package config

type InventoryConfig interface {
	GRPC() GRPCConfig
	Logger() LoggerConfig
	Mongo() MongoConfig
}

type GRPCConfig interface {
	Host() string
	Port() int
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type MongoConfig interface {
	URI() string
	Database() string
}
