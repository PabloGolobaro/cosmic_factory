package v1

import (
	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

func partTypeFromProto(pt inventoryv1.PartType) model.PartType {
	return model.PartType(pt)
}

func partFromProto(p *inventoryv1.Part) model.Part {
	return model.Part{
		UUID:          uuid.MustParse(p.GetUuid()),
		Name:          p.GetName(),
		Description:   p.GetDescription(),
		Price:         p.GetPrice(),
		PartType:      partTypeFromProto(p.GetPartType()),
		StockQuantity: p.GetStockQuantity(),
		CreatedAt:     p.GetCreatedAt().AsTime(),
	}
}
