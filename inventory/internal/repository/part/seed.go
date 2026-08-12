package part

import (
	"context"
	"time"

	repositoryModel "github.com.Saveliy12/microservices2.git/inventory/internal/repository/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SeedParts вставляет стартовые детали в коллекцию, если она пуста
// Работает напрямую с *mongo.Collection, чтобы не расширять интерфейс db.Collection
func SeedParts(ctx context.Context, collection *mongo.Collection) error {
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	parts := make([]any, 0, 4)
	for _, part := range seedParts() {
		parts = append(parts, part)
	}

	_, err = collection.InsertMany(ctx, parts)
	return err
}

func seedParts() []repositoryModel.Part {
	now := time.Now()

	return []repositoryModel.Part{
		{
			UUID:          uuid.NewString(),
			Name:          "Engine V8",
			Description:   "Powerful 8-cylinder engine",
			Price:         50000.00,
			StockQuantity: 10,
			Category:      repositoryModel.CategoryEngine,
			Dimensions: repositoryModel.Dimensions{
				Length: 600,
				Width:  450,
				Height: 550,
				Weight: 250,
			},
			Manufacturer: repositoryModel.Manufacturer{
				Name:    "TechMotor Inc.",
				Country: "USA",
				Website: "https://techmotor.com",
			},
			Tags:      []string{"high-performance", "reliable"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.NewString(),
			Name:          "Turbo Charger",
			Description:   "High-performance turbo charger",
			Price:         15000.00,
			StockQuantity: 25,
			Category:      repositoryModel.CategoryEngine,
			Dimensions: repositoryModel.Dimensions{
				Length: 250,
				Width:  200,
				Height: 280,
				Weight: 45,
			},
			Manufacturer: repositoryModel.Manufacturer{
				Name:    "SpeedForce Ltd.",
				Country: "Germany",
				Website: "https://speedforce.de",
			},
			Tags:      []string{"turbo", "boost"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.NewString(),
			Name:          "Fuel Pump",
			Description:   "Electric fuel pump",
			Price:         3500.00,
			StockQuantity: 50,
			Category:      repositoryModel.CategoryFuel,
			Dimensions: repositoryModel.Dimensions{
				Length: 120,
				Width:  80,
				Height: 100,
				Weight: 8,
			},
			Manufacturer: repositoryModel.Manufacturer{
				Name:    "AutoSupply Corp.",
				Country: "Japan",
				Website: "https://autosupply.jp",
			},
			Tags:      []string{"fuel-system", "essential"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.NewString(),
			Name:          "Carbon Fiber Wing",
			Description:   "Lightweight carbon fiber aerodynamic wing",
			Price:         25000.00,
			StockQuantity: 5,
			Category:      repositoryModel.CategoryWing,
			Dimensions: repositoryModel.Dimensions{
				Length: 1800,
				Width:  1200,
				Height: 150,
				Weight: 12,
			},
			Manufacturer: repositoryModel.Manufacturer{
				Name:    "Carbon Racing",
				Country: "UK",
				Website: "https://carbonracing.co.uk",
			},
			Tags:      []string{"aerodynamic", "premium"},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}
