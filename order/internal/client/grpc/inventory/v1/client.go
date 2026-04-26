package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

type inventoryClient struct {
	client inventoryv1.InventoryServiceClient
}

func NewInventoryClient(client inventoryv1.InventoryServiceClient) *inventoryClient {
	return &inventoryClient{client: client}
}

func (i inventoryClient) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	resp, err := i.client.ListParts(ctx, &inventoryv1.ListPartsRequest{Uuids: uuids})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, errs.ErrPartNotFound
			case codes.InvalidArgument:
				return nil, errs.ErrInvalidUUID
			}
		}
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}
	parts := make([]model.Part, 0, len(resp.GetParts()))
	for _, p := range resp.GetParts() {
		parts = append(parts, partFromProto(p))
	}
	return parts, nil
}
