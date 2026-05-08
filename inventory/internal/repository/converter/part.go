package converter

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/repository/record"
)

var partTypeToStr = map[valueobject.PartType]string{
	valueobject.PartTypeHull:   "HULL",
	valueobject.PartTypeEngine: "ENGINE",
	valueobject.PartTypeShield: "SHIELD",
	valueobject.PartTypeWeapon: "WEAPON",
}

func PartToRecord(p entity.Part) (record.PartRecord, error) {
	propsRec := partPropertiesToRecord(p.Properties())
	propsJSON, err := json.Marshal(propsRec)
	if err != nil {
		return record.PartRecord{}, fmt.Errorf("сериализовать свойства: %w", err)
	}

	return record.PartRecord{
		UUID:          p.UUID(),
		Name:          p.Name(),
		Description:   p.Description(),
		Price:         p.Price(),
		PartType:      partTypeToStr[p.PartType()],
		StockQuantity: p.StockQuantity(),
		Reserved:      p.Reserved(),
		Properties:    propsJSON,
		CreatedAt:     p.CreatedAt(),
	}, nil
}

func PartFromRecord(r record.PartRecord) (entity.Part, error) {
	if _, err := uuid.Parse(r.UUID); err != nil {
		return entity.Part{}, fmt.Errorf("некорректный UUID записи: %w", err)
	}

	pt, err := valueobject.NewPartType(r.PartType)
	if err != nil {
		return entity.Part{}, fmt.Errorf("конвертировать тип детали: %w", err)
	}

	var propsRec record.PartPropertiesRecord
	if len(r.Properties) > 0 {
		if err = json.Unmarshal(r.Properties, &propsRec); err != nil {
			return entity.Part{}, fmt.Errorf("десериализовать свойства: %w", err)
		}
	}

	props, err := partPropertiesFromRecord(propsRec)
	if err != nil {
		return entity.Part{}, fmt.Errorf("конвертировать свойства: %w", err)
	}

	return entity.RestorePart(
		r.UUID, r.Name, r.Description, pt, r.Price,
		int(r.StockQuantity), r.Reserved, props, r.CreatedAt,
	), nil
}

func partPropertiesFromRecord(rec record.PartPropertiesRecord) (valueobject.PartProperties, error) {
	switch {
	case rec.Hull != nil:
		return valueobject.NewHullProperties(rec.Hull.Strength)
	case rec.Engine != nil:
		return valueobject.NewEngineProperties(
			valueobject.EngineClass(rec.Engine.Class),
			rec.Engine.RequiredStrength,
		)
	case rec.Shield != nil:
		return valueobject.NewShieldProperties(valueobject.ShieldType(rec.Shield.ShieldType))
	case rec.Weapon != nil:
		return valueobject.NewWeaponProperties(valueobject.WeaponType(rec.Weapon.WeaponType))
	default:
		return valueobject.PartProperties{}, nil
	}
}

func partPropertiesToRecord(p valueobject.PartProperties) record.PartPropertiesRecord {
	var rec record.PartPropertiesRecord
	switch {
	case p.Hull() != nil:
		rec.Hull = &record.HullPropertiesRecord{Strength: p.Hull().Strength()}
	case p.Engine() != nil:
		rec.Engine = &record.EnginePropertiesRecord{
			Class:            string(p.Engine().Class()),
			RequiredStrength: p.Engine().RequiredStrength(),
		}
	case p.Shield() != nil:
		rec.Shield = &record.ShieldPropertiesRecord{ShieldType: string(p.Shield().ShieldType())}
	case p.Weapon() != nil:
		rec.Weapon = &record.WeaponPropertiesRecord{WeaponType: string(p.Weapon().WeaponType())}
	}
	return rec
}
