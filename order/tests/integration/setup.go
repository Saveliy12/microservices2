package integration

import (
	"context"
	"path/filepath"

	dbpg "github.com.Saveliy12/microservices2.git/order/internal/client/db/pg"
	repositoryorder "github.com.Saveliy12/microservices2.git/order/internal/repository/order"
	"github.com.Saveliy12/microservices2.git/platform/pkg/logger"
	migratorpg "github.com.Saveliy12/microservices2.git/platform/pkg/migrator/pg"
	"github.com.Saveliy12/microservices2.git/platform/pkg/testcontainers/path"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

// setupTestEnvironment поднимает PostgreSQL в контейнере, применяет миграции
// и собирает репозиторий — готовит настоящую БД для тестов
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Запускаем контейнер PostgreSQL и ждём готовности по логу
	// Порт 5432 мапится на случайный порт хоста
	container, err := postgres.Run(ctx, postgresImage,
		postgres.WithDatabase(postgresDB),
		postgres.WithUsername(postgresUser),
		postgres.WithPassword(postgresPassword),

		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(startupTimeout),
		),
	)
	if err != nil {
		logger.Fatal(ctx, "не удалось запустить контейнер PostgreSQL", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер PostgreSQL запущен")

	// Получаем DSN с уже подставленным случайным портом
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Container: container})
		logger.Fatal(ctx, "не удалось получить строку подключения", zap.Error(err))
	}

	// Пул соединений к тестовой БД
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Container: container})
		logger.Fatal(ctx, "не удалось создать пул соединений", zap.Error(err))
	}

	// Применяем те же миграции, что и сервис при старте
	sqlDB := stdlib.OpenDBFromPool(pool)
	migrationsDir := filepath.Join(path.GetProjectRoot(), "order/migrations")

	if err := migratorpg.New(sqlDB, migrationsDir).Up(); err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Container: container, Pool: pool})
		logger.Fatal(ctx, "не удалось применить миграции", zap.Error(err))
	}
	logger.Info(ctx, "✅ Миграции применены")

	// Настоящий репозиторий поверх реальной БД
	repo := repositoryorder.NewRepository(dbpg.NewPGClient(pool))

	logger.Info(ctx, "🎉 Тестовое окружение готово")

	return &TestEnvironment{
		Container: container,
		Pool:      pool,
		Repo:      repo,
	}
}
