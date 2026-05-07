package part

import (
	"context"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repositoryConverter "github.com.Saveliy12/microservices2.git/inventory/internal/repository/converter"
)

func (r *Repository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return repositoryConverter.RepoPartToModel(part), nil
}
