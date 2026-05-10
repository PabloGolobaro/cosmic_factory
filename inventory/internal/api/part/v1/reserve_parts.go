package v1

import (
	"context"

	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ReserveParts(ctx context.Context, req *inventoryv1.ReservePartsRequest) (*inventoryv1.ReservePartsResponse, error) {
	if err := a.PartService.ReserveParts(ctx, req.GetUuids()); err != nil {
		return nil, mapToGRPCError(err)
	}

	return &inventoryv1.ReservePartsResponse{}, nil
}
