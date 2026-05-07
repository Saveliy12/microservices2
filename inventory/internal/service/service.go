package service

import (
	"context"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
