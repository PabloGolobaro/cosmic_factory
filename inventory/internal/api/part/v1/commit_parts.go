package v1

import (
	"context"

	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

func (a *api) CommitParts(ctx context.Context, req *inventoryv1.CommitPartsRequest) (*inventoryv1.CommitPartsResponse, error) {
	if err := a.PartService.CommitParts(ctx, req.GetUuids()); err != nil {
		return nil, mapToGRPCError(err)
	}
	return &inventoryv1.CommitPartsResponse{}, nil
}
