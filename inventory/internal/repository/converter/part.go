package converter

import (
	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/record"
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

func PartToRecord(p model.Part) record.PartRecord {
	return record.PartRecord{
		UUID:          p.UUID.String(),
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      partTypeToStr[p.PartType],
		StockQuantity: p.StockQuantity,
		CreatedAt:     p.CreatedAt,
	}
}

func PartFromRecord(r record.PartRecord) model.Part {
	return model.Part{
		UUID:          uuid.MustParse(r.UUID),
		Name:          r.Name,
		Description:   r.Description,
		Price:         r.Price,
		PartType:      strToPartType[r.PartType],
		StockQuantity: r.StockQuantity,
		CreatedAt:     r.CreatedAt,
	}
}
