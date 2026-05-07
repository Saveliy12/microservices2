package part

import (
	"sync"
	"time"

	repository "github.com.Saveliy12/microservices2.git/inventory/internal/repository"
	repositoryModel "github.com.Saveliy12/microservices2.git/inventory/internal/repository/model"
	"github.com/google/uuid"
)

var _ repository.PartRepository = (*Repository)(nil)

type Repository struct {
	mu    sync.RWMutex
	parts map[string]repositoryModel.Part
}

func NewRepository() *Repository {
	parts := make(map[string]repositoryModel.Part)

	// Seed data — use the map key as the Part.UUID so lookups by UUID work
	id1 := uuid.NewString()
	parts[id1] = repositoryModel.Part{
		UUID:          id1,
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id2 := uuid.NewString()
	parts[id2] = repositoryModel.Part{
		UUID:          id2,
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id3 := uuid.NewString()
	parts[id3] = repositoryModel.Part{
		UUID:          id3,
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id4 := uuid.NewString()
	parts[id4] = repositoryModel.Part{
		UUID:          id4,
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &Repository{
		parts: parts,
	}
}
