package v1

import (
	"context"
	"errors"

	"github.com.Saveliy12/microservices2.git/inventory/internal/converter"
	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	part, err := a.inventoryService.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &inventoryv1.GetResponse{Part: converter.PartToProto(part)}, nil
}
