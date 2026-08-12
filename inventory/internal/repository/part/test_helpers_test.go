package part_test

import (
	"time"

	repositoryModel "github.com.Saveliy12/microservices2.git/inventory/internal/repository/model"
)

func testPart() repositoryModel.Part {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	return repositoryModel.Part{
		UUID:          "uuid-1",
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
	}
}
