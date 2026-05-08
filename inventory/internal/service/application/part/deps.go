package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

type PartRepository interface {
	Get(ctx context.Context, id uuid.UUID) (entity.Part, error)
	GetBatch(ctx context.Context, filter valueobject.PartFilter) ([]entity.Part, error)
}

type CompatibilityChecker interface {
	Check(parts []entity.Part) error
}
