package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, id uuid.UUID) (model.Part, error)
	GetBatch(ctx context.Context, ids []uuid.UUID) ([]model.Part, error)
	GetAll(ctx context.Context) ([]model.Part, error)
}
