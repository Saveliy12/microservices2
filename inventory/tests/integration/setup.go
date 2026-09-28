package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
	"github.com.Saveliy12/microservices2.git/platform/pkg/testcontainers/app"
	"github.com.Saveliy12/microservices2.git/platform/pkg/testcontainers/mongo"
	"github.com.Saveliy12/microservices2.git/platform/pkg/testcontainers/network"
	"github.com.Saveliy12/microservices2.git/platform/pkg/testcontainers/path"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

// setupTestEnvironment поднимает сеть, MongoDB и контейнер inventory-сервиса
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(mongoContainerName),
		mongo.WithImageName(mongoImage),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithAuthDB(mongoAuthDB),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер MongoDB запущен")

	projectRoot := path.GetProjectRoot()

	appEnv, err := loadAppEnv(projectRoot, generatedMongo.Config().ContainerName)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "не удалось загрузить переменные окружения", zap.Error(err))
	}

	appContainer, err := app.NewContainer(ctx,
		app.WithName(appName),
		app.WithPort(appPort),
		app.WithDockerfile(projectRoot, appDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(wait.ForListeningPort(appPort+"/tcp").WithStartupTimeout(startupTimeout)),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")

	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

func loadAppEnv(projectRoot, mongoContainer string) (map[string]string, error) {
	envPath := filepath.Join(projectRoot, "deploy", "compose", "inventory", ".env")

	envVars, err := godotenv.Read(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", envPath, err)
	}

	envVars["MONGO_URI"] = fmt.Sprintf(
		"mongodb://%s:%s@%s:27017/%s?authSource=%s",
		mongoUsername, mongoPassword, mongoContainer, mongoDatabase, mongoAuthDB,
	)

	return envVars, nil
}
