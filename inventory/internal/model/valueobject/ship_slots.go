package valueobject

// ShipSlots описывает UUID деталей по слотам корабля для проверки совместимости.
// HullUUID и EngineUUID — обязательны; ShieldUUID и WeaponUUID — опциональны (пустая строка = слот не используется).
type ShipSlots struct {
	HullUUID   string
	EngineUUID string
	ShieldUUID string
	WeaponUUID string
}
