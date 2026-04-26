package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/converter"
	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	part, err := a.PartService.Get(ctx, req.GetUuid())
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidUUID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errs.ErrPartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &inventoryv1.GetPartResponse{Part: converter.PartToProto(part)}, nil
}
