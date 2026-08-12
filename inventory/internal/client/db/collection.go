package db

import "context"

type Collection interface {
	FindOne(ctx context.Context, filter, result any) error
	Find(ctx context.Context, filter, results any) error
}
