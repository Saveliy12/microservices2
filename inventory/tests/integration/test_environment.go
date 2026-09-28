package integration

import (
	"context"
	"time"

	repositoryModel "github.com.Saveliy12/microservices2.git/inventory/internal/repository/model"
	"github.com.Saveliy12/microservices2.git/platform/pkg/testcontainers/app"
	"github.com.Saveliy12/microservices2.git/platform/pkg/testcontainers/mongo"
	"github.com.Saveliy12/microservices2.git/platform/pkg/testcontainers/network"
	"github.com/google/uuid"
)

// TestEnvironment хранит ресурсы тестового окружения: сеть, MongoDB и контейнер приложения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

// InsertTestPart вставляет деталь напрямую в MongoDB и возвращает её UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := uuid.NewString()
	now := time.Now()

	part := repositoryModel.Part{
		UUID:          partUUID,
		Name:          "Test Part",
		Description:   "Part inserted directly by e2e test",
		Price:         123.45,
		StockQuantity: 7,
		Category:      repositoryModel.CategoryEngine,
		Dimensions: repositoryModel.Dimensions{
			Length: 10,
			Width:  20,
			Height: 30,
			Weight: 40,
		},
		Manufacturer: repositoryModel.Manufacturer{
			Name:    "Test Manufacturer Inc.",
			Country: "Testland",
			Website: "https://test.example",
		},
		Tags:      []string{"e2e", "test"},
		Metadata:  map[string]repositoryModel.MetadataValue{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := env.Mongo.Client().
		Database(mongoDatabase).
		Collection(partsCollectionName).
		InsertOne(ctx, part)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}
