package integration

import (
	"context"

	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
	"go.uber.org/zap"
)

// teardownTestEnvironment освобождает все ресурсы тестового окружения
func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	logger.Info(ctx, "🧹 Очистка тестового окружения...")

	cleanupTestEnvironment(ctx, env)

	logger.Info(ctx, "✅ Тестовое окружение успешно очищено")
}

// cleanupTestEnvironment закрывает пул соединений и удаляет контейнер
func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.Pool != nil {
		env.Pool.Close()
	}

	if env.Container != nil {
		if err := env.Container.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер PostgreSQL", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер PostgreSQL остановлен")
		}
	}
}
