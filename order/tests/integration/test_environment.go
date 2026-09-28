package integration

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	"github.com.Saveliy12/microservices2.git/order/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// TestEnvironment хранит ресурсы тестового окружения
type TestEnvironment struct {
	Container *postgres.PostgresContainer
	Pool      *pgxpool.Pool
	Repo      repository.OrderRepository
}

// ClearOrders очищает таблицу между тестами, чтобы данные не текли между кейсами
func (env *TestEnvironment) ClearOrders(ctx context.Context) error {
	_, err := env.Pool.Exec(ctx, "TRUNCATE TABLE orders")

	return err
}

// BuildTestOrder создаёт заказ с валидными данными для тестов
// TransactionUUID и PaymentMethod не заполняются
func BuildTestOrder() model.Order {
	return model.Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   uuid.NewString(),
		PartUUIDs:  []string{uuid.NewString(), uuid.NewString()},
		TotalPrice: 150.5,
		Status:     model.StatusPendingPayment,
	}
}
