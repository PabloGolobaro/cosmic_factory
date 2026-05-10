package v1

import (
	"context"

	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ReleaseParts(ctx context.Context, req *inventoryv1.ReleasePartsRequest) (*inventoryv1.ReleasePartsResponse, error) {
	if err := a.PartService.ReleaseParts(ctx, req.GetUuids()); err != nil {
		return nil, mapToGRPCError(err)
	}

	return &inventoryv1.ReleasePartsResponse{}, nil
}
