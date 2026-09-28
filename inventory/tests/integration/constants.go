package integration

import "time"

const (
	projectName = "inventory-service"

	appName       = "inventory-app-test"
	appPort       = "50051"
	appDockerfile = "inventory/Dockerfile"

	// Параметры тестового MongoDB
	// Имя контейнера отличается от compose-контейнера (mongo-inventory),
	// чтобы тесты не конфликтовали с запущенным docker compose стеком
	mongoContainerName = "mongo-inventory-test"
	mongoImage         = "mongo:7.0.5"
	mongoDatabase      = "example"
	mongoUsername      = "demo"
	mongoPassword      = "demo" //nolint:gosec
	mongoAuthDB        = "admin"

	partsCollectionName = "parts"

	testsTimeout   = 5 * time.Minute
	startupTimeout = 3 * time.Minute

	loggerLevelValue = "debug"
)
