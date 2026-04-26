package v1

import (
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	PartService PartService
	inventoryv1.UnimplementedInventoryServiceServer
}

func NewApi(partService PartService) *api {
	return &api{PartService: partService}
}
