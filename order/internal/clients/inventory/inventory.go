package inventory

import (
	"context"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
	service "github.com.Saveliy12/microservices2.git/order/internal/service"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ service.InventoryClient = (*Client)(nil)

type Client struct {
	client inventoryv1.InventoryServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{
		client: inventoryv1.NewInventoryServiceClient(conn),
	}
}

func (c *Client) ListParts(ctx context.Context, partUUIDs []string) ([]model.PartQuote, error) {
	req := &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: make([]*wrapperspb.StringValue, 0, len(partUUIDs)),
		},
	}
	for _, uuid := range partUUIDs {
		req.Filter.Uuids = append(req.Filter.Uuids, wrapperspb.String(uuid))
	}

	resp, err := c.client.ListParts(ctx, req)
	if err != nil {
		return nil, err
	}

	parts := make([]model.PartQuote, 0, len(resp.GetParts()))
	for _, part := range resp.GetParts() {
		parts = append(parts, model.PartQuote{
			UUID:  part.GetUuid(),
			Price: part.GetPrice(),
		})
	}

	return parts, nil
}
