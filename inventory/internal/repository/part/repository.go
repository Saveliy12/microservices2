package part

import (
	db "github.com.Saveliy12/microservices2.git/inventory/internal/client/db"
	repository "github.com.Saveliy12/microservices2.git/inventory/internal/repository"
)

var _ repository.PartRepository = (*Repository)(nil)

type Repository struct {
	collection db.Collection
}

func NewRepository(collection db.Collection) *Repository {
	return &Repository{
		collection: collection,
	}
}
