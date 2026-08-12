package order

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	repositoryConverter "github.com.Saveliy12/microservices2.git/order/internal/repository/converter"
	repositoryModel "github.com.Saveliy12/microservices2.git/order/internal/repository/model"
	"github.com/Masterminds/squirrel"
)

func (r *Repository) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	builder := r.builder.Select("order_uuid", "user_uuid", "total_price", "status", "transaction_uuid", "payment_method", "part_uuids").
		From("orders").
		Where(squirrel.Eq{"order_uuid": orderUUID})

	query, args, err := builder.ToSql()
	if err != nil {
		return model.Order{}, err
	}

	var repoOrder repositoryModel.Order
	err = r.client.QueryRow(ctx, query, args...).Scan(
		&repoOrder.OrderUUID,
		&repoOrder.UserUUID,
		&repoOrder.TotalPrice,
		&repoOrder.Status,
		&repoOrder.TransactionUUID,
		&repoOrder.PaymentMethod,
		&repoOrder.PartUUIDs,
	)
	if err != nil {
		return model.Order{}, err
	}

	return repositoryConverter.OrderFromRepoModel(repoOrder), nil
}
