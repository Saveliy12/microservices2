package testcontainers

// MongoDB constants
const (
	// MongoDB container constants
	MongoContainerName = "mongo-inventory" // имя контейнера = host в MONGO_URI
	MongoPort          = "27017"           // внутренний порт MongoDB

	// MongoDB environment variables — ключи, по которым тесты читают параметры из окружения
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_DB"
)
