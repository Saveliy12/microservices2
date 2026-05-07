package part

import (
	"context"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
)

func (s *Service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.partRepository.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
