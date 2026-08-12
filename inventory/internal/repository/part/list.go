package part

import (
	"context"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repositoryConverter "github.com.Saveliy12/microservices2.git/inventory/internal/repository/converter"
	repositoryModel "github.com.Saveliy12/microservices2.git/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	repoFilter := repositoryConverter.ModelFilterToRepoFilter(filter)
	query := buildPartsFilter(repoFilter)

	var parts []repositoryModel.Part
	if err := r.collection.Find(ctx, query, &parts); err != nil {
		return nil, err
	}

	return repositoryConverter.RepoPartsToModel(parts), nil
}

func buildPartsFilter(filter repositoryModel.PartsFilter) bson.M {
	query := bson.M{}

	if len(filter.UUIDs) > 0 {
		query["uuid"] = bson.M{"$in": filter.UUIDs}
	}

	if len(filter.Names) > 0 {
		query["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		query["category"] = bson.M{"$in": filter.Categories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		query["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	return query
}
