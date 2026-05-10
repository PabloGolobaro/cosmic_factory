package model

import "github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"

// ResolvedShipSlots содержит разрешённые детали по именованным слотам корабля.
// Hull и Engine — обязательны; Shield и Weapon — nil если слот не используется.
type ResolvedShipSlots struct {
	Hull   entity.Part
	Engine entity.Part
	Shield *entity.Part
	Weapon *entity.Part
}

// ShipSlots описывает UUID деталей по слотам корабля для проверки совместимости.
// HullUUID и EngineUUID — обязательны; ShieldUUID и WeaponUUID — опциональны (пустая строка = слот не используется).
type ShipSlots struct {
	HullUUID   string
	EngineUUID string
	ShieldUUID string
	WeaponUUID string
}
