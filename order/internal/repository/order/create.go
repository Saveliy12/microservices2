package order

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repositoryConverter "github.com.Saveliy12/microservices2.git/order/internal/repository/converter"
)

func (r *Repository) CreateOrder(ctx context.Context, order model.Order) error {
	repoOrder := repositoryConverter.OrderToRepoModel(order)

	builder := r.builder.Insert("orders").
		Columns("order_uuid", "user_uuid", "total_price", "status", "part_uuids").
		Values(repoOrder.OrderUUID, repoOrder.UserUUID, repoOrder.TotalPrice, repoOrder.Status, repoOrder.PartUUIDs)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.client.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
