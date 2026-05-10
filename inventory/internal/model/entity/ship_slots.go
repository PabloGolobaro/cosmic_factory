package entity

// ResolvedShipSlots содержит разрешённые детали по именованным слотам корабля.
// Hull и Engine — обязательны; Shield и Weapon — nil если слот не используется.
type ResolvedShipSlots struct {
	Hull   Part
	Engine Part
	Shield *Part
	Weapon *Part
}
