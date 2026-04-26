package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

var modelToProtoPartType = map[model.PartType]inventoryv1.PartType{
	model.PartTypeHull:   inventoryv1.PartType_PART_TYPE_HULL,
	model.PartTypeEngine: inventoryv1.PartType_PART_TYPE_ENGINE,
	model.PartTypeShield: inventoryv1.PartType_PART_TYPE_SHIELD,
	model.PartTypeWeapon: inventoryv1.PartType_PART_TYPE_WEAPON,
}

func PartToProto(p model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          p.UUID.String(),
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      modelToProtoPartType[p.PartType],
		StockQuantity: p.StockQuantity,
		CreatedAt:     timestamppb.New(p.CreatedAt),
	}
}
