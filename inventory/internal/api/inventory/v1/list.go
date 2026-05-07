package v1

import (
	"context"

	"github.com.Saveliy12/microservices2.git/inventory/internal/converter"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	parts, err := a.inventoryService.ListParts(ctx, converter.FilterFromProto(req.GetFilter()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &inventoryv1.ListPartsResponse{Parts: converter.PartsToProto(parts)}, nil
}
