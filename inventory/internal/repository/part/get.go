package part

import (
	"context"
	"errors"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repositoryConverter "github.com.Saveliy12/microservices2.git/inventory/internal/repository/converter"
	repositoryModel "github.com.Saveliy12/microservices2.git/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	var part repositoryModel.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}, &part)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Part{}, model.ErrPartNotFound
	}
	if err != nil {
		return model.Part{}, err
	}

	return repositoryConverter.RepoPartToModel(part), nil
}
