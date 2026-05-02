package converter

import (
	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/record"
)

var partTypeToStr = map[model.PartType]string{
	model.PartTypeHull:   "HULL",
	model.PartTypeEngine: "ENGINE",
	model.PartTypeShield: "SHIELD",
	model.PartTypeWeapon: "WEAPON",
}

var strToPartType = map[string]model.PartType{
	"HULL":   model.PartTypeHull,
	"ENGINE": model.PartTypeEngine,
	"SHIELD": model.PartTypeShield,
	"WEAPON": model.PartTypeWeapon,
}

func OrderItemToRecord(i model.OrderItem) record.OrderItemRecord {
	return record.OrderItemRecord{
		UUID:      i.UUID.String(),
		OrderUUID: i.OrderUUID.String(),
		PartUUID:  i.PartUUID.String(),
		PartType:  partTypeToStr[i.PartType],
		Price:     i.Price,
		CreatedAt: i.CreatedAt,
	}
}

func OrderItemFromRecord(r record.OrderItemRecord) model.OrderItem {
	return model.OrderItem{
		UUID:      uuid.MustParse(r.UUID),
		OrderUUID: uuid.MustParse(r.OrderUUID),
		PartUUID:  uuid.MustParse(r.PartUUID),
		PartType:  strToPartType[r.PartType],
		Price:     r.Price,
		CreatedAt: r.CreatedAt,
	}
}
