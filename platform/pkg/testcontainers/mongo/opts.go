package mongo

// Option — функциональная опция для настройки Config
type Option func(*Config)

func WithNetworkName(network string) Option {
	return func(c *Config) {
		c.NetworkName = network
	}
}

func WithContainerName(containerName string) Option {
	return func(c *Config) {
		c.ContainerName = containerName
	}
}

func WithImageName(image string) Option {
	return func(c *Config) {
		c.ImageName = image
	}
}

func WithDatabase(database string) Option {
	return func(c *Config) {
		c.Database = database
	}
}

// WithAuth задаёт root-пользователя MongoDB
func WithAuth(username, password string) Option {
	return func(c *Config) {
		c.Username = username
		c.Password = password
	}
}

// WithAuthDB задаёт базу для аутентификации (обычно admin)
func WithAuthDB(authDB string) Option {
	return func(c *Config) {
		c.AuthDB = authDB
	}
}

// WithLogger подключает логгер для сообщений обёртки
func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}
