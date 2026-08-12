package order

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repositoryConverter "github.com.Saveliy12/microservices2.git/order/internal/repository/converter"
	"github.com/Masterminds/squirrel"
)

func (r *Repository) UpdateOrder(ctx context.Context, order model.Order) error {
	repoOrder := repositoryConverter.OrderToRepoModel(order)

	builder := r.builder.Update("orders").
		Set("user_uuid", repoOrder.UserUUID).
		Set("total_price", repoOrder.TotalPrice).
		Set("status", repoOrder.Status).
		Set("transaction_uuid", repoOrder.TransactionUUID).
		Set("payment_method", repoOrder.PaymentMethod).
		Set("part_uuids", repoOrder.PartUUIDs).
		Where(squirrel.Eq{"order_uuid": repoOrder.OrderUUID})

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
