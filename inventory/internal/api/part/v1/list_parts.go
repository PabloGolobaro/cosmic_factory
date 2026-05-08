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

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	partType, err := converter.PartTypeFromProto(req.GetPartType())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	parts, err := a.PartService.List(ctx, req.GetUuids(), partType)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidUUID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errs.ErrPartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	proto := make([]*inventoryv1.Part, 0, len(parts))
	for _, p := range parts {
		proto = append(proto, converter.PartToProto(p))
	}
	return &inventoryv1.ListPartsResponse{Parts: proto}, nil
}
