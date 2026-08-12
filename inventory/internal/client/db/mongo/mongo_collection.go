package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoCollection struct {
	coll *mongo.Collection
}

// NewMongoCollection оборачивает коллекцию MongoDB в интерфейс db.Collection
func NewMongoCollection(coll *mongo.Collection) *mongoCollection {
	return &mongoCollection{coll: coll}
}

func (m *mongoCollection) FindOne(ctx context.Context, filter, result any) error {
	return m.coll.FindOne(ctx, filter).Decode(result)
}

func (m *mongoCollection) Find(ctx context.Context, filter, results any) error {
	cursor, err := m.coll.Find(ctx, filter)
	if err != nil {
		return err
	}

	return cursor.All(ctx, results)
}
