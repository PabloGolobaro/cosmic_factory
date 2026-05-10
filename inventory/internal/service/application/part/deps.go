package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
)

type PartRepository interface {
	Get(ctx context.Context, id uuid.UUID) (entity.Part, error)
	GetBatch(ctx context.Context, filter model.PartFilter) ([]entity.Part, error)
	UpdateReservedBatch(ctx context.Context, parts []entity.Part) error
}

type CompatibilityChecker interface {
	Check(slots model.ResolvedShipSlots) error
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
