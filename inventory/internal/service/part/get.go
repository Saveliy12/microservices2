package part

import (
	"context"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
)

func (s *Service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.partRepository.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
