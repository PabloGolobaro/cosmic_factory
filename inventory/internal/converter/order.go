package converter

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
)

var modelToProtoPartType = map[valueobject.PartType]inventoryv1.PartType{
	valueobject.PartTypeHull:   inventoryv1.PartType_PART_TYPE_HULL,
	valueobject.PartTypeEngine: inventoryv1.PartType_PART_TYPE_ENGINE,
	valueobject.PartTypeShield: inventoryv1.PartType_PART_TYPE_SHIELD,
	valueobject.PartTypeWeapon: inventoryv1.PartType_PART_TYPE_WEAPON,
}

var protoToModelPartType = map[inventoryv1.PartType]valueobject.PartType{
	inventoryv1.PartType_PART_TYPE_UNSPECIFIED: valueobject.PartTypeUnspecified,
	inventoryv1.PartType_PART_TYPE_HULL:        valueobject.PartTypeHull,
	inventoryv1.PartType_PART_TYPE_ENGINE:      valueobject.PartTypeEngine,
	inventoryv1.PartType_PART_TYPE_SHIELD:      valueobject.PartTypeShield,
	inventoryv1.PartType_PART_TYPE_WEAPON:      valueobject.PartTypeWeapon,
}

func PartTypeFromProto(pt inventoryv1.PartType) (valueobject.PartType, error) {
	mt, ok := protoToModelPartType[pt]
	if !ok {
		return "", fmt.Errorf("неизвестный тип детали %v: %w", pt, errs.ErrInvalidProperties)
	}

	return mt, nil
}

func PartToProto(p entity.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          p.UUID(),
		Name:          p.Name(),
		Description:   p.Description(),
		Price:         p.Price(),
		PartType:      modelToProtoPartType[p.PartType()],
		StockQuantity: p.StockQuantity(),
		CreatedAt:     timestamppb.New(p.CreatedAt()),
	}
}
