package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ValidateCompatibility(ctx context.Context, req *inventoryv1.ValidateCompatibilityRequest) (*inventoryv1.ValidateCompatibilityResponse, error) {
	if err := a.PartService.ValidateCompatibility(ctx, req.GetUuids()); err != nil {
		return nil, mapToGRPCError(err)
	}

	return &inventoryv1.ValidateCompatibilityResponse{}, nil
}

func mapToGRPCError(err error) error {
	switch {
	case errors.Is(err, errs.ErrInvalidUUID):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrPartNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, errs.ErrOutOfStock):
		return status.Error(codes.ResourceExhausted, err.Error())
	case errors.Is(err, errs.ErrNothingToRelease):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, errs.ErrIncompatibleParts):
		return status.Error(codes.FailedPrecondition, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
