package order

import (
	"github.com.Saveliy12/microservices2.git/order/internal/client/db"
	repository "github.com.Saveliy12/microservices2.git/order/internal/repository"
	"github.com/Masterminds/squirrel"
)

var _ repository.OrderRepository = (*Repository)(nil)

type Repository struct {
	client  db.Client
	builder squirrel.StatementBuilderType
}

func NewRepository(client db.Client) *Repository {
	return &Repository{
		client:  client,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
