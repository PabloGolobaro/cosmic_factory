package model

import "github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"

// PartFilter задаёт параметры фильтрации деталей.
type PartFilter struct {
	UUIDs    []string
	PartType valueobject.PartType
}
