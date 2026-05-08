package converter

import (
	"fmt"

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

func PartFromRecord(r record.PartRecord) (model.Part, error) {
	id, err := uuid.Parse(r.UUID)
	if err != nil {
		return model.Part{}, fmt.Errorf("некорректный UUID записи: %w", err)
	}

	pt, err := model.NewPartType(r.PartType)
	if err != nil {
		return model.Part{}, err
	}

	return model.Part{
		UUID:          id,
		Name:          r.Name,
		Description:   r.Description,
		Price:         r.Price,
		PartType:      pt,
		StockQuantity: r.StockQuantity,
		CreatedAt:     r.CreatedAt,
	}, nil
}
