package valueobject

// PartFilter задаёт параметры фильтрации деталей.
type PartFilter struct {
	UUIDs    []string
	PartType PartType
}
